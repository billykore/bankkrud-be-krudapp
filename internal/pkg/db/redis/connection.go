package redis

import (
	"github.com/redis/go-redis/v9"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

// New creates a new Redis client.
func New(cfg *config.Configs) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}
