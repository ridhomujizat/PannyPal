package config

import (
	"context"
	"pannypal/internal/common/enum"
	"pannypal/internal/pkg/ai-connector"
	database "pannypal/internal/pkg/db"
	"pannypal/internal/pkg/rabbitmq"
	"pannypal/internal/pkg/redis"
	s3aws "pannypal/internal/pkg/storage/s3"
	"sync"
)

type Config struct {
	AppEnv        enum.EnvEnum `env:"APP_ENV" envDefault:"development"`
	AppPort       int          `env:"APP_PORT" envDefault:"9001"`
	RedisHost     string       `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     int          `env:"REDIS_PORT" envDefault:"6379"`
	RedisUser     string       `env:"REDIS_USER" envDefault:"default"`
	RedisPass     string       `env:"REDIS_PASS" envDefault:""`
	RedisPoolSize int          `env:"REDIS_POOL_SIZE" envDefault:"10"`
	RabbitHost    string       `env:"RABBIT_HOST" envDefault:"localhost"`
	RabbitPort    int          `env:"RABBIT_PORT" envDefault:"5672"`
	RabbitUser    string       `env:"RABBIT_USER" envDefault:"guest"`
	RabbitPass    string       `env:"RABBIT_PASS" envDefault:"guest"`
	DBHost        string       `env:"DB_HOST" envDefault:"localhost"`
	DBPort        int          `env:"DB_PORT" envDefault:"5432"`
	DBUser        string       `env:"DB_USER" envDefault:"postgres"`
	DBPass        string       `env:"DB_PASS" envDefault:""`
	DBName        string       `env:"DB_NAME" envDefault:"postgres"`
	GeminiAPIKey  string       `env:"GEMINI_API_KEY" envDefault:""`
	GeminiModel   string       `env:"GEMINI_MODEL" envDefault:"gemini-2.5-flash"`
	// AWSACCESSKEYID     string       `env:"AWS_ACCESS_KEY_ID" envDefault:""`
	// AWSSECRETACCESSKEY string       `env:"AWS_SECRET_ACCESS_KEY" envDefault:""`
	// AWSREGION          string       `env:"AWS_REGION" envDefault:"ap-southeast-3"`
	// AWSBUCKETNAME      string       `env:"AWS_BUCKET_NAME" envDefault:"sxored_development"`
}

type SetupServerDto struct {
	Ctx    *context.Context
	Cancel context.CancelFunc
	Wg     *sync.WaitGroup
	Env    *Config
	Db     *database.Database
	Rds    redis.IRedis
	Rb     *rabbitmq.ConnectionManager
	S3     *s3aws.Is3
	Ai     *ai.AiClient
}
