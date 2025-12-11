package bot

import (
	"context"
	"pannypal/internal/common/models"
	"pannypal/internal/pkg/redis"

	database "pannypal/internal/pkg/db"

	"gorm.io/gorm"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    *database.Database
}

type IRepository interface {
	CreateMessageToReply(d models.MessageToReply) (*models.MessageToReply, error)
	MessageToReplyMessage(messageID string) (*models.MessageToReply, error)
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) CreateMessageToReply(d models.MessageToReply) (*models.MessageToReply, error) {
	if err := r.db.WithContext(r.ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *Repository) MessageToReplyMessage(messageID string) (*models.MessageToReply, error) {
	var data models.MessageToReply
	err := r.db.WithContext(r.ctx).Where("message_id = ?", messageID).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}
