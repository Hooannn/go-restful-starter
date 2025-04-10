package routes

import (
	"github.com/Hooannn/EventPlatform/internal/handler"
	"github.com/Hooannn/EventPlatform/internal/repository"
	"github.com/Hooannn/EventPlatform/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRoute struct {
	AuthHandler *handler.AuthHandler
}

func NewAuthRoute(gin *gin.RouterGroup, db *gorm.DB) {
	v1 := gin.Group("/v1/auth")

	userRepo := repository.NewUserRepository(db)
	service := service.NewAuthService(userRepo)

	h := handler.NewAuthHandler(service)

	v1.POST("/login", h.Login)
}
