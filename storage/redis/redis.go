package redis

import (
	"context"
	"rate-limit-request/config"

	"github.com/go-redis/redis/v8"
)

// NewCacheRepository func implements the storage interface for app
func NewCacheRepository(config config.Config) *CacheRepository {
	return &CacheRepository{
		config: config,
	}
}

// CacheRepository is interface structure
type CacheRepository struct {
	config config.Config
	client *redis.Client
}

// Init initial
func (repo *CacheRepository) Init(ctx context.Context) error {
	config := repo.config
	repo.client = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	return repo.client.Ping(ctx).Err()
}

// FlushDB flush db
func (repo *CacheRepository) FlushDB(ctx context.Context) error {
	return repo.client.FlushDB(ctx).Err()
}
