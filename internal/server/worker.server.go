package serverApp

import (
	"context"
	"fmt"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/logger"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/redis"
	s3aws "pannypal/internal/pkg/storage/s3"

	"time"

	"github.com/panjf2000/ants"
)

func InitWorker(ctx context.Context, redis redis.IRedis, db *database.Database, rb *rabbitmq.ConnectionManager, publisher *rabbitmq.Publisher, s3 *s3aws.Is3) {
	// init repo
	// rp := repository.IRepository{}
	// init service
	// init handlers
	poolOpts := ants.Options{
		ExpiryDuration: time.Hour,
		PreAlloc:       true,
		Nonblocking:    true,
		PanicHandler: func(i interface{}) {
			logger.Error.Printf("Worker microservice panic: %v\n", i)
		},
	}

	pool, err := ants.NewPool(100, ants.WithOptions(poolOpts))
	if err != nil {
		panic(fmt.Errorf("failed to create worker pool: %w", err))
	}
	defer pool.Release()

	err = pool.Submit(func() {
		// Initialize the RabbitMQ subscriber for kyb
		// if err := kybAnalyzerHandler.SubscribeKybAnalyzer(); err != nil {
		// 	logger.Error.Printf("Failed to initialize O Analyzer subscriber: %v\n", err)
		// }

	})
	if err != nil {
		panic(fmt.Errorf("failed to submit task to pool: %w", err))
	}
}
