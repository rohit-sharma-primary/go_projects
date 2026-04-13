package search

import (
	"testing"

	"search_engine/model"
)

func TestFilterAndRankAndQuery(t *testing.T) {
	idxMap := map[string][]model.Pair{
		"golang": {
			{File: "doc1.txt", TokenCount: 2},
			{File: "doc2.txt", TokenCount: 1},
		},
		"fast": {
			{File: "doc1.txt", TokenCount: 1},
			{File: "doc3.txt", TokenCount: 3},
		},
	}

	results := FilterAndRank([]string{"golang", "fast"}, "and", idxMap)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Document != "doc1.txt" || results[0].Score != 3 {
		t.Fatalf("unexpected result: %+v", results[0])
	}
}

func TestFilterAndRankOrQuery(t *testing.T) {
	idxMap := map[string][]model.Pair{
		"golang": {
			{File: "doc3.txt", TokenCount: 2},
			{File: "doc1.txt", TokenCount: 1},
			{File: "doc2.txt", TokenCount: 1},
		},
	}

	results := FilterAndRank([]string{"golang"}, "or", idxMap)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if results[0].Document != "doc3.txt" || results[0].Score != 2 {
		t.Fatalf("unexpected top result: %+v", results[0])
	}
}
