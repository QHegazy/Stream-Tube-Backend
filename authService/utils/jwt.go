package utils

import (
	"authService/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

type Claims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

func init() {
	
	secretKey = []byte(config.GetJWTConfig())
}

func GenerateToken[T any](data T, duration time.Duration) (string, error) {
	claims := Claims[T]{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "authService",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "usernamess",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func VerifyToken[T any](tokenString string) (*Claims[T], error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims[T]{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims[T])
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
