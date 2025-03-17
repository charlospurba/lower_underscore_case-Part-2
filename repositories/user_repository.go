package repositories

import (
	"gin-user-app/models"

	"gorm.io/gorm"
)

// Definisikan interface UserRepository
type UserRepository interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

// Struct untuk implementasi UserRepository
type userRepository struct {
	db *gorm.DB
}

// Constructor untuk inisialisasi UserRepository dengan dependency injection
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Implementasi fungsi-fungsi repository
func (r *userRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}
