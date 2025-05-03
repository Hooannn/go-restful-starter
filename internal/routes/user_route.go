package routes

import (
	"github.com/Hooannn/EventPlatform/internal/factory"
	"github.com/Hooannn/EventPlatform/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewUserRoute(group *gin.RouterGroup, f *factory.Factory) {
	v1 := group.Group("/v1/users")
	v1.GET("/", f.UserHandler.GetAuthenticatedUser)

	v2 := group.Group("/v2/users")
	v2.GET("/", middleware.WithPermissions([]string{"read:users"}), f.UserHandler.GetAuthenticatedUser)
}
