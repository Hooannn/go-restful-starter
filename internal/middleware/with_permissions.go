package middleware

import (
	"net/http"
	"slices"

	"github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/gin-gonic/gin"
)

func WithPermissions(requiredPermissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		getPermissions, exist := c.Get("x-user-permissions")
		userPermissions := make([]string, len(getPermissions.([]interface{})))

		for i, v := range getPermissions.([]interface{}) {
			userPermissions[i] = v.(string)
		}

		forbiddenException := api.NewForbiddenException(http.StatusText(http.StatusForbidden), nil)

		if !exist {
			forbiddenException.Send(c)
			c.Abort()
			return
		}

		for _, p := range requiredPermissions {
			if !slices.Contains(userPermissions, p) {
				{
					forbiddenException.Send(c)
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
