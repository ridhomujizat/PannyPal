package chatbot

import (
	"context"
	"pannypal/internal/common/models"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/redis"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    *database.Database
}

type IRepository interface {
	// Conversations
	CreateConversation(conv models.ChatConversation) (*models.ChatConversation, error)
	GetConversationBySession(sessionID string) (*models.ChatConversation, error)
	GetAllConversations(limit int) ([]models.ChatConversation, error)
	UpdateConversation(conv models.ChatConversation) error

	// Messages
	CreateMessage(msg models.ChatMessage) (*models.ChatMessage, error)
	GetConversationMessages(conversationID uint, limit int) ([]models.ChatMessage, error)
	GetMessageCount(conversationID uint) (int64, error)
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

// CreateConversation creates a new conversation
func (r *Repository) CreateConversation(conv models.ChatConversation) (*models.ChatConversation, error) {
	if err := r.db.WithContext(r.ctx).Create(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

// GetConversationBySession retrieves a conversation by session ID
func (r *Repository) GetConversationBySession(sessionID string) (*models.ChatConversation, error) {
	var conv models.ChatConversation
	if err := r.db.WithContext(r.ctx).Where("session_id = ?", sessionID).First(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

// GetAllConversations retrieves all conversations with limit
func (r *Repository) GetAllConversations(limit int) ([]models.ChatConversation, error) {
	var convs []models.ChatConversation
	query := r.db.WithContext(r.ctx).Order("updated_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&convs).Error; err != nil {
		return nil, err
	}
	return convs, nil
}

// UpdateConversation updates an existing conversation
func (r *Repository) UpdateConversation(conv models.ChatConversation) error {
	return r.db.WithContext(r.ctx).Save(&conv).Error
}

// CreateMessage creates a new message
func (r *Repository) CreateMessage(msg models.ChatMessage) (*models.ChatMessage, error) {
	if err := r.db.WithContext(r.ctx).Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetConversationMessages retrieves messages for a conversation
func (r *Repository) GetConversationMessages(conversationID uint, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := r.db.WithContext(r.ctx).Where("conversation_id = ?", conversationID).Order("created_at ASC")

	if limit > 0 {
		// Get the last N messages
		query = query.Limit(limit).Order("created_at DESC")
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	if limit > 0 {
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
	}

	return messages, nil
}

// GetMessageCount counts messages in a conversation
func (r *Repository) GetMessageCount(conversationID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(r.ctx).Model(&models.ChatMessage{}).Where("conversation_id = ?", conversationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
