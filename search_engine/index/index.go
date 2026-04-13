package index

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"search_engine/model"
)

func Build(dataDir string) (map[string][]model.Pair, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("read data directory: %w", err)
	}

	idxMap := make(map[string][]model.Pair)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dataDir, entry.Name())
		fileIdx, err := Tokenize(filePath)
		if err != nil {
			return nil, fmt.Errorf("tokenize %s: %w", entry.Name(), err)
		}

		for token, count := range fileIdx {
			idxMap[token] = append(idxMap[token], model.Pair{
				File:       entry.Name(),
				TokenCount: count,
			})
		}
	}

	for token := range idxMap {
		sort.Slice(idxMap[token], func(i, j int) bool {
			if idxMap[token][i].TokenCount == idxMap[token][j].TokenCount {
				return idxMap[token][i].File < idxMap[token][j].File
			}
			return idxMap[token][i].TokenCount > idxMap[token][j].TokenCount
		})
	}

	return idxMap, nil
}
