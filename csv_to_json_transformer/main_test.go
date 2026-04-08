package main

import (
	"csv_to_json_transformer/model"
	"testing"
	"time"
)

// ─── getTransaction tests ───

func TestGetTransaction_ValidRecord(t *testing.T) {
	record := []string{"1", "100", "completed", "food", "2026-01-01T10:00:00Z"}
	txn, err := getTransaction(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txn.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", txn.UserID)
	}
	if txn.Amount != 100 {
		t.Errorf("expected Amount=100, got %d", txn.Amount)
	}
	if txn.Status != "completed" {
		t.Errorf("expected Status=completed, got %s", txn.Status)
	}
	if txn.Category != "food" {
		t.Errorf("expected Category=food, got %s", txn.Category)
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2026-01-01T10:00:00Z")
	if !txn.Timestamp.Equal(expectedTime) {
		t.Errorf("expected Timestamp=%v, got %v", expectedTime, txn.Timestamp)
	}
}

func TestGetTransaction_TrimsWhitespace(t *testing.T) {
	record := []string{" 1 ", " 100 ", " completed ", " food ", " 2026-01-01T10:00:00Z "}
	txn, err := getTransaction(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txn.UserID != 1 {
		t.Errorf("expected UserID=1, got %d", txn.UserID)
	}
	if txn.Status != "completed" {
		t.Errorf("expected Status=completed, got %s", txn.Status)
	}
	if txn.Category != "food" {
		t.Errorf("expected Category=food, got %s", txn.Category)
	}
}

func TestGetTransaction_TooFewFields(t *testing.T) {
	record := []string{"1", "100", "completed", "food"}
	_, err := getTransaction(record)
	if err == nil {
		t.Fatal("expected error for record with 4 fields")
	}
}

func TestGetTransaction_TooManyFields(t *testing.T) {
	record := []string{"1", "100", "completed", "food", "2026-01-01T10:00:00Z", "extra"}
	_, err := getTransaction(record)
	if err == nil {
		t.Fatal("expected error for record with 6 fields")
	}
}

func TestGetTransaction_InvalidUserID(t *testing.T) {
	record := []string{"abc", "100", "completed", "food", "2026-01-01T10:00:00Z"}
	_, err := getTransaction(record)
	if err == nil {
		t.Fatal("expected error for non-numeric user_id")
	}
}

func TestGetTransaction_InvalidAmount(t *testing.T) {
	record := []string{"1", "xyz", "completed", "food", "2026-01-01T10:00:00Z"}
	_, err := getTransaction(record)
	if err == nil {
		t.Fatal("expected error for non-numeric amount")
	}
}

func TestGetTransaction_InvalidTimestamp(t *testing.T) {
	record := []string{"1", "100", "completed", "food", "not-a-timestamp"}
	_, err := getTransaction(record)
	if err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}

func TestGetTransaction_EmptyRecord(t *testing.T) {
	record := []string{}
	_, err := getTransaction(record)
	if err == nil {
		t.Fatal("expected error for empty record")
	}
}

// ─── updateResult tests ───

func TestUpdateResult_NewUser(t *testing.T) {
	result := make(map[int]model.UserSummary)
	txn := model.Transaction{UserID: 1, Amount: 100, Status: "completed", Category: "food"}

	updateResult(txn, result)

	summary, found := result[1]
	if !found {
		t.Fatal("expected user 1 in result")
	}
	if summary.TotalAmount != 100 {
		t.Errorf("expected TotalAmount=100, got %d", summary.TotalAmount)
	}
	if summary.TransactionCount != 1 {
		t.Errorf("expected TransactionCount=1, got %d", summary.TransactionCount)
	}
	if summary.CategorySpend["food"] != 100 {
		t.Errorf("expected food=100, got %d", summary.CategorySpend["food"])
	}
}

func TestUpdateResult_ExistingUser_SameCategory(t *testing.T) {
	result := make(map[int]model.UserSummary)
	txn1 := model.Transaction{UserID: 2, Amount: 200, Status: "completed", Category: "shopping"}
	txn2 := model.Transaction{UserID: 2, Amount: 300, Status: "completed", Category: "shopping"}

	updateResult(txn1, result)
	updateResult(txn2, result)

	summary := result[2]
	if summary.TotalAmount != 500 {
		t.Errorf("expected TotalAmount=500, got %d", summary.TotalAmount)
	}
	if summary.TransactionCount != 2 {
		t.Errorf("expected TransactionCount=2, got %d", summary.TransactionCount)
	}
	if summary.CategorySpend["shopping"] != 500 {
		t.Errorf("expected shopping=500, got %d", summary.CategorySpend["shopping"])
	}
}

func TestUpdateResult_ExistingUser_DifferentCategories(t *testing.T) {
	result := make(map[int]model.UserSummary)
	txn1 := model.Transaction{UserID: 1, Amount: 100, Status: "completed", Category: "food"}
	txn2 := model.Transaction{UserID: 1, Amount: 150, Status: "completed", Category: "travel"}

	updateResult(txn1, result)
	updateResult(txn2, result)

	summary := result[1]
	if summary.TotalAmount != 250 {
		t.Errorf("expected TotalAmount=250, got %d", summary.TotalAmount)
	}
	if summary.TransactionCount != 2 {
		t.Errorf("expected TransactionCount=2, got %d", summary.TransactionCount)
	}
	if summary.CategorySpend["food"] != 100 {
		t.Errorf("expected food=100, got %d", summary.CategorySpend["food"])
	}
	if summary.CategorySpend["travel"] != 150 {
		t.Errorf("expected travel=150, got %d", summary.CategorySpend["travel"])
	}
}

func TestUpdateResult_MultipleUsers(t *testing.T) {
	result := make(map[int]model.UserSummary)
	txn1 := model.Transaction{UserID: 1, Amount: 100, Status: "completed", Category: "food"}
	txn2 := model.Transaction{UserID: 2, Amount: 200, Status: "completed", Category: "shopping"}

	updateResult(txn1, result)
	updateResult(txn2, result)

	if len(result) != 2 {
		t.Errorf("expected 2 users, got %d", len(result))
	}
	if result[1].TotalAmount != 100 {
		t.Errorf("expected user 1 TotalAmount=100, got %d", result[1].TotalAmount)
	}
	if result[2].TotalAmount != 200 {
		t.Errorf("expected user 2 TotalAmount=200, got %d", result[2].TotalAmount)
	}
}

// ─── resolveFilePath tests ───

func TestResolveFilePath_DirectFile(t *testing.T) {
	path, err := resolveFilePath("main.go")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "main.go" {
		t.Errorf("expected main.go, got %s", path)
	}
}

func TestResolveFilePath_DataDirectory(t *testing.T) {
	path, err := resolveFilePath("transactions.csv")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "data/transactions.csv" {
		t.Errorf("expected data/transactions.csv, got %s", path)
	}
}

func TestResolveFilePath_NotFound(t *testing.T) {
	_, err := resolveFilePath("nonexistent_file.xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
