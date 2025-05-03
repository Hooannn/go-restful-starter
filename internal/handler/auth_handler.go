package handler

import (
	"github.com/Hooannn/EventPlatform/internal/constant"
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
	deviceID := c.GetHeader("x-device-id")

	var request types.LoginRequest

	err := c.ShouldBind(&request)

	if err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.Login(c, deviceID, request.Username, request.Password)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.LoginSuccess, data)
	res.Send(c)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var request types.RefreshRequest

	err := c.ShouldBind(&request)

	if err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.Refresh(c, request.RefreshToken)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.RefreshSuccess, data)
	res.Send(c)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	deviceID := c.GetHeader("x-device-id")
	accessToken := c.GetHeader("Authorization")
	accessToken = accessToken[len("Bearer "):]

	if accessToken == "" {
		err := api.NewBadRequestException(constant.MissingToken, nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.Logout(c, deviceID, accessToken)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.LogoutSuccess, &data)
	res.Send(c)
}
