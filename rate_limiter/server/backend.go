package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const maxRequestsPerSecond = 3

type backendLimiter struct {
	mu          sync.Mutex
	windowStart time.Time
	count       int
}

func newBackendLimiter() *backendLimiter {
	return &backendLimiter{windowStart: time.Now()}
}

func (l *backendLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.Sub(l.windowStart) >= time.Second {
		l.windowStart = now
		l.count = 0
	}

	l.count++
	return l.count <= maxRequestsPerSecond
}

var reqLimiter = newBackendLimiter()

func shouldSend500Error() bool {
	return rand.Float64() < 0.20
}

func getDelay() time.Duration {
	return time.Duration(100+rand.IntN(901)) * time.Millisecond
}

func readPageData(page int) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileName := filepath.Join(dir, "data", "page"+strconv.Itoa(page)+".json")
	return os.ReadFile(fileName)
}

func paginationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !r.URL.Query().Has("page") {
		http.Error(w, "missing page query param", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 || page > 3 {
		http.Error(w, "incorrect page number", http.StatusBadRequest)
		return
	}

	if !reqLimiter.Allow() {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	time.Sleep(getDelay())

	if shouldSend500Error() {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data, err := readPageData(page)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}

func main() {
	http.HandleFunc("/transactions", paginationHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
