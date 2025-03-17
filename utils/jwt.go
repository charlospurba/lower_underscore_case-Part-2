package utils

import (
	"errors"
	"fmt"
	"time"

	"gin-user-app/config"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat JWT token
func GenerateToken(userID int) (string, error) {
	// Ambil secret key dari config
	jwtSecret := []byte(config.AppConfig.JWTSecret)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Token berlaku 24 jam
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken memverifikasi JWT token
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Ambil secret key dari config
	jwtSecret := []byte(config.AppConfig.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
