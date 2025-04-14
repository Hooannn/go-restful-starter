package util

import (
	"errors"
	"slices"
	"time"

	"github.com/Hooannn/EventPlatform/configs"
	"github.com/Hooannn/EventPlatform/internal/constant"
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/types"
	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user *entity.User) (string, error) {
	cfg := configs.LoadConfig(".env")
	exp := time.Now().Add(time.Second * time.Duration(cfg.JWTAccessTokenExpireTime)).Unix()

	roles := make([]string, 0)
	permissions := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)

		for _, permission := range role.Permissions {
			if !slices.Contains(permissions, permission.Name) {
				permissions = append(permissions, permission.Name)
			}
		}
	}

	claims := &types.JwtCustomClaims{
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.AppName,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{cfg.AppName},
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.JWTAccessTokenSecret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user *entity.User) (string, error) {
	cfg := configs.LoadConfig(".env")
	exp := time.Now().Add(time.Second * time.Duration(cfg.JWTRefreshTokenExpireTime)).Unix()
	claimsRefresh := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
		Issuer:    cfg.AppName,
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  []string{cfg.AppName},
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	refreshToken, err := token.SignedString([]byte(cfg.JWTRefreshTokenSecret))
	if err != nil {
		return "", err
	}
	return refreshToken, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(constant.InvalidSigningMethod)
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractToken(requestToken string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(constant.InvalidSigningMethod)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return nil, errors.New(constant.InvalidToken)
	}

	return claims, nil
}
