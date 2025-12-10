package dto

import (
	"pannypal/internal/common/models"
	"time"
)

type CreateBudgetRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty"`
	CategoryID  int     `json:"category_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Month       int     `json:"month" validate:"required,min=1,max=12"`
	Year        int     `json:"year" validate:"required,min=2020"`
}

type UpdateBudgetRequest struct {
	CategoryID *int     `json:"category_id" validate:"omitempty"`
	Amount     *float64 `json:"amount" validate:"omitempty,gt=0"`
	Month      *int     `json:"month" validate:"omitempty,min=1,max=12"`
	Year       *int     `json:"year" validate:"omitempty,min=2020"`
}

type GetBudgetsRequest struct {
	PhoneNumber *string `form:"phone_number,omitempty" validate:"omitempty"`
	Month       *int    `form:"month" validate:"omitempty,min=1,max=12"`
	Year        *int    `form:"year" validate:"omitempty,min=2020"`
	CategoryID  *int    `form:"category_id" validate:"omitempty"`
}

type BudgetStatusRequest struct {
	PhoneNumber *string `form:"phone_number,omitempty" validate:"omitempty"`
	Month       *int    `form:"month" validate:"omitempty,min=1,max=12"`
	Year        *int    `form:"year" validate:"omitempty,min=2020"`
}

type BudgetResponse struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"user_id"`
	CategoryID uint            `json:"category_id"`
	Category   models.Category `json:"category"`
	Amount     float64         `json:"amount"`
	Month      int             `json:"month"`
	Year       int             `json:"year"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type BudgetListResponse struct {
	Budgets []BudgetResponse `json:"budgets"`
}

type BudgetStatusResponse struct {
	BudgetID        uint    `json:"budget_id"`
	CategoryID      uint    `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	BudgetAmount    float64 `json:"budget_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	PercentageUsed  float64 `json:"percentage_used"`
	IsOverBudget    bool    `json:"is_over_budget"`
	Month           int     `json:"month"`
	Year            int     `json:"year"`
}

type BudgetStatusListResponse struct {
	BudgetStatuses []BudgetStatusResponse `json:"budget_statuses"`
	TotalBudget    float64                `json:"total_budget"`
	TotalSpent     float64                `json:"total_spent"`
	TotalRemaining float64                `json:"total_remaining"`
	Month          int                    `json:"month"`
	Year           int                    `json:"year"`
}
