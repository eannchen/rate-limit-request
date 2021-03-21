# Rate Limit Request

## Story
請實作一個 server 並滿足以下要求:
- 每個 IP 每分鐘僅能接受 60 個 requests
- 在首頁顯示目前的 request 量，超過限制的話則顯示 “Error” ，例如在一分鐘內第 30 個 request 則顯示 30，第 61 個 request 則顯示 Error
- 可以使用任意資料庫，也可以自行設計 in-memory 資料結構，並在文件中說明理由
- 請附上測試
- 請不要使用任何現成的 rate limit library

你不需要實作:
- 資料持久化
- 設計網頁

## Implement
rate limit middleware 的功能要盡可能的快速，盡可能不要影響商業邏輯的速度，並且限制訪問的資料是允許遺失的，所以使用 Redis 這個 in-memory key-value database 作為儲存，本題目可以 client IP 作為 key 辨識請求次數，相當方便。

Redis 的操作邏輯，我做了兩個方法：
### 方法一
使用 String type，當有請求時即寫入 value(請求次數) +1，若 key 剛被創建則也設置 ttl，時間到了就取消訪問限制。在此的指令直接設值，並且 Redis 對指令是 single thread，具 atomicity，避免 race condition 問題。
```go
count, err := repo.client.Incr(ctx, key).Result()

// key has just been created
if err == nil && count == 1 {
    err = repo.client.ExpireAt(ctx, key, time.Now().Add(model.RateLimitExpireDuration)).Err()
}
```

### 方法二
使用 Hash type，先讀取請求次數，再決定要不要 value(請求次數) +1，並且也儲存 rate limit 的結束時間，若有需要 response 就可以使用。由於這裡會拿請求次數做 +1 的判斷，所以使用 Lua script 讓整個 Redis 操作是 atomicity 的，避免 race condition 而使實際請求大於計數。
```lua
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
```


## Usage
### 啟動 Redis 及 phpRedisAdmin
1. docker-compose up
```s
$ make dkrps-up
```
2. 訪問 phpRedisAdmin
```
http://localhost:8081
```

### 關閉 Redis 及 phpRedisAdmin
1. docker-compose down
```s
$ make dkrps-down
```

### 運行 Server
1. 啟動 Go API server
```s
$ make run
```

2. 請求
```s
$ curl -i http://localhost:8080/app
```
response
```
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
X-Ratelimit-Current: 1
X-Ratelimit-Maximum: 60
Date: Sun, 21 Mar 2021 16:10:50 GMT
Content-Length: 1

1
```
```
HTTP/1.1 429 Too Many Requests
Content-Type: text/plain; charset=utf-8
X-Ratelimit-Current: 61
X-Ratelimit-Maximum: 60
Date: Sun, 21 Mar 2021 16:10:53 GMT
Content-Length: 5

Error
```


### 執行測試
```
make test
```
