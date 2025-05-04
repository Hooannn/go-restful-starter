package middleware

import (
	"net/http"
	"strings"

	"github.com/Hooannn/EventPlatform/configs"
	"github.com/Hooannn/EventPlatform/internal/constant"
	"github.com/Hooannn/EventPlatform/internal/util"
	"github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/gin-gonic/gin"
)

func WithJwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := configs.LoadConfig()
		authHeader := c.Request.Header.Get("Authorization")

		if parts := strings.Split(authHeader, " "); len(parts) == 2 {
			token := parts[1]
			isAuthorized, err := util.IsAuthorized(token, cfg.JWTAccessTokenSecret)
			if err != nil {
				api.NewUnauthorizedException(err.Error(), err).Send(c)
				c.Abort()
				return
			}
			if isAuthorized {
				claims, err := util.ExtractToken(token, cfg.JWTAccessTokenSecret)
				if err != nil {
					api.NewUnauthorizedException(err.Error(), err).Send(c)
					c.Abort()
					return
				}

				deviceID := c.GetHeader("x-device-id")
				if deviceID == "" {
					deviceID = "default"
				}

				c.Set(constant.ContextUserIDKey, claims["sub"])
				c.Set(constant.ContextUserRolesKey, claims["roles"])
				c.Set(constant.ContextUserPermissionsKey, claims["permissions"])
				c.Set(constant.ContextAccessTokenKey, token)
				c.Set(constant.ContextDeviceIDKey, deviceID)
				c.Next()
				return
			}

			api.NewUnauthorizedException(http.StatusText(http.StatusUnauthorized), nil).Send(c)
			c.Abort()
			return
		}

		api.NewUnauthorizedException(constant.MissingToken, nil).Send(c)
		c.Abort()
	}
}
