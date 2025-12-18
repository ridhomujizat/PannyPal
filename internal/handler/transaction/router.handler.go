package transaction

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/transactions")
	group.POST("", h.CreateTransaction)
	group.GET("", h.GetTransactions)
	group.GET("/:id", h.GetTransactionByID)
	group.PUT("/:id", h.UpdateTransaction)
	group.DELETE("/:id", h.DeleteTransaction)
	group.GET("/summary", h.GetTransactionsSummary)
}
