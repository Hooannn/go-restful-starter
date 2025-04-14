package service

import (
	"net/http"

	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/repository"
	exception "github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/google/uuid"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) FindById(userId string) (*entity.User, *exception.HttpException) {
	user, err := s.UserRepo.GetDetails(entity.User{ID: uuid.MustParse(userId)})
	if err != nil {
		return nil, exception.NewNotFoundException(http.StatusText(http.StatusNotFound), err)
	}
	return user, nil
}
