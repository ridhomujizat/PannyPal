package category

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/categories")
	group.POST("/", h.CreateCategory)
	group.GET("/", h.GetCategories)
	group.GET("/:id", h.GetCategoryByID)
	group.PUT("/:id", h.UpdateCategory)
	group.DELETE("/:id", h.DeleteCategory)
}
