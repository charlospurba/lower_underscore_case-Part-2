package repositories

import (
	"gin-user-app/models"
	"gorm.io/gorm"
)

// AuthRepository adalah interface untuk autentikasi
type AuthRepository interface {
	FindUserByUsername(username string) (*models.User, error)
	FindUserByID(userID int) (*models.User, error) // Tambahkan ini
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository membuat instance baru dari authRepository
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// FindUserByUsername mencari user berdasarkan username
func (r *authRepository) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByID mencari user berdasarkan ID
func (r *authRepository) FindUserByID(userID int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID from the database.
func (r *authRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
