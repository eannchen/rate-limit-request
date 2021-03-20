package storage

import (
	"context"
	"rate-limit-request/model"
)

type CacheRepository interface {
	Init(context.Context) error
	FlushDB(context.Context) error

	IncreaseRateLimit(context.Context, string) (rateLimit model.RateLimit, err error)
	IncreaseRateLimitByLua(ctx context.Context, key string) (rateLimit model.RateLimit, err error)
}
