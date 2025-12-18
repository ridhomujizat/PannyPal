package analytics

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
	GetMonthlyAnalytics(userID *uint, year int) ([]MonthlyAnalyticsData, error)
	GetYearlyAnalytics(userID *uint, startYear, endYear int) ([]YearlyAnalyticsData, error)
	GetCategoryAnalytics(userID *uint, filters CategoryAnalyticsFilters) ([]CategoryAnalyticsData, error)
	GetDashboardAnalytics(userID *uint, startDate, endDate time.Time) (*DashboardAnalyticsData, error)
}

type MonthlyAnalyticsData struct {
	Month            int
	Year             int
	TotalIncome      float64
	TotalExpense     float64
	TransactionCount int64
}

type YearlyAnalyticsData struct {
	Year             int
	TotalIncome      float64
	TotalExpense     float64
	TransactionCount int64
}

type CategoryAnalyticsData struct {
	CategoryID   uint
	CategoryName string
	Type         models.TransactionType
	TotalAmount  float64
	Count        int64
}

type CategoryAnalyticsFilters struct {
	StartDate *time.Time
	EndDate   *time.Time
	Type      *models.TransactionType
}

type DashboardAnalyticsData struct {
	CurrentIncome       float64
	CurrentExpense      float64
	PreviousIncome      float64
	PreviousExpense     float64
	TotalIncomeAllTime  float64
	TotalExpenseAllTime float64
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}

func (r *Repository) GetMonthlyAnalytics(userID *uint, year int) ([]MonthlyAnalyticsData, error) {
	var data []MonthlyAnalyticsData

	query := `
		SELECT 
			EXTRACT(MONTH FROM transaction_date) as month,
			EXTRACT(YEAR FROM transaction_date) as year,
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as total_expense,
			COUNT(*) as transaction_count
		FROM transactions`

	args := []interface{}{}
	if userID != nil {
		query += ` WHERE user_id = ? AND`
		args = append(args, *userID)
	} else {
		query += ` WHERE`
	}

	query += ` EXTRACT(YEAR FROM transaction_date) = ?
		GROUP BY EXTRACT(MONTH FROM transaction_date), EXTRACT(YEAR FROM transaction_date)
		ORDER BY month`

	if userID != nil {
		args = append(args, year)
	} else {
		args = append(args, year)
	}

	err := r.db.WithContext(r.ctx).Raw(query, args...).Scan(&data).Error

	return data, err
}

func (r *Repository) GetYearlyAnalytics(userID *uint, startYear, endYear int) ([]YearlyAnalyticsData, error) {
	var data []YearlyAnalyticsData

	query := `
		SELECT 
			EXTRACT(YEAR FROM transaction_date) as year,
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as total_expense,
			COUNT(*) as transaction_count
		FROM transactions`

	args := []interface{}{}
	if userID != nil {
		query += ` WHERE user_id = ?`
		args = append(args, *userID)
		query += ` AND EXTRACT(YEAR FROM transaction_date) BETWEEN ? AND ?`
	} else {
		query += ` WHERE EXTRACT(YEAR FROM transaction_date) BETWEEN ? AND ?`
	}
	args = append(args, startYear, endYear)

	query += ` GROUP BY EXTRACT(YEAR FROM transaction_date) ORDER BY year`

	err := r.db.WithContext(r.ctx).Raw(query, args...).Scan(&data).Error

	return data, err
}

func (r *Repository) GetCategoryAnalytics(userID *uint, filters CategoryAnalyticsFilters) ([]CategoryAnalyticsData, error) {
	var data []CategoryAnalyticsData

	queryStr := `
		SELECT 
			t.category_id,
			c.name as category_name,
			t.type,
			COALESCE(SUM(t.amount), 0) as total_amount,
			COUNT(*) as count
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id`

	args := []interface{}{}
	if userID != nil {
		queryStr += ` WHERE t.user_id = ?`
		args = append(args, *userID)
	}

	hasWhere := userID != nil

	if filters.StartDate != nil {
		if hasWhere {
			queryStr += " AND t.transaction_date >= ?"
		} else {
			queryStr += " WHERE t.transaction_date >= ?"
			hasWhere = true
		}
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		if hasWhere {
			queryStr += " AND t.transaction_date <= ?"
		} else {
			queryStr += " WHERE t.transaction_date <= ?"
			hasWhere = true
		}
		args = append(args, *filters.EndDate)
	}

	if filters.Type != nil {
		if hasWhere {
			queryStr += " AND t.type = ?"
		} else {
			queryStr += " WHERE t.type = ?"
			hasWhere = true
		}
		args = append(args, *filters.Type)
	}

	queryStr += `
		GROUP BY t.category_id, c.name, t.type
		ORDER BY total_amount DESC
	`

	err := r.db.WithContext(r.ctx).Raw(queryStr, args...).Scan(&data).Error
	return data, err
}

func (r *Repository) GetDashboardAnalytics(userID *uint, startDate, endDate time.Time) (*DashboardAnalyticsData, error) {
	var result DashboardAnalyticsData

	// Get current period data
	currentQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as expense
		FROM transactions
		WHERE transaction_date >= ? AND transaction_date <= ?`

	currentArgs := []interface{}{startDate, endDate}
	if userID != nil {
		currentQuery += " AND user_id = ?"
		currentArgs = append(currentArgs, *userID)
	}

	row := r.db.WithContext(r.ctx).Raw(currentQuery, currentArgs...).Row()
	var currentIncome, currentExpense float64
	if err := row.Scan(&currentIncome, &currentExpense); err != nil {
		return nil, err
	}

	result.CurrentIncome = currentIncome
	result.CurrentExpense = currentExpense

	// Get previous period data (same duration, before current period)
	duration := endDate.Sub(startDate)
	previousEndDate := startDate.Add(-time.Second)
	previousStartDate := previousEndDate.Add(-duration)

	previousQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as expense
		FROM transactions
		WHERE transaction_date >= ? AND transaction_date <= ?`

	previousArgs := []interface{}{previousStartDate, previousEndDate}
	if userID != nil {
		previousQuery += " AND user_id = ?"
		previousArgs = append(previousArgs, *userID)
	}

	row = r.db.WithContext(r.ctx).Raw(previousQuery, previousArgs...).Row()
	var previousIncome, previousExpense float64
	if err := row.Scan(&previousIncome, &previousExpense); err != nil {
		return nil, err
	}

	result.PreviousIncome = previousIncome
	result.PreviousExpense = previousExpense

	// Get all time total
	allTimeQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as expense
		FROM transactions`

	allTimeArgs := []interface{}{}
	if userID != nil {
		allTimeQuery += " WHERE user_id = ?"
		allTimeArgs = append(allTimeArgs, *userID)
	}

	row = r.db.WithContext(r.ctx).Raw(allTimeQuery, allTimeArgs...).Row()
	var totalIncome, totalExpense float64
	if err := row.Scan(&totalIncome, &totalExpense); err != nil {
		return nil, err
	}

	result.TotalIncomeAllTime = totalIncome
	result.TotalExpenseAllTime = totalExpense

	return &result, nil
}
