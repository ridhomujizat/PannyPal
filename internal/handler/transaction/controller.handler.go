package transaction

import (
	"context"
	"net/http"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/validation"
	transactionService "pannypal/internal/service/transaction"
	"pannypal/internal/service/transaction/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctx                context.Context
	rabbitmq           *rabbitmq.ConnectionManager
	transactionService transactionService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	CreateTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
	GetTransactionByID(c *gin.Context)
	UpdateTransaction(c *gin.Context)
	DeleteTransaction(c *gin.Context)
	GetTransactionsSummary(c *gin.Context)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, transactionService transactionService.IService) IHandler {
	return &Handler{
		ctx:                ctx,
		rabbitmq:           rabbitmq,
		transactionService: transactionService,
	}
}

// CreateTransaction godoc
// @Summary Create new transaction
// @Description Create new income or expense transaction
// @Tags Transaction APIs
// @Accept json
// @Produce json
// @Param transaction body dto.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} types.Response "Transaction created successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /transactions [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		send(helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    err.Error(),
			Error:   err,
		}))
		return
	}

	send(h.transactionService.CreateTransactionRequest(payload))
}

// GetTransactions godoc
// @Summary Get all transactions
// @Description List all transactions with filters and pagination
// @Tags Transaction APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param type query string false "Filter by type (INCOME/EXPENSE)"
// @Param category_id query int false "Filter by category ID"
// @Param start_date query string false "Filter transactions from this date (ISO format)"
// @Param end_date query string false "Filter transactions until this date (ISO format)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} types.Response "Transactions retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /transactions [get]
func (h *Handler) GetTransactions(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.GetTransactionsRequest

	if err := c.ShouldBindQuery(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		send(helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    err.Error(),
			Error:   err,
		}))
		return
	}

	send(h.transactionService.GetTransactionsRequest(payload))
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Get specific transaction by ID
// @Tags Transaction APIs
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param phone_number query string false "User's phone number"
// @Success 200 {object} types.Response "Transaction retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /transactions/{id} [get]
func (h *Handler) GetTransactionByID(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	phoneNumber := c.Query("phone_number")

	if phoneNumber == "" {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Phone number is required",
			Data:    nil,
		})
		return
	}

	id := c.Param("id")
	transactionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid transaction ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.transactionService.GetTransactionByIDRequest(uint(transactionID), phoneNumber))
}

// UpdateTransaction godoc
// @Summary Update transaction
// @Description Update existing transaction. All fields are optional
// @Tags Transaction APIs
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param phone_number query string false "User's phone number"
// @Param transaction body dto.UpdateTransactionRequest true "Updated transaction data"
// @Success 200 {object} types.Response "Transaction updated successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /transactions/{id} [put]
func (h *Handler) UpdateTransaction(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	phoneNumber := c.Query("phone_number")

	if phoneNumber == "" {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Phone number is required",
			Data:    nil,
		})
		return
	}

	id := c.Param("id")
	transactionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid transaction ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	var payload dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		send(helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    err.Error(),
			Error:   err,
		}))
		return
	}

	send(h.transactionService.UpdateTransactionRequest(uint(transactionID), payload, phoneNumber))
}

// DeleteTransaction godoc
// @Summary Delete transaction
// @Description Delete transaction by ID
// @Tags Transaction APIs
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param phone_number query string false "User's phone number"
// @Success 200 {object} types.Response "Transaction deleted successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /transactions/{id} [delete]
func (h *Handler) DeleteTransaction(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	phoneNumber := c.Query("phone_number")

	if phoneNumber == "" {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Phone number is required",
			Data:    nil,
		})
		return
	}

	id := c.Param("id")
	transactionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid transaction ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.transactionService.DeleteTransactionRequest(uint(transactionID), phoneNumber))
}

// GetTransactionsSummary godoc
// @Summary Get transaction summary
// @Description Get summary of transactions with totals and category breakdown
// @Tags Transaction APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param start_date query string false "Summary from this date"
// @Param end_date query string false "Summary until this date"
// @Param month query int false "Filter by specific month (1-12)"
// @Param year query int false "Filter by specific year"
// @Success 200 {object} dto.TransactionSummaryResponse "Transaction summary retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /transactions/summary [get]
func (h *Handler) GetTransactionsSummary(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.TransactionSummaryRequest

	if err := c.ShouldBindQuery(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(payload); err != nil {
		send(helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Data:    err.Error(),
			Error:   err,
		}))
		return
	}

	send(h.transactionService.GetTransactionsSummaryRequest(payload))
}
