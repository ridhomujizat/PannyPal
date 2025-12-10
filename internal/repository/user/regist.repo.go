package user

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
	CreateUser(d models.User) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	UpdateUser(d models.User) (*models.User, error)
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) CreateUser(d models.User) (*models.User, error) {
	if err := r.db.WithContext(r.ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *Repository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(r.ctx).Where("phone_number = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(d models.User) (*models.User, error) {
	if err := r.db.WithContext(r.ctx).Save(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
