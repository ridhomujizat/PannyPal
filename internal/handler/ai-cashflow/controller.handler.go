package aicashflow

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/validation"
	aicashflowService "pannypal/internal/service/ai-cashflow"
	"pannypal/internal/service/ai-cashflow/dto"
)

type Handler struct {
	ctx      context.Context
	rabbitmq *rabbitmq.ConnectionManager
	ai       aicashflowService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
}

func NewHandler(ctx context.Context, rabbitmq *rabbitmq.ConnectionManager, ai aicashflowService.IService) IHandler {
	return &Handler{
		ctx:      ctx,
		rabbitmq: rabbitmq,
		ai:       ai,
	}
}

// @Summary Create AI transaction
// @Description Create a new transaction using AI processing
// @Tags AI Cashflow APIs
// @Accept json
// @Produce json
// @Param transaction body dto.InputTransaction true "Transaction data"
// @Success 201 {object} dto.TransactionResponseAi "Transaction created successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /ai-cashflow/transaction [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.InputTransaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(&payload); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(h.ai.InputTransaction(payload))
}
