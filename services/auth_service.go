package services

import (
	"errors"
	"time"

	"gin-user-app/config"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"gin-user-app/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(username, password string) (string, error)
	VerifyUser(tokenString string) (*models.User, error)
}

type authService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{authRepo: repo}
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.authRepo.FindUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verifikasi password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// Generate token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token berlaku 24 jam
	})

	// Tanda tangani token
	secret := config.AppConfig.JWTSecret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) VerifyUser(tokenString string) (*models.User, error) {
	// Verifikasi token
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Ambil user ID dari token
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Cari user berdasarkan ID
	user, err := s.authRepo.FindUserByID(int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
