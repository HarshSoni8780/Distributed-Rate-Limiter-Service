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
// which user are we checking->
local key = KEYS[1]

local capacity = tonumber(argv[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

// "HMGET"->hash-map, token->remaining tokens in bucket, timestamp->when we last updated it.
local data = redis.call("HMGET", key, "token", "timestamp")
local tokens = tonumber(data[1])
local last_time = tonumber(data[2])

if token == nil then
	tokens = capacity
	time = now
end

//refill token logic
local delta = math.max(0,now-last_time)
local new_tokens = math.min(capacity, tokens + delta * rate)

if new_tokens < 1
	return 0
else
	new_tokens = new_tokens - 1
	redis.call("HMSET", key, "tokens", new_tokens, "time_stamps", now)
	redis.call("EXPIRE", key, 60)
	return 1
end
`)

func(tb *TokenBucket) Allow(userID string) bool{
	key:= "rate:tb" + userID
	now:= float64(time.Now().Unix())

	res, err := tokenBucketScript.Run(ctx, tb.Rdb, []string{key},
		tb.Capacity,
		tb.Rate,
		now,
	).Int()

	if err != nil {
		return false
	}

	return res == 1
}