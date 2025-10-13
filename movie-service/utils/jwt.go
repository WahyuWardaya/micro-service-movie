package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

func ParseToken(tokenStr string) (*jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = ""
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, ErrInvalidToken
}
