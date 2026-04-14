package main

import (
	"context"
	"encoding/json"
	"errors"
	"event_processor/model"
	"event_processor/utils"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	workersCount       = 5
	eventBufferSize    = 64
	producerMinDelayMs = 50
	producerMaxDelayMs = 200
	serverAddr         = "127.0.0.1:8080"
)

var (
	processorMinDelay = 75 * time.Millisecond
	processorMaxDelay = 200 * time.Millisecond
)

type MetricsStore struct {
	mu      sync.RWMutex
	metrics model.Metrics
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		metrics: model.Metrics{
			EventTypeCount: make(map[model.EventType]int),
			PerUserCount:   make(map[int]int),
		},
	}
}

func (s *MetricsStore) Record(event model.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metrics.TotalEvents++
	s.metrics.EventTypeCount[event.Type]++
	s.metrics.PerUserCount[event.UserId]++
}

func (s *MetricsStore) Snapshot() model.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	eventTypeCount := make(map[model.EventType]int, len(s.metrics.EventTypeCount))
	for eventType, count := range s.metrics.EventTypeCount {
		eventTypeCount[eventType] = count
	}

	perUserCount := make(map[int]int, len(s.metrics.PerUserCount))
	for userID, count := range s.metrics.PerUserCount {
		perUserCount[userID] = count
	}

	return model.Metrics{
		TotalEvents:    s.metrics.TotalEvents,
		EventTypeCount: eventTypeCount,
		PerUserCount:   perUserCount,
	}
}

func randomDuration(minMs, maxMs int) time.Duration {
	return time.Duration(rand.IntN(maxMs-minMs+1)+minMs) * time.Millisecond
}

func processEvent(event model.Event) {
	_ = event
	delayRange := int64(processorMaxDelay - processorMinDelay)
	time.Sleep(processorMinDelay + time.Duration(rand.Int64N(delayRange+1)))
}

func produce(ctx context.Context, eventCh chan<- model.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(eventCh)

	for {
		timer := time.NewTimer(randomDuration(producerMinDelayMs, producerMaxDelayMs))

		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			log.Println("stopping producer")
			return
		case <-timer.C:
		}

		event := model.Event{
			EventId:   utils.GenerateEventId(),
			UserId:    utils.GenerateUserId(),
			Type:      utils.GenerateEventType(),
			Timestamp: time.Now().Unix(),
		}

		select {
		case eventCh <- event:
		case <-ctx.Done():
			log.Println("stopping producer")
			return
		}
	}
}

func worker(workerID int, eventCh <-chan model.Event, store *MetricsStore, wg *sync.WaitGroup) {
	defer wg.Done()

	for event := range eventCh {
		processEvent(event)
		store.Record(event)
	}

	log.Printf("worker %d stopped\n", workerID)
}

func metricsHandler(store *MetricsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(store.Snapshot()); err != nil {
			http.Error(w, fmt.Sprintf("json encode error: %v", err), http.StatusInternalServerError)
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eventCh := make(chan model.Event, eventBufferSize)
	store := NewMetricsStore()

	var producerWG sync.WaitGroup
	producerWG.Add(1)
	go produce(ctx, eventCh, &producerWG)

	var workerWG sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		workerWG.Add(1)
		go worker(i+1, eventCh, store, &workerWG)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metricsHandler(store))

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("metrics server listening on %s\n", serverAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
		close(serverErrCh)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err, ok := <-serverErrCh:
		if ok && err != nil {
			log.Fatalf("http server failed: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("http server shutdown error: %v\n", err)
	}

	producerWG.Wait()
	workerWG.Wait()
	log.Printf("final metrics: %+v\n", store.Snapshot())
}
