package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const TOKEN_SIZE = 4

func createFile() {
	f, err := os.Create("logs.txt")
	if err != nil {
		panic(err)
	}

	length, err := f.WriteString("Just some random text")
	if err != nil {
		panic(err)
	}
	fmt.Printf("File created: %s and Length written : %d", f, length)
}

func readFile(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
		return
	}
	fmt.Printf("File length = %d \n\n\n", len(f))
	// fmt.Printf("File content = %s \n", f)
}

func readFileLineByLine(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	actionCount := make(map[string]int)
	errCount := make(map[int]int)
	incorrectLines := 0
	for scanner.Scan() {
		line := scanner.Text()
		if correct := parseText(line, actionCount, errCount); !correct {
			incorrectLines++
		}
	}
	fmt.Println("No of incorrect lines = ", incorrectLines)
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	maxFreq := 0
	result := ""
	for action, freq := range actionCount {
		if freq > maxFreq {
			maxFreq = freq
			result = action
		}
	}
	fmt.Printf("Top action = %s with count %d\n", result, maxFreq)
	fmt.Printf("Errors per user = %v", errCount)
}

func parseText(text string, actionCount map[string]int, errCount map[int]int) bool {
	// fmt.Println(text)
	tokens := strings.Fields(text)
	if len(tokens) < TOKEN_SIZE {
		return false
	}
	if strings.EqualFold(tokens[1], "error") {
		if !strings.HasPrefix(tokens[2], "user_id=") {
			return true
		}
		_, userId, _ := strings.Cut(tokens[2], "=")
		user_id, err := strconv.Atoi(userId)
		if err != nil {
			return true
		}

		if _, found := errCount[user_id]; !found {
			errCount[user_id] = 0
		}
		errCount[user_id]++
	}
	// Action parsing
	first, second, found := strings.Cut(tokens[3], "=")
	if !found || first != "action" {
		return true
	}
	actionCount[second]++
	return true
}

func main() {
	// file, err := createFile()
	file := "files/logs1.txt"
	readFileLineByLine(file)
}
