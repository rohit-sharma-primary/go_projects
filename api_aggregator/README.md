# API Aggregator

A Go service that aggregates data from multiple backend APIs concurrently and generates consolidated reports. It fetches user and order data in parallel, then produces a per-user spending report.

## Architecture

```
┌────────────┐         ┌──────────────┐
│  Aggregator │ ──────▶ │   Backend    │
│  (main.go)  │         │  (port 8081) │
│  port 8080  │         │  /users      │
│  /report    │         │  /orders     │
└────────────┘         └──────────────┘
```

- **Aggregator** (`main.go`) — Front-facing HTTP server on port `8080`. The `/report` endpoint fetches users and orders concurrently from the backend, then generates a spending report.
- **Backend** (`cmd/server/backend.go`) — Stub API on port `8081` serving `/users` and `/orders` with sample JSON data.

## Project Structure

```
api_aggregator/
├── main.go                  # Aggregator server (port 8080)
├── cmd/server/
│   └── backend.go           # Backend stub server (port 8081)
├── internal/
│   └── model/
│       └── model.go         # Data models (User, Order, Report)
├── pkg/
│   └── utils.go             # Report generation logic
└── go.mod
```

## Models

| Struct   | Fields |
|----------|--------|
| `User`   | `id`, `name` |
| `Order`  | `id`, `user_id`, `status`, `amount` |
| `Report` | `user_id`, `name`, `total_spend`, `completed_orders` |

## Getting Started

### Prerequisites

- Go 1.25+

### Run the Backend

```bash
go run cmd/server/backend.go
```

This starts the stub API on `http://localhost:8081`.

### Run the Aggregator

In a separate terminal:

```bash
go run main.go
```

This starts the aggregator on `http://localhost:8080`.

### Endpoints

| Server     | Endpoint  | Description |
|------------|-----------|-------------|
| Aggregator | `GET /`       | Welcome message |
| Aggregator | `GET /report` | Aggregated user spending report |
| Backend    | `GET /users`  | List of users (JSON) |
| Backend    | `GET /orders` | List of orders (JSON) |

### Example

```bash
curl http://localhost:8080/report
```

Response:

```json
[
  { "user_id": 1, "name": "Alice", "total_spend": 250.75, "completed_orders": 0 },
  { "user_id": 2, "name": "Bob", "total_spend": 300.50, "completed_orders": 0 },
  { "user_id": 3, "name": "Charlie", "total_spend": 150.25, "completed_orders": 0 }
]
```

## Key Concepts

- **Concurrent HTTP calls** — Users and orders are fetched in parallel using goroutines and channels.
- **Context with timeout** — A 2-second timeout is applied to backend requests via `context.WithTimeout`.
- **Error propagation** — Errors from either goroutine are surfaced through a buffered error channel.
