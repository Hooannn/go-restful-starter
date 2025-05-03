package redis

import (
	"github.com/Hooannn/EventPlatform/configs"
	r "github.com/redis/go-redis/v9"
)

func InitRedis() (*r.Client, error) {
	cfg := configs.LoadConfig(".env")
	return r.NewClient(&r.Options{
		Addr:     cfg.RedisAddress,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	}), nil
}
