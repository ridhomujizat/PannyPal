package dto

import (
	"pannypal/internal/common/models"
	"time"
)

type MonthlyAnalyticsRequest struct {
	PhoneNumber *string `form:"phone_number,omitempty" validate:"omitempty"`
	Year        *int    `form:"year" validate:"omitempty,min=2020"`
}

type YearlyAnalyticsRequest struct {
	PhoneNumber *string `form:"phone_number,omitempty" validate:"omitempty"`
	StartYear   *int    `form:"start_year" validate:"omitempty,min=2020"`
	EndYear     *int    `form:"end_year" validate:"omitempty,min=2020"`
}

type CategoryAnalyticsRequest struct {
	PhoneNumber *string    `form:"phone_number,omitempty" validate:"omitempty"`
	StartDate   *time.Time `form:"start_date" validate:"omitempty"`
	EndDate     *time.Time `form:"end_date" validate:"omitempty"`
	Type        *string    `form:"type" validate:"omitempty,oneof=INCOME EXPENSE"`
}

type MonthlyDataPoint struct {
	Month            int     `json:"month"`
	MonthName        string  `json:"month_name"`
	Year             int     `json:"year"`
	TotalIncome      float64 `json:"total_income"`
	TotalExpense     float64 `json:"total_expense"`
	Balance          float64 `json:"balance"`
	TransactionCount int64   `json:"transaction_count"`
}

type YearlyDataPoint struct {
	Year             int     `json:"year"`
	TotalIncome      float64 `json:"total_income"`
	TotalExpense     float64 `json:"total_expense"`
	Balance          float64 `json:"balance"`
	TransactionCount int64   `json:"transaction_count"`
}

type CategoryDataPoint struct {
	CategoryID    uint                   `json:"category_id"`
	CategoryName  string                 `json:"category_name"`
	Type          models.TransactionType `json:"type"`
	TotalAmount   float64                `json:"total_amount"`
	Count         int64                  `json:"count"`
	Percentage    float64                `json:"percentage"`
	AverageAmount float64                `json:"average_amount"`
}

type MonthlyAnalyticsResponse struct {
	Data         []MonthlyDataPoint `json:"data"`
	Year         int                `json:"year"`
	TotalIncome  float64            `json:"total_income"`
	TotalExpense float64            `json:"total_expense"`
	TotalBalance float64            `json:"total_balance"`
	BestMonth    *MonthlyDataPoint  `json:"best_month"`
	WorstMonth   *MonthlyDataPoint  `json:"worst_month"`
}

type YearlyAnalyticsResponse struct {
	Data         []YearlyDataPoint `json:"data"`
	StartYear    int               `json:"start_year"`
	EndYear      int               `json:"end_year"`
	TotalIncome  float64           `json:"total_income"`
	TotalExpense float64           `json:"total_expense"`
	TotalBalance float64           `json:"total_balance"`
	BestYear     *YearlyDataPoint  `json:"best_year"`
	WorstYear    *YearlyDataPoint  `json:"worst_year"`
}

type CategoryAnalyticsResponse struct {
	Data             []CategoryDataPoint `json:"data"`
	TotalAmount      float64             `json:"total_amount"`
	TransactionCount int64               `json:"transaction_count"`
	Period           PeriodInfo          `json:"period"`
	TopCategory      *CategoryDataPoint  `json:"top_category"`
}

type PeriodInfo struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}
