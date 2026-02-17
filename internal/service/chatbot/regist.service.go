package chatbot

import (
	"context"
	types "pannypal/internal/common/type"
	ai "pannypal/internal/pkg/ai-connector"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/repository/analytics"
	"pannypal/internal/service/chatbot/dto"
	"pannypal/internal/service/chatbot/engine"
)

type Service struct {
	ctx            context.Context
	redis          redis.IRedis
	rp             repository.IRepository
	aiClient       *ai.AiClient
	analysisEngine *engine.AnalysisEngine
}

type IService interface {
	SendMessage(payload dto.SendMessageRequest) *types.Response
	GetConversations(limit int) *types.Response
	GetConversation(sessionID string, limit int) *types.Response
	ClearConversation(sessionID string) *types.Response
}

func NewService(
	ctx context.Context,
	redis redis.IRedis,
	repository repository.IRepository,
	db *database.Database,
	aiClient *ai.AiClient,
) IService {
	// Create analytics repository for data fetching
	analyticsRepo := analytics.NewRepo(ctx, redis, db)

	// Create analysis engine
	analysisEngine := engine.NewAnalysisEngine(ctx, aiClient, analyticsRepo)

	return &Service{
		ctx:            ctx,
		redis:          redis,
		rp:             repository,
		aiClient:       aiClient,
		analysisEngine: analysisEngine,
	}
}
