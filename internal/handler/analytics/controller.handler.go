package analytics

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/validation"
	analyticsService "pannypal/internal/service/analytics"
	"pannypal/internal/service/analytics/dto"
)

type Handler struct {
	ctx              context.Context
	rabbitmq         *rabbitmq.ConnectionManager
	analyticsService analyticsService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	GetMonthlyAnalytics(c *gin.Context)
	GetYearlyAnalytics(c *gin.Context)
	GetCategoryAnalytics(c *gin.Context)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, analyticsService analyticsService.IService) IHandler {
	return &Handler{
		ctx:              ctx,
		rabbitmq:         rabbitmq,
		analyticsService: analyticsService,
	}
}

// GetMonthlyAnalytics godoc
// @Summary Get monthly analytics
// @Description Get monthly breakdown for a year
// @Tags Analytics APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param year query int true "Year"
// @Success 200 {object} dto.MonthlyAnalyticsResponse "Monthly analytics retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /analytics/monthly [get]
func (h *Handler) GetMonthlyAnalytics(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.MonthlyAnalyticsRequest

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

	send(h.analyticsService.GetMonthlyAnalyticsRequest(payload))
}

// GetYearlyAnalytics godoc
// @Summary Get yearly analytics
// @Description Get yearly comparison
// @Tags Analytics APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param start_year query int false "Start year"
// @Param end_year query int false "End year"
// @Success 200 {object} dto.YearlyAnalyticsResponse "Yearly analytics retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /analytics/yearly [get]
func (h *Handler) GetYearlyAnalytics(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.YearlyAnalyticsRequest

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

	send(h.analyticsService.GetYearlyAnalyticsRequest(payload))
}

// GetCategoryAnalytics godoc
// @Summary Get category analytics
// @Description Get spending/income by category
// @Tags Analytics APIs
// @Accept json
// @Produce json
// @Param phone_number query string false "User's phone number"
// @Param start_date query string false "Analysis from this date"
// @Param end_date query string false "Analysis until this date"
// @Param type query string false "INCOME or EXPENSE"
// @Success 200 {object} dto.CategoryAnalyticsResponse "Category analytics retrieved successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /analytics/categories [get]
func (h *Handler) GetCategoryAnalytics(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.CategoryAnalyticsRequest

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

	send(h.analyticsService.GetCategoryAnalyticsRequest(payload))
}
