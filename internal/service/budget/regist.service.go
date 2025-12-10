package budget

import (
	"context"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/service/budget/dto"
)

type Service struct {
	ctx   context.Context
	redis redis.IRedis
	rp    repository.IRepository
}

type IService interface {
	CreateBudgetRequest(payload dto.CreateBudgetRequest) *types.Response
	GetBudgetsRequest(payload dto.GetBudgetsRequest) *types.Response
	GetBudgetByIDRequest(id uint, phoneNumber string) *types.Response
	UpdateBudgetRequest(id uint, payload dto.UpdateBudgetRequest, phoneNumber string) *types.Response
	DeleteBudgetRequest(id uint, phoneNumber string) *types.Response
	GetBudgetStatusRequest(payload dto.BudgetStatusRequest) *types.Response
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository) IService {
	return &Service{
		ctx:   ctx,
		redis: redis,
		rp:    repository,
	}
}
