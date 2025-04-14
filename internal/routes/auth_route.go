package routes

import (
	"github.com/Hooannn/EventPlatform/internal/factory"
	"github.com/gin-gonic/gin"
)

func NewAuthRoute(group *gin.RouterGroup, f *factory.Factory) {
	v1 := group.Group("/v1/auth")
	v1.POST("/login", f.AuthHandler.Login)
}
