package incoming

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/incoming")
	group.POST("/baileys", h.WebhookEventBaileys)
}
