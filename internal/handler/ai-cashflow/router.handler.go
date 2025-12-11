package aicashflow

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/ai-cashflow")
	group.POST("/transaction", h.CreateTransaction)
	group.POST("/bot", h.PannyPalBotCashflow)
	group.POST("/bot/reply-action", h.PannyPalBotCashflowReplayAction)
}
