package AI

import (
	"context"

	"pannypal/internal/pkg/ai-connector"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/service/ai/dto"
	outgoingService "pannypal/internal/service/outgoing"
)

type Service struct {
	rp              repository.IRepository
	redis           redis.IRedis
	ctx             context.Context
	ai              *ai.AiClient
	outgoingService outgoingService.IService
}

type IService interface {
	InputTextCashflow(payload dto.InputTextCashflow) (*dto.TransactionResponseAi, string, error)
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository, aiClient *ai.AiClient, outgoingService outgoingService.IService) IService {
	return &Service{
		rp:              repository,
		redis:           redis,
		ctx:             ctx,
		ai:              aiClient,
		outgoingService: outgoingService,
	}
}
