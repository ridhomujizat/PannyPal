package engine

import (
	"encoding/json"
	"fmt"
	"pannypal/internal/common/models"
	"pannypal/internal/repository/analytics"
	"time"
)

// DataFetcher handles fetching data from various sources
type DataFetcher struct {
	analyticsRepo analytics.IRepository
}

// NewDataFetcher creates a new DataFetcher instance
func NewDataFetcher(analyticsRepo analytics.IRepository) *DataFetcher {
	return &DataFetcher{
		analyticsRepo: analyticsRepo,
	}
}

// FetchTransactionSummary fetches transaction summary for a date range
func (d *DataFetcher) FetchTransactionSummary(startDate, endDate time.Time) (string, error) {
	// Since this is personal use, userID is nil to get all data
	data, err := d.analyticsRepo.GetDashboardAnalytics(nil, startDate, endDate)
	if err != nil {
		return "", err
	}

	summary := map[string]interface{}{
		"period": map[string]string{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
		"current": map[string]float64{
			"income":  data.CurrentIncome,
			"expense": data.CurrentExpense,
			"net":     data.CurrentIncome - data.CurrentExpense,
		},
		"previous": map[string]float64{
			"income":  data.PreviousIncome,
			"expense": data.PreviousExpense,
			"net":     data.PreviousIncome - data.PreviousExpense,
		},
		"all_time": map[string]float64{
			"total_income":  data.TotalIncomeAllTime,
			"total_expense": data.TotalExpenseAllTime,
			"net":           data.TotalIncomeAllTime - data.TotalExpenseAllTime,
		},
	}

	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// FetchCategoryBreakdown fetches category breakdown for a specific type and period
func (d *DataFetcher) FetchCategoryBreakdown(txType models.TransactionType, startDate, endDate *time.Time) (string, error) {
	filters := analytics.CategoryAnalyticsFilters{
		StartDate: startDate,
		EndDate:   endDate,
		Type:      &txType,
	}

	data, err := d.analyticsRepo.GetCategoryAnalytics(nil, filters)
	if err != nil {
		return "", err
	}

	var total float64
	categories := make([]map[string]interface{}, 0)
	for _, cat := range data {
		total += cat.TotalAmount
		categories = append(categories, map[string]interface{}{
			"category_id":   cat.CategoryID,
			"category_name": cat.CategoryName,
			"amount":        cat.TotalAmount,
			"count":         cat.Count,
			"type":          cat.Type,
		})
	}

	// Add percentage
	for i := range categories {
		if total > 0 {
			categories[i]["percentage"] = (categories[i]["amount"].(float64) / total) * 100
		} else {
			categories[i]["percentage"] = 0.0
		}
	}

	result := map[string]interface{}{
		"type":       txType,
		"total":      total,
		"categories": categories,
		"period": map[string]interface{}{
			"start": formatTimePtr(startDate),
			"end":   formatTimePtr(endDate),
		},
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// FetchMonthlyTrend fetches monthly trend for a specific year
func (d *DataFetcher) FetchMonthlyTrend(year int) (string, error) {
	data, err := d.analyticsRepo.GetMonthlyAnalytics(nil, year)
	if err != nil {
		return "", err
	}

	months := make([]map[string]interface{}, 0)
	for _, month := range data {
		months = append(months, map[string]interface{}{
			"month":             month.Month,
			"year":              month.Year,
			"income":            month.TotalIncome,
			"expense":           month.TotalExpense,
			"net":               month.TotalIncome - month.TotalExpense,
			"transaction_count": month.TransactionCount,
		})
	}

	result := map[string]interface{}{
		"year":   year,
		"months": months,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// FetchYearlyTrend fetches yearly trend for a range of years
func (d *DataFetcher) FetchYearlyTrend(startYear, endYear int) (string, error) {
	data, err := d.analyticsRepo.GetYearlyAnalytics(nil, startYear, endYear)
	if err != nil {
		return "", err
	}

	years := make([]map[string]interface{}, 0)
	for _, year := range data {
		years = append(years, map[string]interface{}{
			"year":              year.Year,
			"income":            year.TotalIncome,
			"expense":           year.TotalExpense,
			"net":               year.TotalIncome - year.TotalExpense,
			"transaction_count": year.TransactionCount,
		})
	}

	result := map[string]interface{}{
		"period": map[string]int{
			"start_year": startYear,
			"end_year":   endYear,
		},
		"years": years,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// FetchCurrentMonthData fetches data for current month
func (d *DataFetcher) FetchCurrentMonthData() (string, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return d.FetchTransactionSummary(startOfMonth, endOfMonth)
}

// FetchLastNMonths fetches data for last N months
func (d *DataFetcher) FetchLastNMonths(n int) (string, error) {
	now := time.Now()
	endDate := now
	startDate := now.AddDate(0, -n, 0)

	monthlyData := make([]map[string]interface{}, 0)

	for i := 0; i < n; i++ {
		monthStart := now.AddDate(0, -i, 0)
		monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
		monthEnd := monthStart.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		summary, err := d.FetchTransactionSummary(monthStart, monthEnd)
		if err != nil {
			return "", err
		}

		var summaryData map[string]interface{}
		if err := json.Unmarshal([]byte(summary), &summaryData); err != nil {
			return "", err
		}

		monthlyData = append(monthlyData, summaryData)
	}

	result := map[string]interface{}{
		"period": map[string]string{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
		"months": monthlyData,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// FetchAllCategoriesComparison fetches comparison of all categories
func (d *DataFetcher) FetchAllCategoriesComparison(startDate, endDate *time.Time) (string, error) {
	// Fetch expenses
	expenseType := models.TypeExpense
	expenseData, err := d.FetchCategoryBreakdown(expenseType, startDate, endDate)
	if err != nil {
		return "", err
	}

	// Fetch income
	incomeType := models.TypeIncome
	incomeData, err := d.FetchCategoryBreakdown(incomeType, startDate, endDate)
	if err != nil {
		return "", err
	}

	var expenses, income map[string]interface{}
	json.Unmarshal([]byte(expenseData), &expenses)
	json.Unmarshal([]byte(incomeData), &income)

	result := map[string]interface{}{
		"expenses": expenses,
		"income":   income,
		"period": map[string]interface{}{
			"start": formatTimePtr(startDate),
			"end":   formatTimePtr(endDate),
		},
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Helper function to format time pointer
func formatTimePtr(t *time.Time) string {
	if t == nil {
		return "all_time"
	}
	return t.Format("2006-01-02")
}

// BuildContextString builds a context string from data for AI prompt
func (d *DataFetcher) BuildContextString(dataType string, params map[string]interface{}) (string, error) {
	switch dataType {
	case "current_month":
		return d.FetchCurrentMonthData()
	case "transaction_summary":
		startDate := params["start_date"].(time.Time)
		endDate := params["end_date"].(time.Time)
		return d.FetchTransactionSummary(startDate, endDate)
	case "category_breakdown":
		txType := params["type"].(models.TransactionType)
		startDate, _ := params["start_date"].(*time.Time)
		endDate, _ := params["end_date"].(*time.Time)
		return d.FetchCategoryBreakdown(txType, startDate, endDate)
	case "monthly_trend":
		year := params["year"].(int)
		return d.FetchMonthlyTrend(year)
	case "yearly_trend":
		startYear := params["start_year"].(int)
		endYear := params["end_year"].(int)
		return d.FetchYearlyTrend(startYear, endYear)
	case "last_n_months":
		n := params["n"].(int)
		return d.FetchLastNMonths(n)
	case "all_categories":
		startDate, _ := params["start_date"].(*time.Time)
		endDate, _ := params["end_date"].(*time.Time)
		return d.FetchAllCategoriesComparison(startDate, endDate)
	default:
		return "", fmt.Errorf("unknown data type: %s", dataType)
	}
}
