package redis

import (
	"github.com/Hooannn/go-restful-starter/configs"
	r "github.com/redis/go-redis/v9"
)

func InitRedis() (*r.Client, error) {
	cfg := configs.LoadConfig()
	return r.NewClient(&r.Options{
		Addr:     cfg.RedisAddress,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	}), nil
}
