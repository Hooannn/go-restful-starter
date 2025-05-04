package api

import (
	"net/http"

	"github.com/Hooannn/EventPlatform/internal/constant"
	"github.com/gin-gonic/gin"
)

type HttpResponse[T any] struct {
	Code    int    `json:"code"`
	Data    *T     `json:"data"`
	Message string `json:"message"`
}

func (r *HttpResponse[T]) Send(c *gin.Context) {
	c.Set(constant.ContextResponseKey, r)
	c.JSON(r.Code, r)
}

func NewHttpResponse[T any](code int, message string, data *T) *HttpResponse[T] {
	return &HttpResponse[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewOKResponse[T any](message string, data *T) *HttpResponse[T] {
	return NewHttpResponse(http.StatusOK, message, data)
}

func NewCreatedResponse[T any](message string, data *T) *HttpResponse[T] {
	return NewHttpResponse(http.StatusCreated, message, data)
}

func NewNoContentResponse[T any](message string, data *T) *HttpResponse[T] {
	return NewHttpResponse(http.StatusNoContent, message, data)
}
