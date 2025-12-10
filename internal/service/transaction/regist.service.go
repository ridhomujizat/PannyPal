package transaction

import (
	"context"
	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/service/transaction/dto"
)

type Service struct {
	ctx   context.Context
	redis redis.IRedis
	rp    repository.IRepository
}

type IService interface {
	CreateTransactionRequest(payload dto.CreateTransactionRequest) *types.Response
	GetTransactionsRequest(payload dto.GetTransactionsRequest) *types.Response
	GetTransactionByIDRequest(id uint, phoneNumber string) *types.Response
	UpdateTransactionRequest(id uint, payload dto.UpdateTransactionRequest, phoneNumber string) *types.Response
	DeleteTransactionRequest(id uint, phoneNumber string) *types.Response
	GetTransactionsSummaryRequest(payload dto.TransactionSummaryRequest) *types.Response
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository) IService {
	return &Service{
		ctx:   ctx,
		redis: redis,
		rp:    repository,
	}
}
