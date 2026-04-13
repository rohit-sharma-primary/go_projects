package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"search_engine/index"
	"search_engine/model"
	"search_engine/search"
)

func buildIndex() map[string][]model.Pair {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dataDir := filepath.Join(workingDir, "data")
	idxMap, err := index.Build(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	return idxMap
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No query provided")
	}

	query, err := search.ParseQuery(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	idxMap := buildIndex()
	results := search.FilterAndRank(query.Tokens, query.Condition, idxMap)

	fmt.Printf("Query: %s\n\n", query.Raw)
	if len(results) == 0 {
		fmt.Println("Results:")
		fmt.Println("No matching documents found.")
		return
	}

	fmt.Println("Results:")
	for i, result := range results {
		fmt.Printf("%d. %s (score: %d)\n", i+1, result.Document, result.Score)
	}
}
