package chatbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/service/chatbot/dto"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SendMessage handles sending a message and getting AI response
func (s *Service) SendMessage(payload dto.SendMessageRequest) *types.Response {
	startTime := time.Now()

	// Validate message
	if strings.TrimSpace(payload.Message) == "" {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusBadRequest,
			Message: "Message cannot be empty",
			Data:    nil,
			Error:   fmt.Errorf("empty message"),
		})
	}

	// Get or create conversation
	var conversation *models.ChatConversation
	var err error

	if payload.SessionID != "" {
		// Get existing conversation
		conversation, err = s.rp.Chatbot.GetConversationBySession(payload.SessionID)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "Conversation not found",
				Data:    nil,
				Error:   err,
			})
		}
	} else {
		// Create new conversation
		sessionID := uuid.New().String()
		title := s.generateTitle(payload.Message)

		conversation = &models.ChatConversation{
			SessionID: sessionID,
			Title:     title,
			IsActive:  true,
		}

		conversation, err = s.rp.Chatbot.CreateConversation(*conversation)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create conversation",
				Data:    nil,
				Error:   err,
			})
		}
	}

	// Save user message
	userMessage := models.ChatMessage{
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        payload.Message,
		Metadata:       "",
		TokenUsed:      0,
		ResponseTime:   0,
	}

	_, err = s.rp.Chatbot.CreateMessage(userMessage)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save user message",
			Data:    nil,
			Error:   err,
		})
	}

	// Get conversation history (last 10 messages)
	history, err := s.rp.Chatbot.GetConversationMessages(conversation.ID, 10)
	if err != nil {
		history = []models.ChatMessage{}
	}

	// Analyze query using engine
	metadata, textResponse, tokenUsed, responseTime, err := s.analysisEngine.AnalyzeQuery(
		payload.Message,
		history,
	)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to analyze query",
			Data:    nil,
			Error:   err,
		})
	}

	// Save assistant message
	metadataJSON, _ := json.Marshal(metadata)
	assistantMessage := models.ChatMessage{
		ConversationID: conversation.ID,
		Role:           "assistant",
		Content:        textResponse,
		Metadata:       string(metadataJSON),
		TokenUsed:      tokenUsed,
		ResponseTime:   responseTime,
	}

	savedMessage, err := s.rp.Chatbot.CreateMessage(assistantMessage)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save assistant message",
			Data:    nil,
			Error:   err,
		})
	}

	// Build response
	response := dto.ChatMessageResponse{
		SessionID:    conversation.SessionID,
		Role:         savedMessage.Role,
		Content:      savedMessage.Content,
		Metadata:     metadata,
		TokenUsed:    savedMessage.TokenUsed,
		ResponseTime: int(time.Since(startTime).Milliseconds()),
		CreatedAt:    savedMessage.CreatedAt,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    response,
		Error:   nil,
	})
}

// GetConversations retrieves all conversations
func (s *Service) GetConversations(limit int) *types.Response {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	conversations, err := s.rp.Chatbot.GetAllConversations(limit)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get conversations",
			Data:    nil,
			Error:   err,
		})
	}

	// Build response list
	responseList := make([]dto.ConversationResponse, 0)
	for _, conv := range conversations {
		messageCount, _ := s.rp.Chatbot.GetMessageCount(conv.ID)

		// Get last message time
		messages, _ := s.rp.Chatbot.GetConversationMessages(conv.ID, 1)
		lastMessageTime := conv.UpdatedAt
		if len(messages) > 0 {
			lastMessageTime = messages[0].CreatedAt
		}

		responseList = append(responseList, dto.ConversationResponse{
			SessionID:    conv.SessionID,
			Title:        conv.Title,
			IsActive:     conv.IsActive,
			MessageCount: int(messageCount),
			LastMessage:  lastMessageTime,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    responseList,
		Error:   nil,
	})
}

// GetConversation retrieves a specific conversation with messages
func (s *Service) GetConversation(sessionID string, limit int) *types.Response {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	conversation, err := s.rp.Chatbot.GetConversationBySession(sessionID)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Conversation not found",
			Data:    nil,
			Error:   err,
		})
	}

	messages, err := s.rp.Chatbot.GetConversationMessages(conversation.ID, limit)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get messages",
			Data:    nil,
			Error:   err,
		})
	}

	// Build message responses
	messageResponses := make([]dto.ChatMessageResponse, 0)
	for _, msg := range messages {
		var metadata interface{}
		if msg.Metadata != "" {
			json.Unmarshal([]byte(msg.Metadata), &metadata)
		}

		messageResponses = append(messageResponses, dto.ChatMessageResponse{
			SessionID:    conversation.SessionID,
			Role:         msg.Role,
			Content:      msg.Content,
			Metadata:     metadata,
			TokenUsed:    msg.TokenUsed,
			ResponseTime: msg.ResponseTime,
			CreatedAt:    msg.CreatedAt,
		})
	}

	response := map[string]interface{}{
		"session_id":    conversation.SessionID,
		"title":         conversation.Title,
		"is_active":     conversation.IsActive,
		"message_count": len(messages),
		"messages":      messageResponses,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    response,
		Error:   nil,
	})
}

// ClearConversation archives/deactivates a conversation
func (s *Service) ClearConversation(sessionID string) *types.Response {
	conversation, err := s.rp.Chatbot.GetConversationBySession(sessionID)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusNotFound,
			Message: "Conversation not found",
			Data:    nil,
			Error:   err,
		})
	}

	conversation.IsActive = false
	err = s.rp.Chatbot.UpdateConversation(*conversation)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update conversation",
			Data:    nil,
			Error:   err,
		})
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Conversation cleared successfully",
		Data:    nil,
		Error:   nil,
	})
}

// generateTitle generates a title from the first message
func (s *Service) generateTitle(message string) string {
	// Take first 50 characters or until first newline
	title := message
	if len(title) > 50 {
		title = title[:50] + "..."
	}

	// Remove newlines
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.TrimSpace(title)

	if title == "" {
		title = "New Chat"
	}

	return title
}
