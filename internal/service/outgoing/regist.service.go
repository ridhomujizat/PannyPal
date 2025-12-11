package outgoing

import (
	"context"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
)

type Service struct {
	ctx   context.Context
	redis redis.IRedis
	rp    repository.IRepository
}
type IService interface {
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository) IService {
	return &Service{
		ctx:   ctx,
		redis: redis,
		rp:    repository,
	}
}
