package model

import "time"

type Transaction struct {
	UserID    int
	Amount    int
	Status    string
	Category  string
	Timestamp time.Time
}

type UserSummary struct {
	UserID           int            `json:"user_id"`
	TotalAmount      int            `json:"total_amount"`
	CategorySpend    map[string]int `json:"category_spend"`
	TransactionCount int            `json:"transaction_count"`
}
