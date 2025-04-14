package factory

import (
	"github.com/Hooannn/EventPlatform/internal/handler"
	"github.com/Hooannn/EventPlatform/internal/repository"
	"github.com/Hooannn/EventPlatform/internal/service"
	"gorm.io/gorm"
)

type Factory struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func NewFactory(db *gorm.DB) *Factory {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)

	return &Factory{
		UserHandler: handler.NewUserHandler(userService),
		AuthHandler: handler.NewAuthHandler(authService),
	}
}
