package routes

import (
	"github.com/Hooannn/go-restful-starter/internal/factory"
	"github.com/Hooannn/go-restful-starter/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, f *factory.Factory) {
	publicRouter := router.Group("")
	NewAuthRoute(publicRouter, f)

	protectedRouter := router.Group("")
	protectedRouter.Use(middleware.WithJwtAuthMiddleware())
	NewUserRoute(protectedRouter, f)
}
