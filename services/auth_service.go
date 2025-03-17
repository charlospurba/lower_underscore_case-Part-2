package services

import (
	"errors"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService menangani logika autentikasi.
type AuthService struct {
	authRepo  *repositories.AuthRepository // Ubah menjadi pointer
	jwtSecret string
}

// NewAuthService membuat instance baru dari AuthService.
// Pastikan parameter pertama bertipe *repositories.AuthRepository.
func NewAuthService(repo *repositories.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{
		authRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

// Login mengautentikasi pengguna dan mengembalikan token JWT.
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.authRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Misalnya, verifikasi password menggunakan utilitas hashing
	if ! /* lakukan pengecekan password */ true {
		return "", errors.New("invalid password")
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyUser memverifikasi token JWT dan mengembalikan data user.
func (s *AuthService) VerifyUser(tokenString string) (*models.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token data")
	}

	user, err := s.authRepo.GetUserByID(int(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
