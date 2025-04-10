package routes

import (
	"github.com/Hooannn/EventPlatform/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	publicRouter := router.Group("")
	NewAuthRoute(publicRouter, db)

	protectedRouter := router.Group("")
	protectedRouter.Use(middleware.WithJwtAuthMiddleware())
}
