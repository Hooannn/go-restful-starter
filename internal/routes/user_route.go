package routes

import (
	"github.com/Hooannn/go-restful-starter/internal/factory"
	"github.com/Hooannn/go-restful-starter/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewUserRoute(group *gin.RouterGroup, f *factory.Factory) {
	v1 := group.Group("/v1/users")
	v1.GET("/me", f.UserHandler.GetAuthenticatedUser)
	v1.GET("/", middleware.WithPermissions([]string{"read:users"}), middleware.WithCache(f.RedisClient), f.UserHandler.GetAllUsers)
	v1.POST("/", middleware.WithPermissions([]string{"create:users"}), f.UserHandler.CreateUser, middleware.InvalidateCache(f.RedisClient, "/v1/users/"))
}
