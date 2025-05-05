package handler

import (
	"github.com/Hooannn/go-restful-starter/internal/constant"
	"github.com/Hooannn/go-restful-starter/internal/service"
	"github.com/Hooannn/go-restful-starter/internal/types"
	api "github.com/Hooannn/go-restful-starter/pkg/api"
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
	deviceID := c.GetHeader(constant.ContextDeviceIDKey)

	if deviceID == "" {
		deviceID = "default"
	}

	var request types.LoginRequest

	if err := c.ShouldBind(&request); err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.Login(c, deviceID, request)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.LoginSuccess, data)
	res.Send(c)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var request types.RefreshRequest

	if err := c.ShouldBind(&request); err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.Refresh(c, request)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.RefreshSuccess, data)
	res.Send(c)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	deviceID := c.GetString(constant.ContextDeviceIDKey)
	accessToken := c.GetString(constant.ContextAccessTokenKey)

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

func (h *AuthHandler) ForgotPasswordOTP(c *gin.Context) {
	var request types.ForgotPasswordOTPRequest

	if err := c.ShouldBind(&request); err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.ForgotPasswordOTP(c, request)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.ForgotPasswordOTPSuccess, &data)
	res.Send(c)
}

func (h *AuthHandler) ResetPasswordOTP(c *gin.Context) {
	var request types.ResetPasswordOTPRequest

	if err := c.ShouldBind(&request); err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	data, ex := h.AuthService.ResetPasswordOTP(c, request)

	if ex != nil {
		ex.Send(c)
		return
	}

	res := api.NewOKResponse(constant.ResetPasswordOTPSuccess, &data)
	res.Send(c)
}
