package middleware

import (
	"fmt"
	"log"

	"github.com/Hooannn/go-restful-starter/internal/constant"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InvalidateCache(redisClient *redis.Client, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		userID := c.GetString(constant.ContextUserIDKey)

		if userID == "" {
			log.Println("User ID not found in context, skipping cache invalidation")
			return
		}

		pattern := fmt.Sprintf("%s:user_id:*:path:%s:query:*", constant.CacheKeyPrefix, path)

		iter := redisClient.Scan(c, 0, pattern, 0).Iterator()

		for iter.Next(c) {
			err := redisClient.Del(c, iter.Val()).Err()
			if err != nil {
				log.Printf("Failed to delete cache key %s: %v", iter.Val(), err)
			}
		}

		if err := iter.Err(); err != nil {
			log.Printf("Error iterating over cache keys: %v", err)
		} else {
			log.Println("Cache invalidated for pattern:", pattern)
		}

	}
}
