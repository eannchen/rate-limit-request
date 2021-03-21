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

	count, err := repo.client.Incr(ctx, key).Result()

	// key has just been created
	if err == nil && count == 1 {
		err = repo.client.ExpireAt(ctx, key, time.Now().Add(model.RateLimitExpireDuration)).Err()
	}

	if err != nil {
		return
	}

	rateLimit.Count = int(count)

	return
}

const increaseRateLimitLua = `
local key = KEYS[1]

local nowUnix = tonumber(ARGV[1])
local expireDuration = tonumber(ARGV[2])
local rateLimitMaximum  = tonumber(ARGV[3])

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
	redis.call('EXPIREAT', key, rateLimit["expire"])
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
