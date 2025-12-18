package dto

import (
	"pannypal/internal/common/models"
	"time"
)

type CreateTransactionRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty"`
	Amount      float64 `json:"amount" validate:"required"`
	CategoryID  *int    `json:"category_id" validate:"omitempty"`
	Type        string  `json:"type" validate:"required,oneof=INCOME EXPENSE"`
	Description string  `json:"description" validate:"omitempty"`
}

type UpdateTransactionRequest struct {
	Amount      *float64 `json:"amount" validate:"omitempty,gt=0"`
	CategoryID  *int     `json:"category_id" validate:"omitempty"`
	Type        *string  `json:"type" validate:"omitempty,oneof=INCOME EXPENSE"`
	Description *string  `json:"description" validate:"omitempty"`
}

type GetTransactionsRequest struct {
	PhoneNumber *string    `form:"phone_number,omitempty" validate:"omitempty"`
	Type        *string    `form:"type" validate:"omitempty,oneof=INCOME EXPENSE"`
	CategoryID  *int       `form:"category_id" validate:"omitempty"`
	StartDate   *time.Time `form:"start_date" validate:"omitempty" time_format:"2006-01-02"`
	EndDate     *time.Time `form:"end_date" validate:"omitempty" time_format:"2006-01-02"`
	Page        int        `form:"page" validate:"omitempty,min=1" default:"1"`
	Limit       int        `form:"limit" validate:"omitempty,min=1,max=100" default:"10"`
}

type TransactionSummaryRequest struct {
	PhoneNumber *string    `form:"phone_number,omitempty" validate:"omitempty"`
	StartDate   *time.Time `form:"start_date" validate:"omitempty" time_format:"2006-01-02"`
	EndDate     *time.Time `form:"end_date" validate:"omitempty" time_format:"2006-01-02"`
	Month       *int       `form:"month" validate:"omitempty,min=1,max=12"`
	Year        *int       `form:"year" validate:"omitempty,min=2020"`
}

type TransactionResponse struct {
	ID              uint                   `json:"id"`
	UserID          uint                   `json:"user_id"`
	CategoryID      *uint                  `json:"category_id"`
	Category        models.Category        `json:"category"`
	Amount          float64                `json:"amount"`
	Description     string                 `json:"description"`
	TransactionDate time.Time              `json:"transaction_date"`
	Type            models.TransactionType `json:"type"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Pagination   PaginationResponse    `json:"pagination"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type TransactionSummaryResponse struct {
	TotalIncome      float64           `json:"total_income"`
	TotalExpense     float64           `json:"total_expense"`
	Balance          float64           `json:"balance"`
	TransactionCount int64             `json:"transaction_count"`
	IncomeCount      int64             `json:"income_count"`
	ExpenseCount     int64             `json:"expense_count"`
	CategorySummary  []CategorySummary `json:"category_summary"`
	Period           PeriodInfo        `json:"period"`
}

type CategorySummary struct {
	CategoryID   uint                   `json:"category_id"`
	CategoryName string                 `json:"category_name"`
	Type         models.TransactionType `json:"type"`
	TotalAmount  float64                `json:"total_amount"`
	Count        int64                  `json:"count"`
	Percentage   float64                `json:"percentage"`
}

type PeriodInfo struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Month     *int       `json:"month"`
	Year      *int       `json:"year"`
}
