package handler

import (
	"net/http"

	"github.com/Hooannn/EventPlatform/internal/service"
	api "github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) GetAuthenticatedUser(c *gin.Context) {
	userId := c.MustGet("x-user-id").(string)
	user, err := h.UserService.FindById(userId)

	if err != nil {
		err.Send(c)
	}

	res := api.NewOKResponse(http.StatusText(http.StatusOK), user)
	res.Send(c)
}
