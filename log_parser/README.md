# Log Parser

A Go CLI tool that parses structured log files, extracts action frequencies, and counts errors per user.

## Features

- **Line-by-line parsing** — efficiently reads large log files using buffered scanning
- **Action frequency analysis** — counts occurrences of each action and reports the most frequent one
- **Error tracking per user** — tallies errors grouped by `user_id`
- **Malformed line detection** — identifies and counts lines that don't match the expected format

## Log Format

The parser expects space-delimited log lines with at least 4 tokens:

```
<timestamp> <level> user_id=<id> action=<name>
```

Example:

```
2026-01-01T10:00:00Z INFO user_id=1 action=login
2026-01-01T10:01:00Z ERROR user_id=2 action=payment
```

## Getting Started

### Prerequisites

- Go 1.25+

### Run

```bash
go run main.go
```

By default, the tool parses `files/logs1.txt`. To change the input file, edit the `file` variable in `main()`.

### Build

```bash
go build -o log_parser
./log_parser
```

## Output

The tool prints:

1. Number of incorrectly formatted lines
2. The most frequent action and its count
3. Error counts grouped by user ID

Example:

```
No of incorrect lines =  2
Top action = login with count 5
Errors per user = map[2:1 3:1 5:1]
```

## Project Structure

```
├── main.go          # Entry point and all parsing logic
├── go.mod           # Go module definition
├── logs.txt         # Sample output file (created by createFile)
└── files/
    └── logs1.txt    # Sample log input file
```
