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

// @Summary Process PannyPal Bot Cashflow
// @Description Process cashflow transaction through PannyPal Bot
// @Tags AI Cashflow APIs
// @Accept json
// @Produce json
// @Param payload body dto.PayloadAICashflow true "Bot cashflow payload"
// @Success 200 {object} types.Response "Request processed successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /ai-cashflow/bot [post]
func (h *Handler) PannyPalBotCashflow(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.PayloadAICashflow

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

	// Call the service method (no return value expected)
	h.ai.PannyPalBotCashflow(payload)

	send(&types.Response{
		Code:    http.StatusOK,
		Message: "Bot cashflow request processed successfully",
		Data:    nil,
	})
}

// @Summary Process PannyPal Bot Cashflow Reply Action
// @Description Process reply action for bot cashflow transaction (save, edit, cancel)
// @Tags AI Cashflow APIs
// @Accept json
// @Produce json
// @Param request body dto.ReplayActionSimpleRequest true "Reply action request"
// @Success 200 {object} types.Response "Reply action processed successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /ai-cashflow/bot/reply-action [post]
func (h *Handler) PannyPalBotCashflowReplayAction(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var request dto.ReplayActionSimpleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		})
		return
	}

	if err := validation.Validate(&request); err != nil {
		send(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Data:    nil,
			Error:   err,
		})
		return
	}

	// Call the ReplayAction service method
	response := h.ai.ReplayAction(request.Payload, request.QuotedStanzaID)
	send(response)
}
