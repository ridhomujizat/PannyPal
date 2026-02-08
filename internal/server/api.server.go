package serverApp

import (
	"context"

	aiHandler "pannypal/internal/handler/ai"
	aicashflowHandler "pannypal/internal/handler/ai-cashflow"
	analyticsHandler "pannypal/internal/handler/analytics"
	budgetHandler "pannypal/internal/handler/budget"
	categoryHandler "pannypal/internal/handler/category"
	transactionHandler "pannypal/internal/handler/transaction"
	webhookHandler "pannypal/internal/handler/webhook"
	ai "pannypal/internal/pkg/ai-connector"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/middleware"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/redis"
	s3aws "pannypal/internal/pkg/storage/s3"
	"pannypal/internal/repository"
	"pannypal/internal/repository/analytics"
	"pannypal/internal/repository/bot"
	"pannypal/internal/repository/budget"
	"pannypal/internal/repository/category"
	logdata "pannypal/internal/repository/log-data"
	"pannypal/internal/repository/transaction"
	"pannypal/internal/repository/user"
	aiService "pannypal/internal/service/ai"
	aicashflowService "pannypal/internal/service/ai-cashflow"
	analyticsService "pannypal/internal/service/analytics"
	budgetService "pannypal/internal/service/budget"
	categoryService "pannypal/internal/service/category"
	outgoingService "pannypal/internal/service/outgoing"
	transactionService "pannypal/internal/service/transaction"
	webhookService "pannypal/internal/service/webhook"
	"sync"

	_ "pannypal/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(
	engine *gin.Engine,
	ctx context.Context,
	wg *sync.WaitGroup,
	db *database.Database,
	redis redis.IRedis,
	rb *rabbitmq.ConnectionManager,
	publisher *rabbitmq.Publisher,
	s3 *s3aws.Is3,
	ai *ai.AiClient) {

	InitMiddleware(engine, publisher)
	// Health check

	engine.GET("/health", func(c *gin.Context) {
		rabbitmqHealth := "unhealthy"
		reddistHealth := "unhealthy"
		databaseHealth := "unhealthy"
		rbCon := rb.GetConnection()

		if db != nil && !db.IsCloseConnection() {
			databaseHealth = "healthy"
		}

		if rbCon != nil && !rbCon.IsClosed() {
			rabbitmqHealth = "healthy"
		}
		if redis != nil && redis.Close() == nil {
			reddistHealth = "healthy"
		}
		c.JSON(200, gin.H{
			"status": 200,
			"service": gin.H{
				"rabbitmq": gin.H{
					"status": rabbitmqHealth,
				},
				"redis": gin.H{
					"status": reddistHealth,
				},
				"database": gin.H{
					"status": databaseHealth,
				},
			},
		})
	})

	// Swagger documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	e := engine.Group(BasePath())
	InitRoutes(e, ctx, wg, db, redis, rb, publisher, s3, ai)
}

func BasePath() string {
	return "/api"
}

func InitMiddleware(e *gin.Engine, publisher *rabbitmq.Publisher) {
	e.Use(middleware.CorsMiddleware())
	e.Use(middleware.RequestInit())
	e.Use(middleware.ResponseInit())
}

func InitRoutes(
	e *gin.RouterGroup,
	ctx context.Context,
	wg *sync.WaitGroup,
	db *database.Database,
	redis redis.IRedis,
	rb *rabbitmq.ConnectionManager,
	publisher *rabbitmq.Publisher,
	s3 *s3aws.Is3,
	ai *ai.AiClient) {

	// init repo
	rp := repository.IRepository{
		Category:    category.NewRepo(ctx, redis, db),
		Budget:      budget.NewRepo(ctx, redis, db),
		Transaction: transaction.NewRepo(ctx, redis, db),
		User:        user.NewRepo(ctx, redis, db),
		Analytics:   analytics.NewRepo(ctx, redis, db),
		LogData:     logdata.NewRepo(ctx, redis, db),
		Bot:         bot.NewRepo(ctx, redis, db),
	}
	// init services
	transactionSvc := transactionService.NewService(ctx, redis, rp)
	categorySvc := categoryService.NewService(ctx, redis, rp)
	budgetSvc := budgetService.NewService(ctx, redis, rp)
	analyticsSvc := analyticsService.NewService(ctx, redis, rp, db)
	outgoingSvc := outgoingService.NewService(ctx, redis, rp)
	aiSvc := aiService.NewService(ctx, redis, rp, ai, outgoingSvc)
	aiCashflowSvc := aicashflowService.NewService(ctx, redis, rp, ai, outgoingSvc)
	webhookSvc := webhookService.NewService(ctx, redis, rp, aiCashflowSvc)

	// init handlers
	transactionHandler := transactionHandler.NewHandler(ctx, rb, transactionSvc)
	categoryHandler := categoryHandler.NewHandler(ctx, rb, categorySvc)
	budgetHandler := budgetHandler.NewHandler(ctx, rb, budgetSvc)
	analyticsHandler := analyticsHandler.NewHandler(ctx, rb, analyticsSvc)
	aiHandler := aiHandler.NewHandler(ctx, aiSvc)
	aiCashflowHandler := aicashflowHandler.NewHandler(ctx, rb, aiCashflowSvc)
	webhookHandler := webhookHandler.NewHandler(ctx, rb, webhookSvc)

	// init handler routes
	transactionHandler.NewRoutes(e)
	categoryHandler.NewRoutes(e)
	budgetHandler.NewRoutes(e)
	analyticsHandler.NewRoutes(e)
	aiHandler.NewRoutes(e)
	aiCashflowHandler.NewRoutes(e)
	webhookHandler.NewRoutes(e)
}
