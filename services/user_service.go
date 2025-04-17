package services

import (
	"errors"
	"gin-user-app/dto"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"gin-user-app/utils"
	"regexp"
	"strings"
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

	var userDTOs []dto.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dto.UserDTO{
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

	return userDTOs, nil
}

// GetUserByID mengambil pengguna berdasarkan ID
func (s *UserServiceImpl) GetUserByID(id int) (dto.UserDTO, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, err
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
	// Validasi username (3-20 karakter, hanya huruf dan angka)
	if len(user.Username) < 3 || len(user.Username) > 20 {
		return dto.UserDTO{}, errors.New("username must be 3-20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(user.Username) {
		return dto.UserDTO{}, errors.New("username must contain only letters and numbers")
	}

	// Validasi email harus @gmail.com
	if !strings.HasSuffix(user.Email, "@gmail.com") {
		return dto.UserDTO{}, errors.New("email must be a valid Gmail address (@gmail.com)")
	}

	// Validasi password (minimal 8 karakter)
	if len(user.Password) < 8 {
		return dto.UserDTO{}, errors.New("password must be at least 8 characters")
	}

	// Validasi first name (3-20 karakter)
	if len(user.FirstName) < 3 || len(user.FirstName) > 20 {
		return dto.UserDTO{}, errors.New("first name must be 3-20 characters")
	}

	// Validasi last name (3-20 karakter)
	if len(user.LastName) < 3 || len(user.LastName) > 20 {
		return dto.UserDTO{}, errors.New("last name must be 3-20 characters")
	}

	// Validasi age (> 15)
	if user.Age != nil && *user.Age <= 15 {
		return dto.UserDTO{}, errors.New("age must be greater than 15")
	}

	// Cek apakah username sudah digunakan
	existingUser, _ := s.userRepo.FindByUsername(user.Username)
	if existingUser.ID != 0 {
		return dto.UserDTO{}, errors.New("username already taken")
	}

	// Hash password sebelum disimpan
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

	// Simpan user ke database
	createdUser, err := s.userRepo.Create(newUser)
	if err != nil {
		return dto.UserDTO{}, err
	}

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

// UpdateUser memperbarui pengguna
func (s *UserServiceImpl) UpdateUser(id int, user dto.UpdateUserDTO) (dto.UserDTO, error) {
	// Validasi semua field di awal
	if user.Username != "" {
		if len(user.Username) < 3 || len(user.Username) > 20 {
			return dto.UserDTO{}, errors.New("username must be 3-20 characters")
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(user.Username) {
			return dto.UserDTO{}, errors.New("username must contain only letters and numbers")
		}
	}

	if user.Email != "" && !strings.HasSuffix(user.Email, "@gmail.com") {
		return dto.UserDTO{}, errors.New("email must be a valid Gmail address (@gmail.com)")
	}

	if user.Password != "" && len(user.Password) < 8 {
		return dto.UserDTO{}, errors.New("password must be at least 8 characters")
	}

	if user.FirstName != "" && (len(user.FirstName) < 3 || len(user.FirstName) > 20) {
		return dto.UserDTO{}, errors.New("first name must be 3-20 characters")
	}

	if user.LastName != "" && (len(user.LastName) < 3 || len(user.LastName) > 20) {
		return dto.UserDTO{}, errors.New("last name must be 3-20 characters")
	}

	if user.Age != nil && *user.Age <= 15 {
		return dto.UserDTO{}, errors.New("age must be greater than 15")
	}

	// Ambil user yang ada
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	// Update field yang diisi hanya jika validasi lulus
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return dto.UserDTO{}, err
		}
		existingUser.Password = hashedPassword
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

	// Simpan perubahan
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

// DeleteUser menghapus pengguna
func (s *UserServiceImpl) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}
