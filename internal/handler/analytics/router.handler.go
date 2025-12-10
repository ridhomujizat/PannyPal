package analytics

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/analytics")
	group.GET("/monthly", h.GetMonthlyAnalytics)
	group.GET("/yearly", h.GetYearlyAnalytics)
	group.GET("/categories", h.GetCategoryAnalytics)
}
