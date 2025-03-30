package repositories

import (
	"gin-user-app/models"

	"gorm.io/gorm"
)

// AuthRepository handles data access for authentication
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// GetUserByUsername finds a user by username
func (r *AuthRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID finds a user by ID
func (r *AuthRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
