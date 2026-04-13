package index

import (
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	dataDir := filepath.Join("..", "data")
	idxMap, err := Build(dataDir)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	postings := idxMap["golang"]
	if len(postings) != 3 {
		t.Fatalf("expected 3 postings for golang, got %d", len(postings))
	}
	if postings[0].File != "doc1.txt" || postings[0].TokenCount != 1 {
		t.Fatalf("unexpected first posting: %+v", postings[0])
	}
}
