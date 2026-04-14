# Event Processor

This mini project is a real-time event processing system written in Go. It simulates a continuous stream of user events, processes them concurrently with a worker pool, aggregates metrics in memory, exposes those metrics over HTTP, and shuts down cleanly on `Ctrl+C`.

## Objective

The project is designed to model a small event pipeline that:

- generates events continuously
- processes events concurrently
- tracks metrics in real time
- handles backpressure safely
- exposes metrics through an API
- shuts down gracefully

Example event:

```json
{
  "event_id": "e1",
  "user_id": 1,
  "type": "click",
  "timestamp": 1710000000
}
```

## How It Works

### 1. Event Producer

The producer generates events at random intervals between `50ms` and `200ms`.

Each event includes:

- a unique `event_id`
- a random `user_id`
- a random event `type`
- a Unix timestamp

The producer pushes events into a buffered channel.

### 2. Worker Pool

The system starts multiple workers in parallel.

Each worker:

- reads events from the shared channel
- simulates processing time
- updates the metrics store

This demonstrates concurrent event consumption using goroutines.

### 3. Real-Time Aggregation

The application keeps track of:

- `total_events`
- `event_type_count`
- `per_user_count`

Metrics are stored in memory and protected with `sync.RWMutex` to make reads and writes thread-safe.

### 4. Metrics API

The project exposes:

- `GET /metrics`

This endpoint returns the current metrics snapshot as JSON.

Example response:

```json
{
  "total_events": 1000,
  "event_type_count": {
    "click": 500,
    "view": 300,
    "purchase": 200
  },
  "per_user_count": {
    "1": 100,
    "2": 200
  }
}
```

### 5. Backpressure Handling

The event channel is buffered.

This means:

- short bursts can be absorbed by the buffer
- if workers are slower than the producer, the producer eventually blocks
- the system slows down safely instead of crashing

This is a simple and valid backpressure strategy.

### 6. Graceful Shutdown

When you press `Ctrl+C`:

- the application cancels the shared context
- the producer stops generating events
- the event channel is closed
- workers finish draining pending events
- the HTTP server shuts down cleanly

This prevents abrupt termination and lost in-flight work.

## Key Learning Concepts

This project is useful for learning the following Go concepts:

- goroutines for lightweight concurrency
- channels for communication between producer and workers
- buffered channels for backpressure handling
- worker pool pattern
- context cancellation with `context.Context`
- graceful shutdown using `signal.NotifyContext`
- synchronization with `sync.RWMutex`
- avoiding data races in shared state
- in-memory aggregation of streaming data
- exposing HTTP APIs with `net/http`
- JSON encoding with `encoding/json`
- using `sync.WaitGroup` to wait for goroutines to finish
- simulating load and timing with `time`

## Project Structure

```text
event_processor/
├── main.go
├── model/
│   └── event.go
├── utils/
│   └── utils.go
├── go.mod
└── README.md
```

## Run The Project

From the `event_processor` directory:

```bash
go run .
```

The server will start on:

```text
127.0.0.1:8080
```

Fetch metrics with:

```bash
curl http://127.0.0.1:8080/metrics
```

Stop the program with:

```bash
Ctrl+C
```

## Notes

- Event processing is currently simulated with a sleep to mimic worker load.
- Metrics are stored only in memory, so they reset when the process stops.
- This project is intentionally small so the concurrency patterns are easy to understand.

## Possible Extensions

You can extend this project by adding:

- dropped-event counting
- configurable worker count and buffer size
- request logging
- event validation
- persistent storage
- Prometheus metrics
- unit tests and benchmarks
- separate aggregator goroutine instead of mutex-based aggregation
