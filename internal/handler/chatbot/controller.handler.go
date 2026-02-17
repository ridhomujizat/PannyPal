package chatbot

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/pkg/validation"
	chatbotService "pannypal/internal/service/chatbot"
	"pannypal/internal/service/chatbot/dto"
)

type Handler struct {
	ctx            context.Context
	chatbotService chatbotService.IService
}

type IHandler interface {
	NewRoutes(e *gin.RouterGroup)
	SendMessage(c *gin.Context)
	GetConversations(c *gin.Context)
	GetConversation(c *gin.Context)
	ClearConversation(c *gin.Context)
}

func NewHandler(ctx context.Context, chatbotService chatbotService.IService) IHandler {
	return &Handler{
		ctx:            ctx,
		chatbotService: chatbotService,
	}
}

// SendMessage godoc
// @Summary Send a message to chatbot
// @Description Send a message and get AI response with analysis
// @Tags Chatbot APIs
// @Accept json
// @Produce json
// @Param request body dto.SendMessageRequest true "Send message request"
// @Success 200 {object} dto.ChatMessageResponse "Message sent successfully"
// @Failure 400 {object} types.Response "Bad Request"
// @Failure 500 {object} types.Response "Internal Server Error"
// @Router /chatbot/send [post]
func (h *Handler) SendMessage(c *gin.Context) {
	var payload dto.SendMessageRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Data:    nil,
			Error:   err,
		}))
		return
	}

	if err := validation.Validate(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Data:    nil,
			Error:   err,
		}))
		return
	}

	result := h.chatbotService.SendMessage(payload)
	c.JSON(result.Code, result)
}

// GetConversations godoc
// @Summary Get all conversations
// @Description Retrieve all chat conversations
// @Tags Chatbot APIs
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of conversations" default(20)
// @Success 200 {object} []dto.ConversationResponse "Conversations retrieved successfully"
// @Failure 500 {object} types.Response "Internal Server Error"
// @Router /chatbot/conversations [get]
func (h *Handler) GetConversations(c *gin.Context) {
	limit := 20 // Default
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	result := h.chatbotService.GetConversations(limit)
	c.JSON(result.Code, result)
}

// GetConversation godoc
// @Summary Get a specific conversation
// @Description Retrieve messages from a specific conversation
// @Tags Chatbot APIs
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Param limit query int false "Limit number of messages" default(50)
// @Success 200 {object} types.Response "Conversation retrieved successfully"
// @Failure 404 {object} types.Response "Conversation not found"
// @Router /chatbot/conversations/{session_id} [get]
func (h *Handler) GetConversation(c *gin.Context) {
	sessionID := c.Param("session_id")

	limit := 50 // Default
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	result := h.chatbotService.GetConversation(sessionID, limit)
	c.JSON(result.Code, result)
}

// ClearConversation godoc
// @Summary Clear/archive a conversation
// @Description Mark a conversation as inactive
// @Tags Chatbot APIs
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} types.Response "Conversation cleared successfully"
// @Failure 404 {object} types.Response "Conversation not found"
// @Router /chatbot/conversations/{session_id}/clear [post]
func (h *Handler) ClearConversation(c *gin.Context) {
	sessionID := c.Param("session_id")

	result := h.chatbotService.ClearConversation(sessionID)
	c.JSON(result.Code, result)
}
