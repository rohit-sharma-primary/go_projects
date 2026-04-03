package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rate_limiter/model"
	"sync"
	"time"
)

const (
	maxRetries     = 5
	baseTxnURL     = "http://127.0.0.1:8081/transactions"
	maxRequestsRPS = 3
	workerCount    = 3
	requestTimeout = 2 * time.Second
	baseBackoff    = 100 * time.Millisecond
)

var retryableStatusCodes = map[int]bool{
	http.StatusInternalServerError: true,
	http.StatusTooManyRequests:     true,
}

type pageJob struct {
	Page int
}

type pageResult struct {
	Page model.Page
	Err  error
}

type rateLimiter struct {
	ticker *time.Ticker
	tokens <-chan time.Time
}

func newRateLimiter(requestsPerSecond int) *rateLimiter {
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
	return &rateLimiter{ticker: ticker, tokens: ticker.C}
}

func (r *rateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.tokens:
		return nil
	}
}

func (r *rateLimiter) Stop() {
	r.ticker.Stop()
}

type backendError struct {
	StatusCode int
}

func (e backendError) Error() string {
	return fmt.Sprintf("backend returned status %d", e.StatusCode)
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var be backendError
	if errors.As(err, &be) {
		return retryableStatusCodes[be.StatusCode]
	}

	return errors.Is(err, context.DeadlineExceeded)
}

func backoffDuration(attempt int) time.Duration {
	return baseBackoff << (attempt - 1)
}

func fetchPageWithRetry(parentCtx context.Context, limiter *rateLimiter, page int) (model.Page, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries+1; attempt++ {
		if err := limiter.Wait(parentCtx); err != nil {
			return model.Page{}, err
		}

		reqCtx, cancel := context.WithTimeout(parentCtx, requestTimeout)
		data, err := fetchPage(reqCtx, page)
		cancel()

		if err == nil {
			return data, nil
		}

		lastErr = err
		if attempt == maxRetries+1 || !isRetryableError(err) {
			return model.Page{}, err
		}

		select {
		case <-parentCtx.Done():
			return model.Page{}, parentCtx.Err()
		case <-time.After(backoffDuration(attempt)):
		}
	}

	return model.Page{}, lastErr
}

func fetchPage(ctx context.Context, page int) (model.Page, error) {
	url := fmt.Sprintf("%s?page=%d", baseTxnURL, page)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Page{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Page{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Page{}, backendError{StatusCode: resp.StatusCode}
	}

	var pageData model.Page
	if err := json.NewDecoder(resp.Body).Decode(&pageData); err != nil {
		return model.Page{}, err
	}

	return pageData, nil
}

func worker(ctx context.Context, limiter *rateLimiter, jobs <-chan pageJob, results chan<- pageResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		page, err := fetchPageWithRetry(ctx, limiter, job.Page)
		select {
		case <-ctx.Done():
			return
		case results <- pageResult{Page: page, Err: err}:
		}
	}
}

func aggregate(pages []model.Page) model.Result {
	result := model.Result{
		PerUserTotal: make(map[int]int),
	}

	for _, page := range pages {
		for _, txn := range page.Data {
			result.TotalTransactions++
			result.TotalAmount += txn.Amount
			result.PerUserTotal[txn.UserID] += txn.Amount
		}
	}

	if result.TotalTransactions > 0 {
		result.AverageAmount = float64(result.TotalAmount) / float64(result.TotalTransactions)
	}

	return result
}

func buildSummary(ctx context.Context) (model.Result, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	limiter := newRateLimiter(maxRequestsRPS)
	defer limiter.Stop()

	firstPage, err := fetchPageWithRetry(ctx, limiter, 1)
	if err != nil {
		return model.Result{}, err
	}

	pages := make([]model.Page, firstPage.TotalPages)
	pages[0] = firstPage

	if firstPage.TotalPages == 1 {
		return aggregate(pages), nil
	}

	jobs := make(chan pageJob, firstPage.TotalPages-1)
	results := make(chan pageResult, firstPage.TotalPages-1)

	var wg sync.WaitGroup
	workers := workerCount
	if firstPage.TotalPages-1 < workers {
		workers = firstPage.TotalPages - 1
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(ctx, limiter, jobs, results, &wg)
	}

	for page := 2; page <= firstPage.TotalPages; page++ {
		jobs <- pageJob{Page: page}
	}
	close(jobs)

	for received := 2; received <= firstPage.TotalPages; received++ {
		result := <-results
		if result.Err != nil {
			wg.Wait()
			return model.Result{}, result.Err
		}
		pages[result.Page.Page-1] = result.Page
	}

	wg.Wait()
	return aggregate(pages), nil
}

func generateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	summary, err := buildSummary(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/summary", generateReport)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
