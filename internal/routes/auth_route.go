package routes

import (
	"github.com/Hooannn/go-restful-starter/internal/factory"
	"github.com/Hooannn/go-restful-starter/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewAuthRoute(group *gin.RouterGroup, f *factory.Factory) {
	v1 := group.Group("/v1/auth")
	v1.POST("/login", f.AuthHandler.Login)
	v1.POST("/logout", middleware.WithJwtAuthMiddleware(), f.AuthHandler.Logout)
	v1.POST("/refresh", f.AuthHandler.Refresh)
	v1.POST("/forgot-password/otp", f.AuthHandler.ForgotPasswordOTP)
	v1.POST("/reset-password/otp", f.AuthHandler.ResetPasswordOTP)
}
