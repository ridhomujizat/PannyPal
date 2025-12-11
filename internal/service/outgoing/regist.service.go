package outgoing

import (
	"context"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/service/outgoing/dto"
)

type Service struct {
	ctx   context.Context
	redis redis.IRedis
	rp    repository.IRepository
}
type IService interface {
	HandleWebhookEventWaha(payload dto.PayloadOutgoing) (*dto.ResponseOutgoing, error)
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository) IService {
	return &Service{
		ctx:   ctx,
		redis: redis,
		rp:    repository,
	}
}
