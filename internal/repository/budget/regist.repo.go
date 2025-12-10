package budget

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
	CreateBudget(model models.Budget) (*models.Budget, error)
	UpdateBudget(model models.Budget) (*models.Budget, error)
	DeleteBudget(model models.Budget) error
	GetBudgetByID(id uint) (*models.Budget, error)
	GetBudgetsByUserID(userID *uint, filters BudgetFilters) ([]models.Budget, error)
	GetBudgetStatus(userID *uint, filters BudgetStatusFilters) ([]BudgetStatusData, error)
}

type BudgetFilters struct {
	Month      *int
	Year       *int
	CategoryID *uint
}

type BudgetStatusFilters struct {
	Month *int
	Year  *int
}

type BudgetStatusData struct {
	BudgetID     uint
	CategoryID   uint
	CategoryName string
	BudgetAmount float64
	SpentAmount  float64
	Month        int
	Year         int
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) CreateBudget(model models.Budget) (*models.Budget, error) {
	if err := r.db.WithContext(r.ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
func (r *Repository) UpdateBudget(model models.Budget) (*models.Budget, error) {
	if err := r.db.WithContext(r.ctx).Save(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
func (r *Repository) DeleteBudget(model models.Budget) error {
	if err := r.db.WithContext(r.ctx).Delete(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetBudgetByID(id uint) (*models.Budget, error) {
	var budget models.Budget
	if err := r.db.WithContext(r.ctx).Preload("Category").Where("id = ?", id).First(&budget).Error; err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *Repository) GetBudgetsByUserID(userID *uint, filters BudgetFilters) ([]models.Budget, error) {
	var budgets []models.Budget
	query := r.db.WithContext(r.ctx)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if filters.Month != nil {
		query = query.Where("month = ?", *filters.Month)
	}
	if filters.Year != nil {
		query = query.Where("year = ?", *filters.Year)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if err := query.Preload("Category").Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *Repository) GetBudgetStatus(userID *uint, filters BudgetStatusFilters) ([]BudgetStatusData, error) {
	var statusData []BudgetStatusData

	whereClause := ""
	var args []interface{}

	if userID != nil {
		whereClause = "WHERE b.user_id = ?"
		args = append(args, *userID)
	}

	query := r.db.WithContext(r.ctx).Raw(`
		SELECT 
			b.id as budget_id,
			b.category_id,
			c.name as category_name,
			b.amount as budget_amount,
			COALESCE(SUM(t.amount), 0) as spent_amount,
			b.month,
			b.year
		FROM budgets b
		LEFT JOIN categories c ON b.category_id = c.id
		LEFT JOIN transactions t ON b.category_id = t.category_id 
			AND b.user_id = t.user_id 
			AND t.type = 'EXPENSE'
			AND EXTRACT(MONTH FROM t.transaction_date) = b.month
			AND EXTRACT(YEAR FROM t.transaction_date) = b.year
		`+whereClause+`
		GROUP BY b.id, b.category_id, c.name, b.amount, b.month, b.year
	`, args...)

	if filters.Month != nil {
		query = query.Where("b.month = ?", *filters.Month)
	}
	if filters.Year != nil {
		query = query.Where("b.year = ?", *filters.Year)
	}

	query = query.Group("b.id, b.category_id, c.name, b.amount, b.month, b.year")

	if err := query.Scan(&statusData).Error; err != nil {
		return nil, err
	}

	return statusData, nil
}
