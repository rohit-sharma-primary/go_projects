# Order Management Backend in Go

This project is a small production-style Go backend for learning how to build an order management system with clean layers, validation, middleware, and tests. It uses the standard library and an in-memory store so you can focus on HTTP APIs and application structure before adding a database.

## What It Covers

- REST APIs for users and orders
- In-memory repositories with auto-increment IDs
- Layered architecture: handler, service, repository, model, middleware
- Request validation and JSON error responses
- Authentication, logging, and panic recovery middleware
- Table-driven tests for services, handlers, and middleware

## Project Structure

```text
cmd/
  server/main.go

internal/
  handler/      HTTP handlers and route registration
  middleware/   auth, logging, panic recovery
  model/        request and response models
  repository/   in-memory data storage
  service/      business rules and validation
  utils/        shared helpers for JSON and validation
```

## API Endpoints

All endpoints require:

```http
Authorization: Bearer afifdosa
```

Routes:

- `POST /users`
- `GET /users`
- `GET /users/{id}`
- `POST /orders`
- `GET /orders`
- `GET /orders/{id}`
- `GET /orders?user_id=1`

## Validation Rules

- Create user:
  `name` must not be empty and `email` must contain `@`
- Create order:
  `user_id` must refer to an existing user and `amount` must be greater than `0`

## Run the Server

```bash
go run ./cmd/server
```

Server address:

```text
http://localhost:3000
```

## Example Requests

Create a user:

```bash
curl -X POST http://localhost:3000/users \
  -H "Authorization: Bearer afifdosa" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

Create an order:

```bash
curl -X POST http://localhost:3000/orders \
  -H "Authorization: Bearer afifdosa" \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"amount":149.99}'
```

List orders for one user:

```bash
curl "http://localhost:3000/orders?user_id=1" \
  -H "Authorization: Bearer afifdosa"
```

## Error Handling

- Invalid input returns `400`
- Missing or invalid auth returns `401`
- Missing resources return `404`
- Panics are recovered and return `500`

## Run Tests

```bash
go test ./...
```

## What to Learn Next

- Add persistent storage with MongoDB or PostgreSQL
- Introduce interfaces where multiple implementations are useful
- Add request-scoped context, structured logging, and configuration
- Add integration tests and API documentation
