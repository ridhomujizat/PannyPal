package user

import (
	"context"
	"pannypal/internal/pkg/redis"

	database "pannypal/internal/pkg/db"
)

type Repository struct {
	ctx   context.Context
	redis redis.IRedis
	db    *database.Database
}

type IRepository interface {
}

func NewRepo(ctx context.Context, redis redis.IRedis, db *database.Database) IRepository {
	return &Repository{
		ctx:   ctx,
		redis: redis,
		db:    db,
	}
}
