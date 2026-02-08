package ai

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/validation"
	aiService "pannypal/internal/service/ai"
	"pannypal/internal/service/ai/dto"
)

type Handler struct {
	ctx context.Context
	ai  aiService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
}

func NewHandler(ctx context.Context, ai aiService.IService) IHandler {
	return &Handler{
		ctx: ctx,
		ai:  ai,
	}
}

// @Summary Process text input for cashflow using AI
// @Description Extract financial transactions from text input using AI
// @Tags AI APIs
// @Accept json
// @Produce json
// @Param request body dto.InputTextCashflow true "Text cashflow input"
// @Success 200 {object} types.Response "Transactions extracted successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Router /ai/cashflow/text [post]
func (h *Handler) InputTextCashflow(c *gin.Context) {
	send := c.MustGet("send").(func(r *types.Response))
	var payload dto.InputTextCashflow

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

	result, message, err := h.ai.InputTextCashflow(payload)
	if err != nil {
		send(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process cashflow input",
			Data:    nil,
			Error:   err,
		})
		return
	}

	send(&types.Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    result,
	})
}
