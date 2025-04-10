package handler

import (
	"net/http"

	"github.com/Hooannn/EventPlatform/internal/service"
	"github.com/Hooannn/EventPlatform/internal/types"
	api "github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request types.LoginRequest

	err := c.ShouldBind(&request)

	if err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	var data *types.LoginResponse
	data, err = h.AuthService.Login(request.Username, request.Password)

	if err := err.(*api.HttpException); err != nil {
		err.Send(c)
		return
	}

	res := api.NewOKResponse(http.StatusText(http.StatusOK), data)
	res.Send(c)
}
