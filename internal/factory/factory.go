package factory

import (
	"github.com/Hooannn/EventPlatform/internal/handler"
	"github.com/Hooannn/EventPlatform/internal/repository"
	"github.com/Hooannn/EventPlatform/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler

	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewFactory(db *gorm.DB, redisClient *redis.Client) *Factory {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, redisClient)

	return &Factory{
		UserHandler: handler.NewUserHandler(userService),
		AuthHandler: handler.NewAuthHandler(authService),
		DB:          db,
		RedisClient: redisClient,
	}
}
