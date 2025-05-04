package service

import (
	"fmt"
	"log"
	"net/smtp"
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
	cfg := configs.LoadConfig()
	accessToken, err := util.CreateAccessToken(user)

	if err != nil {
		return "", "", exception.NewInteralServerError(err.Error(), err)
	}

	refreshToken, err := util.CreateRefreshToken(user, deviceID)

	if err != nil {
		return "", "", exception.NewInteralServerError(err.Error(), err)
	}

	if err := s.RedisClient.Set(
		c,
		fmt.Sprintf("%s:user_id:%v:device_id:%s", constant.RefreshTokenKeyPrefix, user.ID, deviceID),
		refreshToken,
		time.Duration(cfg.JWTRefreshTokenExpireHours)*time.Hour).Err(); err != nil {
		return "", "", exception.NewInteralServerError(err.Error(), err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(c *gin.Context, deviceID string, request types.LoginRequest) (*types.LoginResponse, *exception.HttpException) {
	user, err := s.UserRepo.GetDetails(entity.User{Email: request.Username})

	invalidException := exception.NewBadRequestException(constant.InvalidCredentials, nil)

	if err != nil {
		return nil, invalidException
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
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

func (s *AuthService) Refresh(c *gin.Context, request types.RefreshRequest) (*types.RefreshResponse, *exception.HttpException) {
	cfg := configs.LoadConfig()

	token := request.RefreshToken

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

		storedToken, err := s.RedisClient.Get(c, fmt.Sprintf("%s:user_id:%v:device_id:%s", constant.RefreshTokenKeyPrefix, userID, deviceID)).Result()

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
	cfg := configs.LoadConfig()

	claims, err := util.ExtractToken(accessToken, cfg.JWTAccessTokenSecret)

	if err != nil {
		return false, exception.NewForbiddenException(err.Error(), err)
	}

	userID := claims["sub"].(string)

	if err := s.RedisClient.Del(c, fmt.Sprintf("%s:user_id:%v:device_id:%s", constant.RefreshTokenKeyPrefix, userID, deviceID)).Err(); err != nil {
		return false, exception.NewInteralServerError(err.Error(), err)
	}

	return true, nil
}

func (s *AuthService) ForgotPasswordOTP(c *gin.Context, request types.ForgotPasswordOTPRequest) (bool, *exception.HttpException) {
	cfg := configs.LoadConfig()
	username := request.Username
	if exists := s.UserRepo.ExistsByUsername(username); !exists {
		return false, exception.NewBadRequestException(constant.InvalidCredentials, nil)
	}

	otp := util.GenerateOTP()
	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)

	if err != nil {
		return false, exception.NewInteralServerError(err.Error(), err)
	}

	if err := s.RedisClient.Set(c, fmt.Sprintf("%s:username:%s", constant.ResetPasswordOTPKeyPrefix, username), string(hashed), time.Duration(cfg.ResetPasswordOTPExpireMinutes)*time.Minute).Err(); err != nil {
		return false, exception.NewInteralServerError(err.Error(), err)
	}

	go s.sendForgotPasswordOTPMail(username, otp)

	return true, nil
}

func (s *AuthService) ResetPasswordOTP(c *gin.Context, request types.ResetPasswordOTPRequest) (bool, *exception.HttpException) {
	username := request.Username
	otp := request.OTP
	password := request.Password

	if exists := s.UserRepo.ExistsByUsername(username); !exists {
		return false, exception.NewBadRequestException(constant.InvalidCredentials, nil)
	}

	storedOTP, err := s.RedisClient.Get(c, fmt.Sprintf("%s:username:%s", constant.ResetPasswordOTPKeyPrefix, username)).Result()

	if err != nil {
		if err == redis.Nil {
			return false, exception.NewBadRequestException(constant.InvalidCredentials, nil)
		}
		return false, exception.NewInteralServerError(err.Error(), err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedOTP), []byte(otp)); err != nil {
		return false, exception.NewBadRequestException(constant.InvalidCredentials, nil)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return false, exception.NewInteralServerError(err.Error(), err)
	}

	if err := s.RedisClient.Del(c, fmt.Sprintf("%s:username:%s", constant.ResetPasswordOTPKeyPrefix, username)).Err(); err != nil {
		return false, exception.NewInteralServerError(err.Error(), err)
	} else {
		success := s.UserRepo.UpdatePasswordByUsername(username, string(hashed))
		return success, nil
	}
}

func (s *AuthService) sendForgotPasswordOTPMail(username, otp string) {
	cfg := configs.LoadConfig()
	subject := "Forgot Password OTP"
	body := fmt.Sprintf("Your OTP is: %s", otp)
	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPassword, cfg.SMTPHost)

	if err := smtp.SendMail(fmt.Sprintf("%s:%v", cfg.SMTPHost, cfg.SMTPPort), auth, cfg.EmailSender, []string{username}, message); err != nil {
		log.Println("Failed to send forgot password otp email to", username, "err:", err)
	} else {
		log.Println("Forgot password otp email sent to", username)
	}
}
