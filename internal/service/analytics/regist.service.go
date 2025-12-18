package analytics

import (
	"context"
	types "pannypal/internal/common/type"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/repository/analytics"
	"pannypal/internal/service/analytics/dto"
)

type Service struct {
	ctx           context.Context
	redis         redis.IRedis
	rp            repository.IRepository
	analyticsRepo analytics.IRepository
}

type IService interface {
	GetMonthlyAnalyticsRequest(payload dto.MonthlyAnalyticsRequest) *types.Response
	GetYearlyAnalyticsRequest(payload dto.YearlyAnalyticsRequest) *types.Response
	GetCategoryAnalyticsRequest(payload dto.CategoryAnalyticsRequest) *types.Response
	GetDashboardAnalyticsRequest(payload dto.DashboardAnalyticsRequest) *types.Response
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository, db *database.Database) IService {
	analyticsRepo := analytics.NewRepo(ctx, redis, db)
	return &Service{
		ctx:           ctx,
		redis:         redis,
		rp:            repository,
		analyticsRepo: analyticsRepo,
	}
}
