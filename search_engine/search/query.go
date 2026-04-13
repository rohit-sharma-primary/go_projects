package search

import (
	"fmt"
	"strings"

	"search_engine/index"
)

type ParsedQuery struct {
	Raw       string
	Tokens    []string
	Condition string
}

func ParseQuery(args []string) (ParsedQuery, error) {
	raw := strings.TrimSpace(strings.Join(args, " "))
	if raw == "" {
		return ParsedQuery{}, fmt.Errorf("empty query")
	}

	condition := "and"
	conditionExplicit := false
	seen := make(map[string]struct{})
	cleanSlice := make([]string, 0)

	for _, arg := range strings.Fields(raw) {
		token := index.NormalizeToken(arg)
		if token == "" {
			continue
		}

		if token == "or" || token == "and" {
			if conditionExplicit && condition != token {
				return ParsedQuery{}, fmt.Errorf("mixed query operators are not supported")
			}
			condition = token
			conditionExplicit = true
			continue
		}

		if index.ShouldSkipToken(token) {
			continue
		}
		if _, found := seen[token]; found {
			continue
		}
		seen[token] = struct{}{}
		cleanSlice = append(cleanSlice, token)
	}

	if len(cleanSlice) == 0 {
		if conditionExplicit {
			return ParsedQuery{}, fmt.Errorf("query must include at least one searchable term")
		} else {
			return ParsedQuery{}, fmt.Errorf("query did not contain searchable terms")
		}
	}

	return ParsedQuery{
		Raw:       raw,
		Tokens:    cleanSlice,
		Condition: condition,
	}, nil
}
