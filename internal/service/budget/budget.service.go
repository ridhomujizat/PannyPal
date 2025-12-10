package budget

import (
	"net/http"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/repository/budget"
	"pannypal/internal/service/budget/dto"
	"time"

	"gorm.io/gorm"
)

func (s *Service) CreateBudgetRequest(payload dto.CreateBudgetRequest) *types.Response {
	var user *models.User
	var err error

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err = s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				createdUser, err := s.rp.User.CreateUser(models.User{
					PhoneNumber: *payload.PhoneNumber,
				})
				if err != nil {
					return helper.ParseResponse(&types.Response{
						Code:    http.StatusInternalServerError,
						Message: "Failed to create user",
						Error:   err,
						Data:    nil,
					})
				}
				user = createdUser
			} else {
				return helper.ParseResponse(&types.Response{
					Code:    http.StatusInternalServerError,
					Message: "Database error occurred",
					Data:    nil,
					Error:   err,
				})
			}
		}
	} else {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Phone number is required for creating budget",
			Data:    nil,
		})
	}

	budgetModel := models.Budget{
		UserID:     user.ID,
		CategoryID: uint(payload.CategoryID),
		Amount:     payload.Amount,
		Month:      payload.Month,
		Year:       payload.Year,
	}

	createdBudget, err := s.rp.Budget.CreateBudget(budgetModel)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create budget",
			Data:    nil,
			Error:   err,
		})
	}

	// Get the budget with category preloaded
	budgetWithCategory, err := s.rp.Budget.GetBudgetByID(createdBudget.ID)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get created budget",
			Data:    nil,
			Error:   err,
		})
	}

	response := dto.BudgetResponse{
		ID:         budgetWithCategory.ID,
		UserID:     budgetWithCategory.UserID,
		CategoryID: budgetWithCategory.CategoryID,
		Category:   budgetWithCategory.Category,
		Amount:     budgetWithCategory.Amount,
		Month:      budgetWithCategory.Month,
		Year:       budgetWithCategory.Year,
		CreatedAt:  budgetWithCategory.CreatedAt,
		UpdatedAt:  budgetWithCategory.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusCreated,
		Message: "Budget created successfully",
		Data:    response,
	})
}

func (s *Service) GetBudgetsRequest(payload dto.GetBudgetsRequest) *types.Response {
	var userID *uint

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err := s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   err,
			})
		}
		userID = &user.ID
	}

	filters := budget.BudgetFilters{
		Month: payload.Month,
		Year:  payload.Year,
	}
	if payload.CategoryID != nil {
		categoryID := uint(*payload.CategoryID)
		filters.CategoryID = &categoryID
	}

	budgets, err := s.rp.Budget.GetBudgetsByUserID(userID, filters)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get budgets",
			Data:    nil,
			Error:   err,
		})
	}

	budgetResponses := make([]dto.BudgetResponse, len(budgets))
	for i, b := range budgets {
		budgetResponses[i] = dto.BudgetResponse{
			ID:         b.ID,
			UserID:     b.UserID,
			CategoryID: b.CategoryID,
			Category:   b.Category,
			Amount:     b.Amount,
			Month:      b.Month,
			Year:       b.Year,
			CreatedAt:  b.CreatedAt,
			UpdatedAt:  b.UpdatedAt,
		}
	}

	response := dto.BudgetListResponse{
		Budgets: budgetResponses,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Budgets retrieved successfully",
		Data:    response,
	})
}

func (s *Service) GetBudgetByIDRequest(id uint, phoneNumber string) *types.Response {
	user, err := s.rp.User.GetUserByPhone(phoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   err,
		})
	}

	budgetModel, err := s.rp.Budget.GetBudgetByID(id)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Budget not found",
			Data:    nil,
			Error:   err,
		})
	}

	// Check if budget belongs to user
	if budgetModel.UserID != user.ID {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
	}

	response := dto.BudgetResponse{
		ID:         budgetModel.ID,
		UserID:     budgetModel.UserID,
		CategoryID: budgetModel.CategoryID,
		Category:   budgetModel.Category,
		Amount:     budgetModel.Amount,
		Month:      budgetModel.Month,
		Year:       budgetModel.Year,
		CreatedAt:  budgetModel.CreatedAt,
		UpdatedAt:  budgetModel.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Budget retrieved successfully",
		Data:    response,
	})
}

