package service

import (
	"fmt"
	"time"

	"github.com/Hooannn/EventPlatform/configs"
	"github.com/Hooannn/EventPlatform/internal/constant"
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/repository"
	"github.com/Hooannn/EventPlatform/internal/types"
	"github.com/Hooannn/EventPlatform/internal/util"
	exception "github.com/Hooannn/EventPlatform/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    *repository.UserRepository
	RedisClient *redis.Client
}

func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{
		UserRepo:    userRepo,
		RedisClient: redisClient,
	}
}

func (s *AuthService) generateTokens(c *gin.Context, user *entity.User, deviceID string) (string, string, *exception.HttpException) {
	cfg := configs.LoadConfig(".env")
	accessToken, err := util.CreateAccessToken(user)

	if err != nil {
		return "", "", exception.NewInteralServerError(err.Error(), err)
	}

	refreshToken, err := util.CreateRefreshToken(user, deviceID)

	if err != nil {
		return "", "", exception.NewInteralServerError(err.Error(), err)
	}

	err = s.RedisClient.Set(
		c,
		fmt.Sprintf("refresh_token:user_id:%v:device_id:%v", user.ID, deviceID),
		refreshToken,
		time.Duration(cfg.JWTRefreshTokenExpireHours)*time.Hour).Err()

	if err != nil {
		return "", "", exception.NewInteralServerError(err.Error(), err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(c *gin.Context, deviceID, username, password string) (*types.LoginResponse, *exception.HttpException) {
	user, err := s.UserRepo.GetDetails(entity.User{Email: username})

	invalidException := exception.NewBadRequestException(constant.InvalidCredentials, nil)

	if err != nil {
		return nil, invalidException
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, invalidException
	}

	accessToken, refreshToken, ex := s.generateTokens(c, user, deviceID)

	if ex != nil {
		return nil, ex
	}

	return &types.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(c *gin.Context, token string) (*types.RefreshResponse, *exception.HttpException) {
	cfg := configs.LoadConfig(".env")

	isAuthorized, err := util.IsAuthorized(token, cfg.JWTRefreshTokenSecret)

	if err != nil {
		return nil, exception.NewForbiddenException(err.Error(), err)
	}

	if isAuthorized {
		claims, err := util.ExtractToken(token, cfg.JWTRefreshTokenSecret)

		if err != nil {
			return nil, exception.NewForbiddenException(err.Error(), err)
		}

		userID := claims["sub"].(string)
		deviceID := claims["device_id"].(string)

		storedToken, err := s.RedisClient.Get(c, fmt.Sprintf("refresh_token:user_id:%v:device_id:%v", userID, deviceID)).Result()

		if err != nil {
			if err == redis.Nil {
				return nil, exception.NewForbiddenException(constant.InvalidToken, err)
			}
			return nil, exception.NewInteralServerError(err.Error(), err)
		}

		if storedToken != token {
			return nil, exception.NewForbiddenException(constant.InvalidToken, err)
		}

		user, err := s.UserRepo.GetDetails(entity.User{ID: uuid.MustParse(userID)})

		if err != nil {
			return nil, exception.NewInteralServerError(err.Error(), err)
		}

		accessToken, refreshToken, ex := s.generateTokens(c, user, deviceID)

		if ex != nil {
			return nil, ex
		}

		return &types.RefreshResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}

	return nil, exception.NewForbiddenException(err.Error(), err)
}

func (s *AuthService) Logout(c *gin.Context, deviceID, accessToken string) (bool, *exception.HttpException) {
	cfg := configs.LoadConfig(".env")

	claims, err := util.ExtractToken(accessToken, cfg.JWTAccessTokenSecret)

	if err != nil {
		return false, exception.NewForbiddenException(err.Error(), err)
	}

	userID := claims["sub"].(string)

	if err := s.RedisClient.Del(c, fmt.Sprintf("refresh_token:user_id:%v:device_id:%v", userID, deviceID)).Err(); err != nil {
		return false, exception.NewInteralServerError(err.Error(), err)
	}

	return true, nil
}
