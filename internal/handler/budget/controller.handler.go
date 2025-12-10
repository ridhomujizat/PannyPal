package budget

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/validation"
	budgetService "pannypal/internal/service/budget"
	"pannypal/internal/service/budget/dto"
)

type Handler struct {
	ctx           context.Context
	rabbitmq      *rabbitmq.ConnectionManager
	budgetService budgetService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	CreateBudget(c *gin.Context)
	GetBudgets(c *gin.Context)
	GetBudgetByID(c *gin.Context)
	UpdateBudget(c *gin.Context)
	DeleteBudget(c *gin.Context)
	GetBudgetStatus(c *gin.Context)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, budgetService budgetService.IService) IHandler {
	return &Handler{
		ctx:           ctx,
		rabbitmq:      rabbitmq,
		budgetService: budgetService,
	}
}

// CreateBudget godoc
// @Summary Create budget
// @Description Set monthly budget for a category
// @Tags Budget APIs
// @Accept json
// @Produce json
// @Param budget body dto.CreateBudgetRequest true "Budget data"
// @Success 201 {object} types.Response "Budget created successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /budgets [post]
func (h *Handler) CreateBudget(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.CreateBudgetRequest

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

	send(h.budgetService.CreateBudgetRequest(payload))
}

// GetBudgets godoc
// @Summary Get budgets
// @Description Get list of budgets with optional filters
// @Tags Budget APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param month query int false "Filter by month"
// @Param year query int false "Filter by year"
// @Param category_id query int false "Filter by category"
// @Success 200 {object} types.Response "Budgets retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /budgets [get]
func (h *Handler) GetBudgets(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.GetBudgetsRequest

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

	send(h.budgetService.GetBudgetsRequest(payload))
}

// GetBudgetByID godoc
// @Summary Get budget by ID
// @Description Get specific budget by ID
// @Tags Budget APIs
// @Accept json
// @Produce json
// @Param id path int true "Budget ID"
// @Param phone_number query string false "User's phone number"
// @Success 200 {object} types.Response "Budget retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /budgets/{id} [get]
func (h *Handler) GetBudgetByID(c *gin.Context) {
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
	budgetID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid budget ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.budgetService.GetBudgetByIDRequest(uint(budgetID), phoneNumber))
}

// UpdateBudget godoc
// @Summary Update budget
// @Description Update existing budget
// @Tags Budget APIs
// @Accept json
// @Produce json
// @Param id path int true "Budget ID"
// @Param phone_number query string false "User's phone number"
// @Param budget body dto.UpdateBudgetRequest true "Updated budget data"
// @Success 200 {object} types.Response "Budget updated successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /budgets/{id} [put]
func (h *Handler) UpdateBudget(c *gin.Context) {
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
	budgetID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid budget ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	var payload dto.UpdateBudgetRequest
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

	send(h.budgetService.UpdateBudgetRequest(uint(budgetID), payload, phoneNumber))
}

// DeleteBudget godoc
// @Summary Delete budget
// @Description Delete budget by ID
// @Tags Budget APIs
// @Accept json
// @Produce json
// @Param id path int true "Budget ID"
// @Param phone_number query string false "User's phone number"
// @Success 200 {object} types.Response "Budget deleted successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 404 {object} types.Response "Not Found"
// @Router /budgets/{id} [delete]
func (h *Handler) DeleteBudget(c *gin.Context) {
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
	budgetID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid budget ID",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.budgetService.DeleteBudgetRequest(uint(budgetID), phoneNumber))
}

// GetBudgetStatus godoc
// @Summary Get budget status
// @Description Check budget usage status
// @Tags Budget APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param month query int false "Filter by month (1-12)"
// @Param year query int false "Filter by year"
// @Success 200 {object} dto.BudgetStatusListResponse "Budget status retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /budgets/status [get]
func (h *Handler) GetBudgetStatus(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.BudgetStatusRequest

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

	send(h.budgetService.GetBudgetStatusRequest(payload))
}
