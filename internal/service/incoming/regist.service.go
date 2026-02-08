package incoming

import (
	"context"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	AI "pannypal/internal/service/ai"
	"pannypal/internal/service/outgoing"
)

type Service struct {
	ctx      context.Context
	redis    redis.IRedis
	rp       repository.IRepository
	ai       AI.IService
	outgoing outgoing.IService
}
type IService interface {
	HandleWebhookEventBaileys(payload interface{}) *types.Response
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository, ai AI.IService, outgoing outgoing.IService) IService {
	return &Service{
		ctx:      ctx,
		redis:    redis,
		rp:       repository,
		ai:       ai,
		outgoing: outgoing,
	}
}
