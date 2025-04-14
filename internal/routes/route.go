package routes

import (
	"github.com/Hooannn/EventPlatform/internal/factory"
	"github.com/Hooannn/EventPlatform/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, f *factory.Factory) {
	publicRouter := router.Group("")
	NewAuthRoute(publicRouter, f)

	protectedRouter := router.Group("")
	protectedRouter.Use(middleware.WithJwtAuthMiddleware())
	NewUserRoute(protectedRouter, f)
}
