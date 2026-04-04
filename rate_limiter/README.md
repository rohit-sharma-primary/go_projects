# Rate Limiter Learning Project

This project implements a rate-limited API client in Go along with a mock backend service.

The goal is to simulate a real integration where the backend is unreliable and rate-limited, while the client still fetches all pages safely and produces an aggregated summary.

## Problem Covered

The mock backend exposes:

- `GET /transactions?page=<n>`
- Random latency between `100ms` and `1000ms`
- Random `500` failures roughly `20%` of the time
- `429 Too Many Requests` when more than `3` requests hit the backend within `1` second

The client exposes:

- `GET /summary`

The client:

- Starts from page `1`
- Discovers `total_pages`
- Fetches the remaining pages concurrently
- Limits outbound traffic to `3 requests/second`
- Retries on `500`, `429`, and timeout failures
- Uses exponential backoff: `100ms`, `200ms`, `400ms`, `800ms`, `1600ms`
- Uses a `2 second` timeout per request attempt
- Aggregates the fetched transactions into a summary response

## Project Structure

- [main.go](/Users/rohitsharma/Documents/personal/learning/go&mongo/rate_limiter/main.go): client service and `/summary` API
- [server/backend.go](/Users/rohitsharma/Documents/personal/learning/go&mongo/rate_limiter/server/backend.go): mock backend and `/transactions` API
- [model/model.go](/Users/rohitsharma/Documents/personal/learning/go&mongo/rate_limiter/model/model.go): shared response models
- [data/page1.json](/Users/rohitsharma/Documents/personal/learning/go&mongo/rate_limiter/data/page1.json): sample paginated data
- [data/page2.json](/Users/rohitsharma/Documents/personal/learning/go&mongo/rate_limiter/data/page2.json): sample paginated data
- [data/page3.json](/Users/rohitsharma/Documents/personal/learning/go&mongo/rate_limiter/data/page3.json): sample paginated data

## Requirements

- Go `1.25.8` or compatible Go `1.25.x`

## How To Run

Start the backend from the repository root:

```bash
go run server/backend.go
```

In another terminal, start the client from the repository root:

```bash
go run main.go
```

Then call the summary API:

```bash
curl http://127.0.0.1:8080/summary
```

## Example Response

With the bundled JSON files, a successful `/summary` response looks like:

```json
{
  "total_transactions": 14,
  "total_amount": 2750,
  "avg_amount": 196.42857142857142,
  "per_user_total": {
    "1": 750,
    "2": 1100,
    "3": 900
  }
}
```

Because the backend intentionally injects latency and failures, repeated runs may take different amounts of time, but the final aggregation should remain the same when all pages are fetched successfully.

## Implementation Notes

- The client uses a fixed worker pool to avoid uncontrolled goroutines.
- Rate limiting is handled separately from concurrency.
- Each retry attempt gets its own `2s` timeout.
- The backend rate limit is implemented as a simple per-second window counter.
- The backend reads page data from the `data/` directory, so run commands from the project root.
