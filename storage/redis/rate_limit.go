package redis

import (
	"context"
	"encoding/json"
	"rate-limit-request/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// IncreaseRateLimit IncreaseRateLimit
func (repo *CacheRepository) IncreaseRateLimit(ctx context.Context, key string) (rateLimit model.RateLimit, err error) {
	for i := 0; i < repo.config.Redis.MaxRetries; i++ {
		err = repo.client.Watch(ctx, func(tx *redis.Tx) error {

			HGetRes := tx.HGetAll(ctx, key)

			if err := HGetRes.Err(); err != nil && err != redis.Nil {
				return err
			}

			if err := HGetRes.Scan(&rateLimit); err != nil {
				return err
			}

			// create or reset flag
			if rateLimit.Count == 0 || time.Unix(int64(rateLimit.Expire), 0).Before(time.Now()) {
				tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
					rateLimit = model.RateLimit{
						Count:  1,
						Expire: int(time.Now().Add(model.RateLimitExpireDuration).Unix()),
					}
					return pipe.HMSet(ctx, key, map[string]interface{}{
						"count":  rateLimit.Count,
						"expire": rateLimit.Expire,
					}).Err()
				})
				return nil
			}

			// increase count
			if rateLimit.Count <= model.RateLimitMaximum {
				tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
					if pipe.HIncrBy(ctx, key, "count", 1).Err() != nil {
						return err
					}
					rateLimit.Count++
					return nil
				})
			}

			return err
		}, key)
		if err == nil {
			// Success.
			return
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return
	}
	return
}

const increaseRateLimitLua = `
local key = KEYS[1]

local nowUnix = tonumber(ARGV[1])
local expireDuration = tonumber(ARGV[2])
local rateLimitMaximum  = tonumber(ARGV[3])

local ipLimit = tonumber(ARGV[2])


local HGetRes = redis.call('HGETALL', key)
local rateLimit = {}
if #HGetRes ~= 0 then
	rateLimit["count"] = tonumber(HGetRes[2])
	rateLimit["expire"] = tonumber(HGetRes[4])
end

if #HGetRes == 0 or rateLimit["expire"] < nowUnix then
	rateLimit["count"] = 1
	rateLimit["expire"] = nowUnix + expireDuration

    redis.call('HMSET', key, "count", rateLimit["count"], "expire", rateLimit["expire"])
    return cjson.encode(rateLimit)
end


if rateLimit["count"] <= rateLimitMaximum then
    rateLimit["count"] = redis.call('HINCRBY', key, "count", 1)
end

return cjson.encode(rateLimit)
`

// IncreaseRateLimitByLua IncreaseRateLimitByLua
func (repo *CacheRepository) IncreaseRateLimitByLua(ctx context.Context, key string) (rateLimit model.RateLimit, err error) {
	luaScript := redis.NewScript(increaseRateLimitLua)

	nowUnix := time.Now().Unix()
	expireDuration := model.RateLimitExpireDuration.Seconds()
	rateLimitMaximum := model.RateLimitMaximum
	args := []interface{}{nowUnix, expireDuration, rateLimitMaximum}

	resp, err := luaScript.Run(ctx, repo.client, []string{key}, args...).Result()
	if err != nil {
		err = errors.Wrap(err, "run lua script")
		return
	}

	if err = json.Unmarshal([]byte(resp.(string)), &rateLimit); err != nil {
		err = errors.Wrap(err, "json unmarshal")
		return
	}
	return
}
