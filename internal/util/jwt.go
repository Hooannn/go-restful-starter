package util

import (
	"errors"
	"math/rand"
	"slices"
	"strconv"
	"time"

	"github.com/Hooannn/go-restful-starter/configs"
	"github.com/Hooannn/go-restful-starter/internal/constant"
	"github.com/Hooannn/go-restful-starter/internal/entity"
	"github.com/Hooannn/go-restful-starter/internal/types"
	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user *entity.User) (string, error) {
	cfg := configs.LoadConfig()
	exp := time.Now().Add(time.Hour * time.Duration(cfg.JWTAccessTokenExpireHours)).Unix()

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

	claims := &types.JwtAccessTokenClaims{
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

func CreateRefreshToken(user *entity.User, deviceID string) (string, error) {
	cfg := configs.LoadConfig()
	exp := time.Now().Add(time.Hour * time.Duration(cfg.JWTRefreshTokenExpireHours)).Unix()

	claims := &types.JwtRefreshTokenClaims{
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
			Issuer:    cfg.AppName,
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{cfg.AppName},
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return false, errors.New(constant.InvalidTokenSignature)
		case errors.Is(err, jwt.ErrTokenMalformed):
			return false, errors.New(constant.InvalidToken)
		case errors.Is(err, jwt.ErrTokenExpired):
			return false, errors.New(constant.ExpiredToken)
		default:
			return false, err
		}
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
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return nil, errors.New(constant.InvalidTokenSignature)
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New(constant.InvalidToken)
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.New(constant.ExpiredToken)
		default:
			return nil, err
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return nil, errors.New(constant.InvalidToken)
	}

	return claims, nil
}

func GenerateOTP() string {
	otp := ""
	for range 6 {
		otp += strconv.Itoa(rand.Intn(10))
	}
	return otp
}
