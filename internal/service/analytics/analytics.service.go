package analytics

import (
	"net/http"
	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/repository/analytics"
	"pannypal/internal/service/analytics/dto"
	"time"
)

func (s *Service) GetMonthlyAnalyticsRequest(payload dto.MonthlyAnalyticsRequest) *types.Response {
	var userID *uint

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err := s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   err,
			})
		}
		userID = &user.ID
	}

	year := time.Now().Year()
	if payload.Year != nil {
		year = *payload.Year
	}

	data, err := s.analyticsRepo.GetMonthlyAnalytics(userID, year)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get monthly analytics",
			Data:    nil,
			Error:   err,
		})
	}

	monthNames := []string{"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	monthlyData := make([]dto.MonthlyDataPoint, len(data))
	totalIncome := float64(0)
	totalExpense := float64(0)
	var bestMonth, worstMonth *dto.MonthlyDataPoint

	for i, d := range data {
		balance := d.TotalIncome - d.TotalExpense
		monthlyData[i] = dto.MonthlyDataPoint{
			Month:            d.Month,
			MonthName:        monthNames[d.Month],
			Year:             d.Year,
			TotalIncome:      d.TotalIncome,
			TotalExpense:     d.TotalExpense,
			Balance:          balance,
			TransactionCount: d.TransactionCount,
		}

		totalIncome += d.TotalIncome
		totalExpense += d.TotalExpense

		// Find best and worst months by balance
		if bestMonth == nil || balance > bestMonth.Balance {
			bestMonth = &monthlyData[i]
		}
		if worstMonth == nil || balance < worstMonth.Balance {
			worstMonth = &monthlyData[i]
		}
	}

	response := dto.MonthlyAnalyticsResponse{
		Data:         monthlyData,
		Year:         year,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		TotalBalance: totalIncome - totalExpense,
		BestMonth:    bestMonth,
		WorstMonth:   worstMonth,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Monthly analytics retrieved successfully",
		Data:    response,
	})
}

func (s *Service) GetYearlyAnalyticsRequest(payload dto.YearlyAnalyticsRequest) *types.Response {
	var userID *uint

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err := s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   err,
			})
		}
		userID = &user.ID
	}

	currentYear := time.Now().Year()
	startYear := currentYear - 4 // Default to last 5 years
	endYear := currentYear

	if payload.StartYear != nil {
		startYear = *payload.StartYear
	}
	if payload.EndYear != nil {
		endYear = *payload.EndYear
	}

	data, err := s.analyticsRepo.GetYearlyAnalytics(userID, startYear, endYear)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get yearly analytics",
			Data:    nil,
			Error:   err,
		})
	}

	yearlyData := make([]dto.YearlyDataPoint, len(data))
	totalIncome := float64(0)
	totalExpense := float64(0)
	var bestYear, worstYear *dto.YearlyDataPoint

	for i, d := range data {
		balance := d.TotalIncome - d.TotalExpense
		yearlyData[i] = dto.YearlyDataPoint{
			Year:             d.Year,
			TotalIncome:      d.TotalIncome,
			TotalExpense:     d.TotalExpense,
			Balance:          balance,
			TransactionCount: d.TransactionCount,
		}

		totalIncome += d.TotalIncome
		totalExpense += d.TotalExpense

		// Find best and worst years by balance
		if bestYear == nil || balance > bestYear.Balance {
			bestYear = &yearlyData[i]
		}
		if worstYear == nil || balance < worstYear.Balance {
			worstYear = &yearlyData[i]
		}
	}

	response := dto.YearlyAnalyticsResponse{
		Data:         yearlyData,
		StartYear:    startYear,
		EndYear:      endYear,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		TotalBalance: totalIncome - totalExpense,
		BestYear:     bestYear,
		WorstYear:    worstYear,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Yearly analytics retrieved successfully",
		Data:    response,
	})
}

