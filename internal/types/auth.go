package types

import (
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JwtAccessTokenClaims struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type JwtRefreshTokenClaims struct {
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *entity.User `json:"user"`
}
