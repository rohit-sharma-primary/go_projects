package search

import "testing"

func TestParseQueryDefaultAnd(t *testing.T) {
	query, err := ParseQuery([]string{"golang fast"})
	if err != nil {
		t.Fatalf("ParseQuery returned error: %v", err)
	}

	if query.Condition != "and" {
		t.Fatalf("expected condition and, got %s", query.Condition)
	}

	expected := []string{"golang", "fast"}
	if len(query.Tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(query.Tokens))
	}

	for i := range expected {
		if query.Tokens[i] != expected[i] {
			t.Fatalf("expected token %q at %d, got %q", expected[i], i, query.Tokens[i])
		}
	}
}

func TestParseQueryOr(t *testing.T) {
	query, err := ParseQuery([]string{"golang or python"})
	if err != nil {
		t.Fatalf("ParseQuery returned error: %v", err)
	}

	if query.Condition != "or" {
		t.Fatalf("expected condition or, got %s", query.Condition)
	}
}

func TestParseQueryRejectsMixedOperators(t *testing.T) {
	_, err := ParseQuery([]string{"golang and python or rust"})
	if err == nil {
		t.Fatal("expected error for mixed operators")
	}
}
