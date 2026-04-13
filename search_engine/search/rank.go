package search

import (
	"sort"

	"search_engine/model"
)

func FilterAndRank(tokens []string, condition string, idxMap map[string][]model.Pair) []model.Result {
	if len(tokens) == 0 {
		return nil
	}

	scores := make(map[string]int)
	matches := make(map[string]int)

	for _, token := range tokens {
		postings, found := idxMap[token]
		if !found {
			if condition == "and" {
				return nil
			}
			continue
		}

		seenInToken := make(map[string]struct{}, len(postings))
		for _, posting := range postings {
			scores[posting.File] += posting.TokenCount
			if _, alreadyCounted := seenInToken[posting.File]; !alreadyCounted {
				matches[posting.File]++
				seenInToken[posting.File] = struct{}{}
			}
		}
	}

	results := make([]model.Result, 0, len(scores))
	for document, score := range scores {
		if condition == "and" && matches[document] != len(tokens) {
			continue
		}
		results = append(results, model.Result{
			Document: document,
			Score:    score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Document < results[j].Document
		}
		return results[i].Score > results[j].Score
	})

	return results
}
