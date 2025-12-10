package transaction

import (
	"net/http"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/repository/transaction"
	"pannypal/internal/service/transaction/dto"
	"time"

	"gorm.io/gorm"
)

func (s *Service) CreateTransactionRequest(payload dto.CreateTransactionRequest) *types.Response {
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
				// Handle other database errors
				return helper.ParseResponse(&types.Response{
					Code:    http.StatusInternalServerError,
					Message: "Database error occurred",
					Data:    nil,
					Error:   err,
				})
			}
		}
	} else {
		// If no phone number provided, create transaction without user
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Phone number is required for creating transaction",
			Data:    nil,
		})
	}

	transaction := models.Transaction{
		UserID:          user.ID,
		Amount:          payload.Amount,
		TransactionDate: time.Now(),
		Type:            models.TransactionType(payload.Type),
	}
	if payload.CategoryID != nil {
		categoryID := uint(*payload.CategoryID)
		transaction.CategoryID = &categoryID
	}
	createdTransaction, err := s.rp.Transaction.CreateTransaction(transaction)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create transaction",
			Data:    nil,
			Error:   err,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusCreated,
		Message: "Transaction created successfully",
		Data:    createdTransaction,
	})
}

func (s *Service) GetTransactionsRequest(payload dto.GetTransactionsRequest) *types.Response {
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

	// Set default values
	if payload.Page == 0 {
		payload.Page = 1
	}
	if payload.Limit == 0 {
		payload.Limit = 10
	}

	filters := transaction.TransactionFilters{
		Page:  payload.Page,
		Limit: payload.Limit,
	}

	if payload.Type != nil {
		transactionType := models.TransactionType(*payload.Type)
		filters.Type = &transactionType
	}
	if payload.CategoryID != nil {
		categoryID := uint(*payload.CategoryID)
		filters.CategoryID = &categoryID
	}
	if payload.StartDate != nil {
		filters.StartDate = payload.StartDate
	}
	if payload.EndDate != nil {
		filters.EndDate = payload.EndDate
	}

	transactions, total, err := s.rp.Transaction.GetTransactionsByUserID(userID, filters)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get transactions",
			Data:    nil,
			Error:   err,
		})
	}

	// Convert to response format
	transactionResponses := make([]dto.TransactionResponse, len(transactions))
	for i, t := range transactions {
		transactionResponses[i] = dto.TransactionResponse{
			ID:              t.ID,
			UserID:          t.UserID,
			CategoryID:      t.CategoryID,
			Category:        t.Category,
			Amount:          t.Amount,
			Description:     t.Description,
			TransactionDate: t.TransactionDate,
			Type:            t.Type,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
		}
	}

	totalPages := int((total + int64(payload.Limit) - 1) / int64(payload.Limit))
	response := dto.TransactionListResponse{
		Transactions: transactionResponses,
		Pagination: dto.PaginationResponse{
			Page:       payload.Page,
			Limit:      payload.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transactions retrieved successfully",
		Data:    response,
	})
}

func (s *Service) GetTransactionByIDRequest(id uint, phoneNumber string) *types.Response {
	user, err := s.rp.User.GetUserByPhone(phoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   err,
		})
	}

	transaction, err := s.rp.Transaction.GetTransactionByID(id)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Transaction not found",
			Data:    nil,
			Error:   err,
		})
	}

	// Check if transaction belongs to user
	if transaction.UserID != user.ID {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
	}

	response := dto.TransactionResponse{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		CategoryID:      transaction.CategoryID,
		Category:        transaction.Category,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate,
		Type:            transaction.Type,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transaction retrieved successfully",
		Data:    response,
	})
}

func (s *Service) UpdateTransactionRequest(id uint, payload dto.UpdateTransactionRequest, phoneNumber string) *types.Response {
	user, err := s.rp.User.GetUserByPhone(phoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   err,
		})
	}

	transaction, err := s.rp.Transaction.GetTransactionByID(id)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Transaction not found",
			Data:    nil,
			Error:   err,
		})
	}

	// Check if transaction belongs to user
	if transaction.UserID != user.ID {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
	}

	// Update fields if provided
	if payload.Amount != nil {
		transaction.Amount = *payload.Amount
	}
	if payload.CategoryID != nil {
		categoryID := uint(*payload.CategoryID)
		transaction.CategoryID = &categoryID
	}
	if payload.Type != nil {
		transaction.Type = models.TransactionType(*payload.Type)
	}
	if payload.Description != nil {
		transaction.Description = *payload.Description
	}

	updatedTransaction, err := s.rp.Transaction.UpdateTransaction(*transaction)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update transaction",
			Data:    nil,
			Error:   err,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transaction updated successfully",
		Data:    updatedTransaction,
	})
}

func (s *Service) DeleteTransactionRequest(id uint, phoneNumber string) *types.Response {
	user, err := s.rp.User.GetUserByPhone(phoneNumber)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   err,
		})
	}

	transaction, err := s.rp.Transaction.GetTransactionByID(id)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Transaction not found",
			Data:    nil,
			Error:   err,
		})
	}

	// Check if transaction belongs to user
	if transaction.UserID != user.ID {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
	}

	if err := s.rp.Transaction.DeleteTransaction(id); err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete transaction",
			Data:    nil,
			Error:   err,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transaction deleted successfully",
		Data:    nil,
	})
}

func (s *Service) GetTransactionsSummaryRequest(payload dto.TransactionSummaryRequest) *types.Response {
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

	filters := transaction.SummaryFilters{
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
		Month:     payload.Month,
		Year:      payload.Year,
	}

	summary, err := s.rp.Transaction.GetTransactionsSummary(userID, filters)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get transactions summary",
			Data:    nil,
			Error:   err,
		})
	}

	// Calculate percentages and convert to response format
	totalAmount := summary.TotalIncome + summary.TotalExpense
	categorySummary := make([]dto.CategorySummary, len(summary.CategorySummary))
	for i, cs := range summary.CategorySummary {
		percentage := float64(0)
		if totalAmount > 0 {
			percentage = (cs.TotalAmount / totalAmount) * 100
		}
		categorySummary[i] = dto.CategorySummary{
			CategoryID:   cs.CategoryID,
			CategoryName: cs.CategoryName,
			Type:         cs.Type,
			TotalAmount:  cs.TotalAmount,
			Count:        cs.Count,
			Percentage:   percentage,
		}
	}

	response := dto.TransactionSummaryResponse{
		TotalIncome:      summary.TotalIncome,
		TotalExpense:     summary.TotalExpense,
		Balance:          summary.TotalIncome - summary.TotalExpense,
		TransactionCount: summary.TransactionCount,
		IncomeCount:      summary.IncomeCount,
		ExpenseCount:     summary.ExpenseCount,
		CategorySummary:  categorySummary,
		Period: dto.PeriodInfo{
			StartDate: payload.StartDate,
			EndDate:   payload.EndDate,
			Month:     payload.Month,
			Year:      payload.Year,
		},
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Transaction summary retrieved successfully",
		Data:    response,
	})
}
