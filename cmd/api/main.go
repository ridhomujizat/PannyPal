package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	config "pannypal/configs"
	"pannypal/internal/pkg/ai-connector"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/helper"
	"pannypal/internal/pkg/logger"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/redis"
	"pannypal/internal/pkg/validation"
	serverApp "pannypal/internal/server"
	"sync"
	"syscall"

	_ "pannypal/docs"

	"github.com/gin-gonic/gin"
)

// @title Cash Flow API
// @version 1.0
// @description API lengkap untuk monitoring pengeluaran dan pemasukan cash flow
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /api

func main() {

	logger.Setup()
	env, err := config.GetEnv()
	if err != nil {
		logger.Error.Println("Error getting environment", err)
		panic(err)
	}

	// // Set timezone to Indonesia (WIB - Western Indonesia Time)
	// loc, err := time.LoadLocation("Asia/Jakarta")
	// if err != nil {
	// 	logger.Error.Println("Error loading timezone", err)
	// 	// Fallback to UTC if timezone loading fails
	// 	loc = time.UTC
	// }
	// time.Local = loc

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Setup Redis
	redis, err := setupRedis(ctx, env)
	if err != nil {
		logger.Error.Println("Error setting up Redis", err)
		cancel()
		return
	}

	// Setup RabbitMQ
	rabbit, err := setupRabbitMQ(ctx, env)
	if err != nil {
		logger.Error.Println("Error setting up RabbitMQ", err)
		cancel()
		return
	}
	// Setup Database
	db, err := setupDB(env)
	if err != nil {
		logger.Error.Println("Error setting up Database", err)
		cancel()
		return
	}
	// s3Client, err := setupS3(ctx, env, redis)
	// if err != nil {
	// 	logger.Error.Println("Error setting up S3", err)
	// 	cancel()
	// 	return
	// }
	// Setup Server
	aiClient := setupAi(ctx)

	setupServer(&config.SetupServerDto{
		Rds:    redis,
		Env:    env,
		Ctx:    &ctx,
		Cancel: cancel,
		Db:     db,
		Wg:     &wg,
		Rb:     rabbit,
		Ai:     aiClient,
		// S3:     &s3Client,
	})

}

func setupRedis(ctx context.Context, env *config.Config) (redis.IRedis, error) {
	return redis.Setup(ctx, &redis.Config{
		Host:     env.RedisHost,
		Username: env.RedisUser,
		Port:     env.RedisPort,
		Password: env.RedisPass,
		PoolSize: env.RedisPoolSize,
	})
}

func setupRabbitMQ(ctx context.Context, env *config.Config) (*rabbitmq.ConnectionManager, error) {
	return rabbitmq.NewConnectionManager(ctx, &rabbitmq.Config{
		Username: env.RabbitUser,
		Password: env.RabbitPass,
		Host:     env.RabbitHost,
		Port:     env.RabbitPort,
	})
}

func setupDB(env *config.Config) (*database.Database, error) {
	return database.Setup(&database.Config{
		Host:     env.DBHost,
		Port:     env.DBPort,
		User:     env.DBUser,
		Password: env.DBPass,
		Database: env.DBName,
		SSLMode:  "disable",
		Driver:   "postgres",
	})
}

// func setupS3(ctx context.Context, env *config.Config, redis redis.IRedis) (s3aws.Is3, error) {
// 	return s3aws.NewS3Client(ctx, s3aws.S3Config{
// 		AWSRegion:          env.AWSREGION,
// 		AWSAccessKeyID:     env.AWSACCESSKEYID,
// 		AWSSecretAccessKey: env.AWSSECRETACCESSKEY,
// 	}, env.AWSBUCKETNAME, redis)
// }

func setupAi(ctx context.Context) *ai.AiClient {
	apiKey := helper.GetEnv("GEMINI_API_KEY")
	model := helper.GetEnv("GEMINI_MODEL")

	fmt.Printf("Gemini API Key configured: %t\n", apiKey != "")
	fmt.Printf("Gemini Model: %s\n", model)

	return ai.NewAiClient(
		ctx,
		&ai.Config{
			GeminiAPIKey: apiKey,
			GeminiModel:  model,
		},
	)
}

func setupServer(payload *config.SetupServerDto) {
	rds := payload.Rds
	env := payload.Env
	ctx := payload.Ctx
	cancel := payload.Cancel
	wg := payload.Wg
	rb := payload.Rb
	db := payload.Db
	s3 := payload.S3
	ai := payload.Ai

	defer func() {
		if rds != nil {
			_ = rds.Close()
		}
		cancel()
		wg.Wait()
	}()

	err := validation.Setup()
	if err != nil {
		logger.Error.Println("Failed to setup validation")
		panic(err)
	}

	e := gin.Default()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.AppPort),
		Handler: e,
		//IdleTimeout:  1 * time.Minute,
		//ReadTimeout:  1 * time.Minute,
		//WriteTimeout: 1 * time.Minute,
	}

	publisher, err := rabbitmq.NewPublisher(*ctx, rb)
	if err != nil {
		panic(err)
	}

	serverApp.Setup(e, *ctx, wg, db, rds, rb, publisher, s3, ai)
	if payload.Env.AppEnv != "development" {
		serverApp.InitWorker(*ctx, rds, db, rb, publisher, s3)
	}

	go func() {
		logger.HTTP.Println("========= Server Started =========")
		logger.HTTP.Println("=========", env.AppPort, "=========")
		if err := server.ListenAndServe(); err != nil {
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan
	_ = server.Shutdown(*ctx)
}
