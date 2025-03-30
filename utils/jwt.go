package utils

import (
	"fmt"
	"time"

	"gin-user-app/config"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat JWT token
func GenerateToken(userID int) (string, error) {
	jwtSecret := []byte(config.AppConfig.JWTSecret)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken hanya untuk parsing token, bukan buat handle request HTTP
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	jwtSecret := []byte(config.AppConfig.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
