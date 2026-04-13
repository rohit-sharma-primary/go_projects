# Mini Search Engine in Go

This project is a small document search engine written in Go. It builds an inverted index from text files in the `data/` directory, supports fast term lookups, and ranks matching documents using basic term frequency (TF).

## Features

- Indexes plain text documents from `data/`
- Builds an inverted index for fast token-to-document lookup
- Supports `AND` queries by default
- Supports explicit `OR` queries
- Ranks results using summed term frequency across matched query terms
- Normalizes tokens by lowercasing and stripping punctuation
- Ignores common stop words
- Includes unit tests for indexing, parsing, and ranking

## Project Structure

```text
search_engine/
├── data/
│   ├── doc1.txt
│   ├── doc2.txt
│   └── doc3.txt
├── index/
│   ├── index.go
│   └── tokenizer.go
├── model/
│   └── model.go
├── search/
│   ├── query.go
│   └── rank.go
├── main.go
└── README.md
```

## How It Works

### 1. Tokenization

Each document is read from the `data/` directory and split into tokens.

The tokenizer:

- converts text to lowercase
- removes punctuation
- skips common stop words such as `the`, `is`, and `and`
- counts how many times each token appears in a document

Example:

```text
doc1.txt: golang is fast and efficient
```

Produces:

```text
golang -> 1
fast -> 1
efficient -> 1
```

### 2. Inverted Index

The index maps each token to a postings list.

Example:

```text
golang -> [(doc1.txt, 1), (doc2.txt, 1), (doc3.txt, 1)]
python -> [(doc3.txt, 1)]
fast   -> [(doc1.txt, 1)]
```

This allows direct lookup of all documents containing a token without scanning every file during query time.

### 3. Query Parsing

The CLI accepts a query string:

```bash
go run main.go "golang fast"
go run main.go "golang or python"
```

Supported query modes:

- `AND` query, default behavior
  - `"golang fast"`
  - returns only documents containing both `golang` and `fast`
- `OR` query
  - `"golang or python"`
  - returns documents containing either token

### 4. Ranking

Documents are ranked using basic term frequency.

Scoring rule:

- for each matching query term, add that term's frequency in the document
- sort by higher total score first
- break ties by filename for stable output

Example:

If the query is:

```text
golang
```

and the frequencies are:

```text
doc3.txt -> golang appears 2 times
doc1.txt -> golang appears 1 time
doc2.txt -> golang appears 1 time
```

then the output is:

```text
1. doc3.txt (score: 2)
2. doc1.txt (score: 1)
3. doc2.txt (score: 1)
```

## Running the Project

From the `search_engine` directory:

```bash
go run main.go "golang fast"
```

Example output:

```text
Query: golang fast

Results:
1. doc1.txt (score: 2)
```

Run an OR query:

```bash
go run main.go "golang or python"
```

Example output:

```text
Query: golang or python

Results:
1. doc3.txt (score: 2)
2. doc1.txt (score: 1)
3. doc2.txt (score: 1)
```

## Running Tests

```bash
go test ./...
```

## Complexity

### Index Build Time

- proportional to the total number of tokens across all documents

### Query Time

- token lookup in the inverted index is average-case near `O(1)` per token
- total query time depends on the size of the postings lists scanned for the query terms

In practice, query execution is much faster than scanning every document on each search.

## Current Limitations

- no phrase queries
- no wildcard or fuzzy matching
- no TF-IDF or BM25 scoring
- index is built in memory on each run
- no stemming or lemmatization

## Possible Improvements

- persist the index to disk
- support phrase search and quoted queries
- add TF-IDF or BM25 ranking
- add stemming for related word forms
- expose the engine through an HTTP API
- support incremental indexing when documents change

## Summary

This project demonstrates the core mechanics of a search engine:

- document ingestion
- tokenization
- inverted indexing
- boolean retrieval
- TF-based ranking

It is small enough to understand quickly, while still using the same basic ideas as larger production search systems.
