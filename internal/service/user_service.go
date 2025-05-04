package service

import (
	"net/http"

	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/repository"
	"github.com/Hooannn/EventPlatform/internal/types"
	exception "github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) GetById(userId string) (*entity.User, *exception.HttpException) {
	user, err := s.UserRepo.GetDetails(entity.User{ID: uuid.MustParse(userId)})
	if err != nil {
		return nil, exception.NewNotFoundException(http.StatusText(http.StatusNotFound), err)
	}
	return user, nil
}

func (s *UserService) GetAll() (*[]entity.User, *exception.HttpException) {
	users, err := s.UserRepo.GetAll()
	if err != nil {
		return nil, exception.NewNotFoundException(http.StatusText(http.StatusNotFound), err)
	}
	return users, nil
}

func (s *UserService) CreateUser(request types.CreateUserRequest) (*entity.User, *exception.HttpException) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, exception.NewBadRequestException(err.Error(), err)
	}

	user := entity.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  string(hashed),
	}

	if err := s.UserRepo.Create(&user); err != nil {
		return nil, exception.NewBadRequestException(err.Error(), err)
	}

	return &user, nil
}
