package ai

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/ai")
	group.POST("/cashflow/text", h.InputTextCashflow)
}
