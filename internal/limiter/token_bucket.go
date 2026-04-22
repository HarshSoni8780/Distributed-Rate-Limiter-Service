package limiter

import(
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type TokenBucket struct{
	Rdb *redis.Client
	
	//capacity of bucket
	Capacity int

	//rate at which tokens are filled
	Rate float64
}

var ctx = context.Background()

//helper function to ready-made bucket each time it is called
func NewTokenBucket(rdb *redis.Client, capacity int, rate float64) *TokenBucket{
	return &TokenBucket{
		Rdb: rdb,
		Capacity: capacity,
		Rate: rate,
	}
}

//lua-script
var tokenBucketScript = redis.NewScript(`
--which user are we checking->
local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

--"HMGET"->hash-map, token->remaining tokens in bucket, timestamp->when we last updated it.
local data = redis.call("HMGET", key, "tokens", "timestamp")
local tokens = tonumber(data[1])
local last_time = tonumber(data[2])

if tokens == nil or last_time == nil then
	tokens = capacity
	last_time = now
end

--refill token logic
local delta = math.max(0,now-last_time)
local new_tokens = math.min(capacity, tokens + delta * rate)

local allowed = 0
if new_tokens >= 1 then
	new_tokens = new_tokens - 1
	allowed = 1
end

	redis.call("HMSET", key, "tokens", new_tokens, "timestamp", now)
	redis.call("EXPIRE", key, 120)

-- calculate reset time
local reset = now
if new_tokens < capacity then
	reset = now + (1 / rate)
end

return {allowed, new_tokens, reset}
`)
type Result struct {
	Allowed   bool
	Remaining int
	Reset     int64
}

func(tb *TokenBucket) Allow(userID string) Result{
	key := "rate:tb:" + userID
	now := float64(time.Now().UnixNano()) / 1e9

	res, err := tokenBucketScript.Run(ctx, tb.Rdb, []string{key},
		tb.Capacity,
		tb.Rate,
		now,
	).Result()

	if err != nil {
		return Result{Allowed: false}
	}
	values := res.([]interface{})

	allowed := values[0].(int64) == 1
	remaining := int(values[1].(int64))
	reset := int64(values[2].(float64))

	return Result{
		Allowed:   allowed,
		Remaining: remaining,
		Reset:     reset,
	}
}