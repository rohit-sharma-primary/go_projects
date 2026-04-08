package main

import (
	"csv_to_json_transformer/model"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func resolveFilePath(file string) (string, error) {
	if _, err := os.Stat(file); err == nil {
		return file, nil
	}
	dataPath := "data/" + file
	if _, err := os.Stat(dataPath); err == nil {
		return dataPath, nil
	}
	return "", fmt.Errorf("File not found")
}

func getTransaction(record []string) (model.Transaction, error) {
	if len(record) != 5 {
		return model.Transaction{}, errors.New("Invalid record")
	}

	userId, uIderr := strconv.Atoi(strings.TrimSpace(record[0]))
	amount, amterr := strconv.Atoi(strings.TrimSpace(record[1]))

	if uIderr != nil || amterr != nil {
		return model.Transaction{}, errors.New("Invalid record")
	}
	ts, terr := time.Parse(time.RFC3339, strings.TrimSpace(record[4]))
	if terr != nil {
		return model.Transaction{}, errors.New("Invalid record. Error parsing time")
	}
	return model.Transaction{
		UserID:    userId,
		Amount:    amount,
		Status:    strings.TrimSpace(record[2]),
		Category:  strings.TrimSpace(record[3]),
		Timestamp: ts,
	}, nil
}

func updateResult(txn model.Transaction, result map[int]model.UserSummary) {
	if item, found := result[txn.UserID]; found {
		item.TotalAmount += txn.Amount
		item.CategorySpend[txn.Category] += txn.Amount
		item.TransactionCount++
		result[txn.UserID] = item
	} else {
		categorySpent := make(map[string]int)
		categorySpent[txn.Category] = txn.Amount
		result[txn.UserID] = model.UserSummary{
			UserID:           txn.UserID,
			TotalAmount:      txn.Amount,
			CategorySpend:    categorySpent,
			TransactionCount: 1,
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Provide a file name!")
	}

	fileName := os.Args[1]

	filePath, err := resolveFilePath(fileName)
	if err != nil {
		log.Fatalf("File doesn't exist: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("File is empty!")
			return
		}
		log.Fatal("Error reading csv header!")
	}
	errCountRows := 0
	result := make(map[int]model.UserSummary)
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("Error reading csv body")
		}
		// fmt.Println("Length of record: ", len(record))
		txn, err := getTransaction(record)
		if err != nil {
			errCountRows++
			continue
		}
		if txn.Status != "completed" {
			continue
		}
		updateResult(txn, result)
	}
	if errCountRows > 0 {
		fmt.Println("No of invalid rows : ", errCountRows)
	}

	var finalResult []model.UserSummary
	for _, v := range result {
		if v.TotalAmount >= 200 {
			finalResult = append(finalResult, v)
		}
	}

	// sort for deterministic output
	slices.SortFunc(finalResult, func(a, b model.UserSummary) int {
		return a.UserID - b.UserID
	})

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(finalResult)
}
