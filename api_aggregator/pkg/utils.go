package utils

import (
	"api_aggregator/internal/model"
	"fmt"
)

func GenerateReport(users []model.User, orders []model.Order) []model.Report {
	mapReports := make(map[int]*model.Report)
	for _, user := range users {
		mapReports[user.ID] = &model.Report{
			UserID: user.ID,
			Name:   user.Name,
		}
	}
	for _, order := range orders {
		if order.Status != "completed" {
			continue
		}
		if report, exists := mapReports[order.UserID]; exists {
			report.TotalSpend += order.Amount
		} else {
			fmt.Printf("Warning: Order with user_id %d has no corresponding user\n", order.UserID)
		}
	}
	var finalReport []model.Report
	for _, report := range mapReports {
		finalReport = append(finalReport, *report)
	}
	return finalReport
}
