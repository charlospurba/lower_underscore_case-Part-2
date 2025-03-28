package services

import (
	"errors"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic.
type AuthService struct {
	authRepo  *repositories.AuthRepository
	jwtSecret string
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo *repositories.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{
		authRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.authRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT token
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

func (s *AuthService) Logout(token string) error {
	return s.authRepo.BlacklistToken(token)
}

// VerifyUser verifies a JWT token and returns user data.
func (s *AuthService) VerifyUser(tokenString string) (*models.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	// Token invalid or error
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token data")
	}

	// Retrieve user from repository
	user, err := s.authRepo.GetUserByID(int(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
