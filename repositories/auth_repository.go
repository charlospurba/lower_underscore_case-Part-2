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

func (r *AuthRepository) BlacklistToken(token string) error {
	blacklistedToken := models.BlacklistedToken{Token: token}
	return r.db.Create(&blacklistedToken).Error
}

// Cek apakah token sudah di-blacklist
func (r *AuthRepository) IsTokenBlacklisted(token string) bool {
	var count int64
	r.db.Model(&models.BlacklistedToken{}).
		Where("token = ?", token).
		Count(&count)
	return count > 0
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
