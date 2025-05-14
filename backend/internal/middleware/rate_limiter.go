package middleware

import (
	"context"
	"fmt"
	"grab-bootcamp-be-team13-2025/pkg/utils/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(redisClient *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		ctx := context.Background()

		// Get the current count
		count, err := redisClient.Get(ctx, key)
		if err == nil {
			// If count exists and exceeds limit
			if count > "100" { // 100 requests per minute
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Too many requests",
				})
				c.Abort()
				return
			}
		}

		// Increment the count
		err = redisClient.Set(ctx, key, "1", time.Minute)
		if err != nil {
			c.Next() // Continue if Redis fails
			return
		}

		c.Next()
	}
}
