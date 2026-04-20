package store

import (
	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
