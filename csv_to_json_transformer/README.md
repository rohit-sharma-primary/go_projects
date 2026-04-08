# CSV to JSON Transformer

A Go CLI tool that reads a CSV file of transactions, aggregates data per user, applies business rules, and outputs transformed JSON.

## Features

- Parses CSV transaction records into structured data
- Skips invalid rows and logs the count
- Filters out non-completed transactions
- Aggregates per user: total amount, category-wise spend, transaction count
- Applies threshold filter (only users with `total_amount >= 200`)
- Deterministic JSON output sorted by `user_id`

## Usage

```bash
go run main.go <filename>
```

The tool looks for the file in the current directory first, then falls back to `data/`.

```bash
go run main.go transactions.csv
```

## Input Format

CSV with the following columns:

| Column    | Type   | Description                     |
|-----------|--------|---------------------------------|
| user_id   | int    | User identifier                 |
| amount    | int    | Transaction amount              |
| status    | string | `completed` or `failed`         |
| category  | string | Spending category               |
| timestamp | string | RFC 3339 timestamp              |

Example:

```csv
user_id,amount,status,category,timestamp
1,100,completed,food,2026-01-01T10:00:00Z
1,50,failed,food,2026-01-01T11:00:00Z
2,200,completed,shopping,2026-01-02T10:00:00Z
2,300,completed,shopping,2026-01-02T12:00:00Z
3,150,completed,travel,2026-01-03T10:00:00Z
```

## Output Format

JSON array of user summaries, filtered to `total_amount >= 200`:

```json
[
  {
    "user_id": 2,
    "total_amount": 500,
    "category_spend": {
      "shopping": 500
    },
    "transaction_count": 2
  }
]
```

## Project Structure

```
├── main.go                 # CLI entry point, CSV parsing, aggregation logic
├── model/
│   └── model.go            # Transaction and UserSummary structs
├── data/
│   └── transactions.csv    # Sample input data
├── go.mod
└── README.md
```

## Business Rules

1. Rows that fail validation (wrong field count, unparseable numbers/timestamps) are skipped and counted.
2. Only transactions with `status = "completed"` are included.
3. Users are aggregated by `user_id` with total amount, per-category spend, and transaction count.
4. Only users with `total_amount >= 200` appear in the final output.
