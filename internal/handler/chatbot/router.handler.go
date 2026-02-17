package chatbot

import (
	"github.com/gin-gonic/gin"
)

// NewRoutes registers chatbot routes
func (h *Handler) NewRoutes(e *gin.RouterGroup) {
	chatbot := e.Group("/chatbot")
	{
		// Send message to chatbot
		chatbot.POST("/send", h.SendMessage)

		// Get all conversations
		chatbot.GET("/conversations", h.GetConversations)

		// Get specific conversation
		chatbot.GET("/conversations/:session_id", h.GetConversation)

		// Clear/archive conversation
		chatbot.POST("/conversations/:session_id/clear", h.ClearConversation)
	}
}
