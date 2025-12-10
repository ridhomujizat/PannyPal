package ticketing

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
	CreateTicket(model models.Ticket) (*models.Ticket, error)
	UpdateTicket(model models.Ticket) (*models.Ticket, error)
	GetTicketByID(id uint) (*models.Ticket, error)
	DeleteTicket(id uint) error
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}
func (r *Repository) CreateTicket(model models.Ticket) (*models.Ticket, error) {
	if err := r.db.WithContext(r.ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *Repository) UpdateTicket(model models.Ticket) (*models.Ticket, error) {
	if err := r.db.WithContext(r.ctx).Save(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
func (r *Repository) GetTicketByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.WithContext(r.ctx).Where("id = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}
func (r *Repository) DeleteTicket(id uint) error {
	if err := r.db.WithContext(r.ctx).Where("id = ?", id).Delete(&models.Ticket{}).Error; err != nil {
		return err
	}
	return nil
}
