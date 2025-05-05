package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Hooannn/go-restful-starter/configs"
	"github.com/Hooannn/go-restful-starter/internal/constant"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func WithCache(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := configs.LoadConfig()
		userID := c.GetString(constant.ContextUserIDKey)

		if userID == "" {
			log.Println("User ID not found in context, skipping cache")
			c.Next()
			return
		}

		cacheKey := fmt.Sprintf("%s:user_id:%s:path:%s:query:%s", constant.CacheKeyPrefix, userID, c.Request.URL.Path, c.Request.URL.RawQuery)

		cachedResponse, err := redisClient.Get(c, cacheKey).Result()
		if err == redis.Nil {
			log.Println("Cache miss... setting cache")
			c.Next()
			body, exists := c.Get(constant.ContextResponseKey)
			if exists {
				jsonBody, err := json.Marshal(body)
				if err == nil {
					if err := redisClient.Set(c, cacheKey, jsonBody, time.Duration(cfg.DefaultCacheExpireMinutes)*time.Minute).Err(); err != nil {
						log.Println("Failed to set cache with redis", err)
					}
				}
			}
		} else if err != nil {
			log.Println("Failed to get cache from redis", err)
			c.Next()
		} else {
			log.Println("Cache hit")
			c.JSON(200, json.RawMessage(cachedResponse))
			c.Abort()
			return
		}
	}
}
