package models

import (
	"gorm.io/gorm"
)

// ChatConversation represents a chat session between user and AI
type ChatConversation struct {
	gorm.Model
	SessionID string `gorm:"type:varchar(100);uniqueIndex;not null" json:"session_id"` // UUID for one chat session
	Title     string `gorm:"type:varchar(255)" json:"title"`                           // Auto-generated from first query
	IsActive  bool   `gorm:"default:true" json:"is_active"`                            // Session is still active or finished

	// Relations
	Messages []ChatMessage `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	gorm.Model
	ConversationID uint   `gorm:"not null;index" json:"conversation_id"`
	Role           string `gorm:"type:varchar(20);not null" json:"role"` // "user" or "assistant"
	Content        string `gorm:"type:text;not null" json:"content"`     // Message text
	Metadata       string `gorm:"type:text" json:"metadata"`             // JSON for additional data (chart data, query results, etc)
	TokenUsed      int    `gorm:"default:0" json:"token_used"`
	ResponseTime   int    `gorm:"default:0" json:"response_time"` // in milliseconds

	// Relations
	Conversation ChatConversation `gorm:"foreignKey:ConversationID" json:"-"`
}
