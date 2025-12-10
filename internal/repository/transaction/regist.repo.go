package transaction

import (
	"context"
	"pannypal/internal/common/models"
	"pannypal/internal/pkg/redis"
	"time"

	database "pannypal/internal/pkg/db"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    *database.Database
}

type IRepository interface {
	CreateTransaction(model models.Transaction) (*models.Transaction, error)
	UpdateTransaction(model models.Transaction) (*models.Transaction, error)
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByUserID(userID *uint, filters TransactionFilters) ([]models.Transaction, int64, error)
	DeleteTransaction(id uint) error
	GetTransactionsSummary(userID *uint, filters SummaryFilters) (*TransactionSummary, error)
}

type TransactionFilters struct {
	Type       *models.TransactionType
	CategoryID *uint
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	Limit      int
}

type SummaryFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	Month     *int
	Year      *int
}

type TransactionSummary struct {
	TotalIncome      float64
	TotalExpense     float64
	TransactionCount int64
	IncomeCount      int64
	ExpenseCount     int64
	CategorySummary  []CategorySummaryData
}

type CategorySummaryData struct {
	CategoryID   uint
	CategoryName string
	Type         models.TransactionType
	TotalAmount  float64
	Count        int64
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) CreateTransaction(model models.Transaction) (*models.Transaction, error) {
	if err := r.db.WithContext(r.ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *Repository) UpdateTransaction(model models.Transaction) (*models.Transaction, error) {
	if err := r.db.WithContext(r.ctx).Save(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
func (r *Repository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(r.ctx).Preload("Category").Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *Repository) GetTransactionsByUserID(userID *uint, filters TransactionFilters) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var count int64

	query := r.db.WithContext(r.ctx).Model(&models.Transaction{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.StartDate != nil {
		query = query.Where("transaction_date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("transaction_date <= ?", *filters.EndDate)
	}

	// Count total records
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and preload
	offset := (filters.Page - 1) * filters.Limit
	if err := query.Preload("Category").Order("transaction_date DESC").
		Offset(offset).Limit(filters.Limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

func (r *Repository) DeleteTransaction(id uint) error {
	if err := r.db.WithContext(r.ctx).Delete(&models.Transaction{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTransactionsSummary(userID *uint, filters SummaryFilters) (*TransactionSummary, error) {
	query := r.db.WithContext(r.ctx).Model(&models.Transaction{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Apply date filters
	if filters.StartDate != nil {
		query = query.Where("transaction_date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("transaction_date <= ?", *filters.EndDate)
	}
	if filters.Month != nil && filters.Year != nil {
		query = query.Where("EXTRACT(MONTH FROM transaction_date) = ? AND EXTRACT(YEAR FROM transaction_date) = ?",
			*filters.Month, *filters.Year)
	}

	var summary TransactionSummary

	// Get total income
	var totalIncome float64
	if err := query.Where("type = ?", models.TypeIncome).Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome).Error; err != nil {
		return nil, err
	}
	summary.TotalIncome = totalIncome

	// Get total expense
	var totalExpense float64
	if err := query.Where("type = ?", models.TypeExpense).Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense).Error; err != nil {
		return nil, err
	}
	summary.TotalExpense = totalExpense

	// Get counts
	if err := query.Count(&summary.TransactionCount).Error; err != nil {
		return nil, err
	}

	if err := query.Where("type = ?", models.TypeIncome).Count(&summary.IncomeCount).Error; err != nil {
		return nil, err
	}

	if err := query.Where("type = ?", models.TypeExpense).Count(&summary.ExpenseCount).Error; err != nil {
		return nil, err
	}

	// Get category summary
	var categorySummary []CategorySummaryData
	if err := query.Select(`
		t.category_id,
		c.name as category_name,
		t.type,
		COALESCE(SUM(t.amount), 0) as total_amount,
		COUNT(*) as count
	`).Joins("JOIN categories c ON t.category_id = c.id").
		Group("t.category_id, c.name, t.type").
		Scan(&categorySummary).Error; err != nil {
		return nil, err
	}
	summary.CategorySummary = categorySummary

	return &summary, nil
}
