package service

import (
	"github.com/Hooannn/EventPlatform/internal/constant"
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/repository"
	"github.com/Hooannn/EventPlatform/internal/types"
	"github.com/Hooannn/EventPlatform/internal/util"
	exception "github.com/Hooannn/EventPlatform/pkg/api"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (s *AuthService) Login(username, password string) (*types.LoginResponse, *exception.HttpException) {
	user, err := s.UserRepo.GetDetails(entity.User{Email: username})

	invalidException := exception.NewBadRequestException(constant.InvalidCredentials, nil)

	if err != nil {
		return nil, invalidException
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, invalidException
	}

	accessToken, err := util.CreateAccessToken(user)

	if err != nil {
		return nil, exception.NewInteralServerError(err.Error(), err)
	}

	refreshToken, err := util.CreateRefreshToken(user)
	if err != nil {
		return nil, exception.NewInteralServerError(err.Error(), err)
	}

	//Save refresh token to redis

	return &types.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