func (s *Service) GetCategoryAnalyticsRequest(payload dto.CategoryAnalyticsRequest) *types.Response {
	var userID *uint

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err := s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   err,
			})
		}
		userID = &user.ID
	}

	filters := analytics.CategoryAnalyticsFilters{
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
	}
	if payload.Type != nil {
		transactionType := models.TransactionType(*payload.Type)
		filters.Type = &transactionType
	}

	data, err := s.analyticsRepo.GetCategoryAnalytics(userID, filters)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get category analytics",
			Data:    nil,
			Error:   err,
		})
	}

	totalAmount := float64(0)
	totalCount := int64(0)

	// Calculate totals first
	for _, d := range data {
		totalAmount += d.TotalAmount
		totalCount += d.Count
	}

	categoryData := make([]dto.CategoryDataPoint, len(data))
	var topCategory *dto.CategoryDataPoint

	for i, d := range data {
		percentage := float64(0)
		if totalAmount > 0 {
			percentage = (d.TotalAmount / totalAmount) * 100
		}

		averageAmount := float64(0)
		if d.Count > 0 {
			averageAmount = d.TotalAmount / float64(d.Count)
		}

		categoryData[i] = dto.CategoryDataPoint{
			CategoryID:    d.CategoryID,
			CategoryName:  d.CategoryName,
			Type:          d.Type,
			TotalAmount:   d.TotalAmount,
			Count:         d.Count,
			Percentage:    percentage,
			AverageAmount: averageAmount,
		}

		// Find top category
		if topCategory == nil || d.TotalAmount > topCategory.TotalAmount {
			topCategory = &categoryData[i]
		}
	}

	response := dto.CategoryAnalyticsResponse{
		Data:             categoryData,
		TotalAmount:      totalAmount,
		TransactionCount: totalCount,
		Period: dto.PeriodInfo{
			StartDate: payload.StartDate,
			EndDate:   payload.EndDate,
		},
		TopCategory: topCategory,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Category analytics retrieved successfully",
		Data:    response,
	})
}

func (s *Service) GetDashboardAnalyticsRequest(payload dto.DashboardAnalyticsRequest) *types.Response {
	var userID *uint

	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		user, err := s.rp.User.GetUserByPhone(*payload.PhoneNumber)
		if err != nil {
			return helper.ParseResponse(&types.Response{
				Code:    http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   err,
			})
		}
		userID = &user.ID
	}

	// Default to current month if not provided
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	if payload.StartDate != nil && payload.EndDate != nil {
		startDate = *payload.StartDate
		endDate = *payload.EndDate
	}

	data, err := s.analyticsRepo.GetDashboardAnalytics(userID, startDate, endDate)
	if err != nil {
		return helper.ParseResponse(&types.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get dashboard analytics",
			Data:    nil,
			Error:   err,
		})
	}

	// Calculate percentage changes
	incomeChange := float64(0)
	if data.PreviousIncome > 0 {
		incomeChange = ((data.CurrentIncome - data.PreviousIncome) / data.PreviousIncome) * 100
	} else if data.CurrentIncome > 0 {
		incomeChange = 100
	}

	expenseChange := float64(0)
	if data.PreviousExpense > 0 {
		expenseChange = ((data.CurrentExpense - data.PreviousExpense) / data.PreviousExpense) * 100
	} else if data.CurrentExpense > 0 {
		expenseChange = 100
	}

	// Calculate previous period dates
	duration := endDate.Sub(startDate)
	previousEndDate := startDate.Add(-time.Second)
	previousStartDate := previousEndDate.Add(-duration)

	totalBalance := data.TotalIncomeAllTime - data.TotalExpenseAllTime

	response := dto.DashboardAnalyticsResponse{
		TotalBalance:      totalBalance,
		Income:            data.CurrentIncome,
		Expense:           data.CurrentExpense,
		IncomeChange:      incomeChange,
		ExpenseChange:     expenseChange,
		StartDate:         &startDate,
		EndDate:           &endDate,
		PreviousStartDate: &previousStartDate,
		PreviousEndDate:   &previousEndDate,
	}

	return helper.ParseResponse(&types.Response{
		Code:    http.StatusOK,
		Message: "Dashboard analytics retrieved successfully",
		Data:    response,
	})
}
