package index

import (
	"os"
	"strings"
	"unicode"
)

var stopWords = map[string]struct{}{
	"a":    {},
	"an":   {},
	"and":  {},
	"are":  {},
	"as":   {},
	"at":   {},
	"be":   {},
	"but":  {},
	"by":   {},
	"for":  {},
	"in":   {},
	"is":   {},
	"it":   {},
	"of":   {},
	"on":   {},
	"or":   {},
	"the":  {},
	"to":   {},
	"was":  {},
	"were": {},
	"with": {},
}

func NormalizeToken(token string) string {
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, token)
	return normalized
}

func ShouldSkipToken(token string) bool {
	if token == "" {
		return true
	}
	_, found := stopWords[token]
	return found
}

func Tokenize(file string) (map[string]int, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	indexMap := make(map[string]int)
	for _, rawToken := range strings.Fields(string(content)) {
		token := NormalizeToken(rawToken)
		if ShouldSkipToken(token) {
			continue
		}
		indexMap[token]++
	}
	return indexMap, nil
}
