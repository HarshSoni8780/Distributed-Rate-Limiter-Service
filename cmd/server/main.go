package main

import (

	"github.com/gin-gonic/gin"

	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/store"
)

func main(){
	r := gin.Default()

	//redis
	rdb := store.NewRedis()

	//limiter(100req/min)
	fw := limiter.NewTokenBucket(rdb,100,10)

	//apply middlware
	r.Use(middleware.TokenBucketLimit(fw))

	//test route
	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	r.Run(":8080")
}