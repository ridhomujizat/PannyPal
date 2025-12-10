package budget

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/budgets")
	group.POST("/", h.CreateBudget)
	group.GET("/", h.GetBudgets)
	group.GET("/:id", h.GetBudgetByID)
	group.PUT("/:id", h.UpdateBudget)
	group.DELETE("/:id", h.DeleteBudget)
	group.GET("/status", h.GetBudgetStatus)
}
