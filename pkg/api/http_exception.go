package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *HttpException) Error() string {
	return e.Message
}

func (e *HttpException) Send(c *gin.Context) {
	c.JSON(e.Code, e)
}

func NewHttpException(code int, message string, data any) *HttpException {
	return &HttpException{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewBadRequestException(message string, data any) *HttpException {
	return NewHttpException(http.StatusBadRequest, message, data)
}

func NewForbiddenException(message string, data any) *HttpException {
	return NewHttpException(http.StatusForbidden, message, data)
}

func NewBadGatewayException(message string, data any) *HttpException {
	return NewHttpException(http.StatusBadGateway, message, data)
}

func NewInteralServerError(message string, data any) *HttpException {
	return NewHttpException(http.StatusInternalServerError, message, data)
}

func NewNotFoundException(message string, data any) *HttpException {
	return NewHttpException(http.StatusNotFound, message, data)
}

func NewUnauthorizedException(message string, data any) *HttpException {
	return NewHttpException(http.StatusUnauthorized, message, data)
}
