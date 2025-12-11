package webhook

import (
	"context"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	aiCashFLowService "pannypal/internal/service/ai-cashflow"
)

type Service struct {
	ctx               context.Context
	redis             redis.IRedis
	rp                repository.IRepository
	aiCashFlowService aiCashFLowService.IService
}
type IService interface {
	HandleWebhookEventWaha(payload interface{}) *types.Response
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository, aiCashFlowService aiCashFLowService.IService) IService {
	return &Service{
		ctx:               ctx,
		redis:             redis,
		rp:                repository,
		aiCashFlowService: aiCashFlowService,
	}
}
