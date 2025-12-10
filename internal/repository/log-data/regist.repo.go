package logdata

import (
	"context"
	"pannypal/internal/common/models"
	"pannypal/internal/pkg/redis"

	database "pannypal/internal/pkg/db"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    *database.Database
}

type IRepository interface {
	CreateLogWaha(d models.LogWaha) (*models.LogWaha, error)
	CreateLogPrompt(d models.LogPrompt) (*models.LogPrompt, error)
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) CreateLogWaha(d models.LogWaha) (*models.LogWaha, error) {
	if err := r.db.WithContext(r.ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *Repository) CreateLogPrompt(d models.LogPrompt) (*models.LogPrompt, error) {
	if err := r.db.WithContext(r.ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
