package aicashflow

import (
	"context"

	"pannypal/internal/common/models"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/ai-connector"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/service/ai-cashflow/dto"
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
	InputTransaction(payload dto.InputTransaction) *types.Response
	PannyPalBotCashflow(payload dto.PayloadAICashflow)
	PannyPalBotCashflowReplayAction(payload dto.PayloadAICashflow, messageToReply models.MessageToReply)
	ReplayAction(payload dto.PayloadAICashflow, quotedStanzaID string) *types.Response
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