func (s *Service) UpdateBudgetRequest(id uint, payload dto.UpdateBudgetRequest, phoneNumber string) *types.Response {
	user, err := s.rp.User.GetUserByPhone(phoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   err,
		})
	}

	budgetModel, err := s.rp.Budget.GetBudgetByID(id)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Budget not found",
			Data:    nil,
			Error:   err,
		})
	}

	// Check if budget belongs to user
	if budgetModel.UserID != user.ID {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
	}

	// Update fields if provided
	if payload.CategoryID != nil {
		budgetModel.CategoryID = uint(*payload.CategoryID)
	}
	if payload.Amount != nil {
		budgetModel.Amount = *payload.Amount
	}
	if payload.Month != nil {
		budgetModel.Month = *payload.Month
	}
	if payload.Year != nil {
		budgetModel.Year = *payload.Year
	}

	updatedBudget, err := s.rp.Budget.UpdateBudget(*budgetModel)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update budget",
			Data:    nil,
			Error:   err,
		})
	}

	// Get updated budget with category
	budgetWithCategory, err := s.rp.Budget.GetBudgetByID(updatedBudget.ID)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get updated budget",
			Data:    nil,
			Error:   err,
		})
	}

	response := dto.BudgetResponse{
		ID:         budgetWithCategory.ID,
		UserID:     budgetWithCategory.UserID,
		CategoryID: budgetWithCategory.CategoryID,
		Category:   budgetWithCategory.Category,
		Amount:     budgetWithCategory.Amount,
		Month:      budgetWithCategory.Month,
		Year:       budgetWithCategory.Year,
		CreatedAt:  budgetWithCategory.CreatedAt,
		UpdatedAt:  budgetWithCategory.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Budget updated successfully",
		Data:    response,
	})
}

func (s *Service) DeleteBudgetRequest(id uint, phoneNumber string) *types.Response {
	user, err := s.rp.User.GetUserByPhone(phoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   err,
		})
	}

	budgetModel, err := s.rp.Budget.GetBudgetByID(id)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Budget not found",
			Data:    nil,
			Error:   err,
		})
	}

	// Check if budget belongs to user
	if budgetModel.UserID != user.ID {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
	}

	if err := s.rp.Budget.DeleteBudget(*budgetModel); err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete budget",
			Data:    nil,
			Error:   err,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Budget deleted successfully",
		Data:    nil,
	})
}

func (s *Service) GetBudgetStatusRequest(payload dto.BudgetStatusRequest) *types.Response {
	var userID *uint

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err := s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   err,
			})
		}
		userID = &user.ID
	}

	// Use current month/year if not provided
	now := time.Now()
	month := payload.Month
	year := payload.Year
	if month == nil {
		currentMonth := int(now.Month())
		month = &currentMonth
	}
	if year == nil {
		currentYear := now.Year()
		year = &currentYear
	}

	filters := budget.BudgetStatusFilters{
		Month: month,
		Year:  year,
	}

	statusData, err := s.rp.Budget.GetBudgetStatus(userID, filters)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get budget status",
			Data:    nil,
			Error:   err,
		})
	}

	budgetStatuses := make([]dto.BudgetStatusResponse, len(statusData))
	totalBudget := float64(0)
	totalSpent := float64(0)

	for i, status := range statusData {
		remainingAmount := status.BudgetAmount - status.SpentAmount
		percentageUsed := float64(0)
		if status.BudgetAmount > 0 {
			percentageUsed = (status.SpentAmount / status.BudgetAmount) * 100
		}
		isOverBudget := status.SpentAmount > status.BudgetAmount

		budgetStatuses[i] = dto.BudgetStatusResponse{
			BudgetID:        status.BudgetID,
			CategoryID:      status.CategoryID,
			CategoryName:    status.CategoryName,
			BudgetAmount:    status.BudgetAmount,
			SpentAmount:     status.SpentAmount,
			RemainingAmount: remainingAmount,
			PercentageUsed:  percentageUsed,
			IsOverBudget:    isOverBudget,
			Month:           status.Month,
			Year:            status.Year,
		}

		totalBudget += status.BudgetAmount
		totalSpent += status.SpentAmount
	}

	response := dto.BudgetStatusListResponse{
		BudgetStatuses: budgetStatuses,
		TotalBudget:    totalBudget,
		TotalSpent:     totalSpent,
		TotalRemaining: totalBudget - totalSpent,
		Month:          *month,
		Year:           *year,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Budget status retrieved successfully",
		Data:    response,
	})
}
