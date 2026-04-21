package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SlidingWindow struct{
	Rdb *redis.Client
	Limit int
	Window time.Duration
}

var ctx = context.Background()

func NewSlidingWindow(rdb *redis.Client, limit int, window time.Duration) *SlidingWindow{
	return &SlidingWindow{
		Rdb: rdb,
		Limit: limit,
		Window: window,
	}
}

func (sw *SlidingWindow) Allow(userID string) bool{
	ctx := context.Background()
	now := time.Now()
	key := fmt.Sprintf("rate:sw:%s",userID)

	sw.Rdb.ZRemRangeByScore(ctx,key,"0", fmt.Sprintf("%d", now.Add(-sw.Window).UnixNano()))

	count, err := sw.Rdb.ZCard(ctx,key).Result()
	if err != nil{
		return false;
	}

	if count >= int64(sw.Limit){
		return false;
	}

	sw.Rdb.ZAdd(ctx,key,redis.Z{
		Score: float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	sw.Rdb.Expire(ctx,key,sw.Window)

	return true
}