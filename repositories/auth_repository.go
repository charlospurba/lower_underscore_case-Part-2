package repositories

import (
	"gin-user-app/models"

	"gorm.io/gorm"
)

// AuthRepository menangani akses data untuk autentikasi
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository membuat instance baru dari AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// GetUserByUsername mencari user berdasarkan username
func (r *AuthRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID mencari user berdasarkan ID
func (r *AuthRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
