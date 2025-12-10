package category

import (
	"context"

	types "pannypal/internal/common/type"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/repository"
	"pannypal/internal/service/category/dto"
)

type Service struct {
	rp    repository.IRepository
	redis redis.IRedis
	ctx   context.Context
}

type IService interface {
	CreateCategoryRequest(payload dto.CreateCategoryRequest) *types.Response
	GetCategoriesRequest() *types.Response
	GetCategoryByIDRequest(id uint) *types.Response
	UpdateCategoryRequest(id uint, payload dto.UpdateCategoryRequest) *types.Response
	DeleteCategoryRequest(id uint) *types.Response
}

func NewService(ctx context.Context, redis redis.IRedis, repository repository.IRepository) IService {
	return &Service{
		rp:    repository,
		redis: redis,
		ctx:   ctx,
	}
}
