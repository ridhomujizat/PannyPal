package category

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
	CreateCategory(model models.Category) (*models.Category, error)
	UpdateCategory(model models.Category) (*models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	DeleteCategory(model models.Category) error
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) CreateCategory(model models.Category) (*models.Category, error) {
	if err := r.db.WithContext(r.ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *Repository) UpdateCategory(model models.Category) (*models.Category, error) {
	if err := r.db.WithContext(r.ctx).Save(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
func (r *Repository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(r.ctx).Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.WithContext(r.ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) DeleteCategory(model models.Category) error {
	if err := r.db.WithContext(r.ctx).Delete(&model).Error; err != nil {
		return err
	}
	return nil
}
