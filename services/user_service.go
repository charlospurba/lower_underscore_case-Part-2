package services

import (
	"errors"
	"gin-user-app/dto"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"gin-user-app/utils"
	"time"
)

// UserService interface untuk layanan pengguna
type UserService interface {
	GetAllUsers() ([]dto.UserDTO, error)
	GetUserByID(id int) (dto.UserDTO, error)
	CreateUser(user dto.CreateUserDTO) (dto.UserDTO, error)
	UpdateUser(id int, user dto.UpdateUserDTO) (dto.UserDTO, error)
	DeleteUser(id int) error
}

// UserServiceImpl implementasi UserService
type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService membuat instance baru UserService
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// GetAllUsers mengambil semua pengguna
func (s *UserServiceImpl) GetAllUsers() ([]dto.UserDTO, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var usersDTO []dto.UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, dto.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return usersDTO, nil
}

// GetUserByID mendapatkan user berdasarkan ID
func (s *UserServiceImpl) GetUserByID(id int) (dto.UserDTO, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, errors.New("user not found")
	}

	return dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// CreateUser membuat pengguna baru
func (s *UserServiceImpl) CreateUser(user dto.CreateUserDTO) (dto.UserDTO, error) {
	// **1. Cek apakah username sudah digunakan**
	existingUser, _ := s.userRepo.FindByUsername(user.Username)
	if existingUser.ID != 0 {
		return dto.UserDTO{}, errors.New("username already taken")
	}

	// **2. Hash password sebelum disimpan**
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return dto.UserDTO{}, err
	}

	newUser := models.User{
		Username:  user.Username,
		Password:  hashedPassword,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// **3. Simpan user ke database**
	createdUser, err := s.userRepo.Create(newUser)
	if err != nil {
		return dto.UserDTO{}, err
	}

	// **4. Return data user tanpa password**
	return dto.UserDTO{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Age:       createdUser.Age,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}

// UpdateUser memperbarui data pengguna
func (s *UserServiceImpl) UpdateUser(id int, user dto.UpdateUserDTO) (dto.UserDTO, error) {
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, errors.New("user not found")
	}

	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}
	if user.Age != nil {
		existingUser.Age = user.Age
	}

	existingUser.UpdatedAt = time.Now()

	updatedUser, err := s.userRepo.Update(existingUser)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Age:       updatedUser.Age,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

// DeleteUser menghapus pengguna berdasarkan ID
func (s *UserServiceImpl) DeleteUser(id int) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
