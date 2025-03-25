package repositories

import (
	"errors"
	"gin-user-app/models"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan operasi database untuk User
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) error
	FindByUsername(username string) (models.User, error) 
}

// userRepositoryImpl implementasi dari UserRepository
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// GetAll mengambil semua user dari database
func (r *userRepositoryImpl) GetAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetByID mengambil user berdasarkan ID
func (r *userRepositoryImpl) GetByID(id int) (models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, errors.New("user not found")
	}
	return user, result.Error
}

// FindByUsername mencari user berdasarkan username
func (r *userRepositoryImpl) FindByUsername(username string) (models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, errors.New("user not found")
	}
	return user, result.Error
}

// Create menambahkan user baru ke database
func (r *userRepositoryImpl) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

// Update memperbarui data user yang ada di database
func (r *userRepositoryImpl) Update(user models.User) (models.User, error) {
	result := r.db.Save(&user)
	return user, result.Error
}

// Delete menghapus user berdasarkan ID
func (r *userRepositoryImpl) Delete(id int) error {
	result := r.db.Delete(&models.User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}
