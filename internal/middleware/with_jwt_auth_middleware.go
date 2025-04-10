package middleware

import (
	"net/http"
	"strings"

	"github.com/Hooannn/EventPlatform/configs"
	"github.com/Hooannn/EventPlatform/internal/util"
	"github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/gin-gonic/gin"
)

func WithJwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := configs.LoadConfig(".env")
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

				c.Set("x-user-id", claims["sub"])
				c.Set("x-user-roles", claims["roles"])
				c.Set("x-user-permissions", claims["permissions"])
				c.Next()
				return
			}

			api.NewUnauthorizedException(http.StatusText(http.StatusUnauthorized), nil).Send(c)
			c.Abort()
			return
		}

		api.NewUnauthorizedException(http.StatusText(http.StatusUnauthorized), nil).Send(c)
		c.Abort()
	}
}
