package model

type Result struct {
	TotalTransactions int         `json:"total_transactions"`
	TotalAmount       int         `json:"total_amount"`
	AverageAmount     float64     `json:"avg_amount"`
	PerUserTotal      map[int]int `json:"per_user_total"`
}

type Page struct {
	Page       int           `json:"page"`
	TotalPages int           `json:"total_pages"`
	Data       []Transaction `json:"data"`
}

type Transaction struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
}
