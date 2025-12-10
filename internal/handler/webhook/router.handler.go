package webhook

import "github.com/gin-gonic/gin"

func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	group := e.Group("/webhook")
	group.POST("/waha", h.WebhookEventWaha)
}
