package handler

import (
	"net/http"

	"github.com/Hooannn/go-restful-starter/internal/service"
	"github.com/Hooannn/go-restful-starter/internal/types"
	api "github.com/Hooannn/go-restful-starter/pkg/api"
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
	user, err := h.UserService.GetById(userId)

	if err != nil {
		err.Send(c)
	}

	res := api.NewOKResponse(http.StatusText(http.StatusOK), user)
	res.Send(c)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserService.GetAll()

	if err != nil {
		err.Send(c)
	}

	res := api.NewOKResponse(http.StatusText(http.StatusOK), users)
	res.Send(c)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request types.CreateUserRequest

	if err := c.ShouldBind(&request); err != nil {
		err := api.NewBadRequestException(err.Error(), nil)
		err.Send(c)
		return
	}

	user, ex := h.UserService.CreateUser(request)

	if ex != nil {
		ex.Send(c)
	}

	res := api.NewCreatedResponse(http.StatusText(http.StatusCreated), user)
	res.Send(c)
}
