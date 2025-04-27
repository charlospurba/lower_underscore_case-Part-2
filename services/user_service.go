package services

import (
	"errors"
	"gin-user-app/dto"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"gin-user-app/utils"
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
	// Validasi username
	if len(user.Username) < utils.UsernameMinLength || len(user.Username) > utils.UsernameMaxLength {
		return dto.UserDTO{}, errors.New("username must be 3-20 characters")
	}
	if !utils.UsernameRegex.MatchString(user.Username) {
		return dto.UserDTO{}, errors.New("username must contain only letters and numbers")
	}

	// Validasi email
	if utils.EmailRequired && user.Email == "" {
		return dto.UserDTO{}, errors.New("email is required")
	}
	if !utils.EmailRegex.MatchString(user.Email) {
		return dto.UserDTO{}, errors.New("invalid email format")
	}
	if !strings.HasSuffix(user.Email, utils.EmailGmailSuffix) {
		return dto.UserDTO{}, errors.New("email must be a valid Gmail address (@gmail.com)")
	}

	// Validasi password
	if len(user.Password) < utils.PasswordMinLength {
		return dto.UserDTO{}, errors.New("password must be at least 8 characters")
	}

	// Validasi first name
	if user.FirstName != "" {
		if len(user.FirstName) < utils.NameMinLength || len(user.FirstName) > utils.NameMaxLength {
			return dto.UserDTO{}, errors.New("first name must be 3-20 characters")
		}
	}

	// Validasi last name
	if user.LastName != "" {
		if len(user.LastName) < utils.NameMinLength || len(user.LastName) > utils.NameMaxLength {
			return dto.UserDTO{}, errors.New("last name must be 3-20 characters")
		}
	}

	// Validasi age
	if user.Age != nil && *user.Age <= utils.AgeMin {
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
	// Validasi username
	if user.Username != "" {
		if len(user.Username) < utils.UsernameMinLength || len(user.Username) > utils.UsernameMaxLength {
			return dto.UserDTO{}, errors.New("username must be 3-20 characters")
		}
		if !utils.UsernameRegex.MatchString(user.Username) {
			return dto.UserDTO{}, errors.New("username must contain only letters and numbers")
		}
	}

	// Validasi email
	if user.Email != "" {
		if !utils.EmailRegex.MatchString(user.Email) {
			return dto.UserDTO{}, errors.New("invalid email format")
		}
		if !strings.HasSuffix(user.Email, utils.EmailGmailSuffix) {
			return dto.UserDTO{}, errors.New("email must be a valid Gmail address (@gmail.com)")
		}
	}

	// Validasi password
	if user.Password != "" {
		if len(user.Password) < utils.PasswordMinLength {
			return dto.UserDTO{}, errors.New("password must be at least 8 characters")
		}
	}

	// Validasi first name
	if user.FirstName != "" {
		if len(user.FirstName) < utils.NameMinLength || len(user.FirstName) > utils.NameMaxLength {
			return dto.UserDTO{}, errors.New("first name must be 3-20 characters")
		}
	}

	// Validasi last name
	if user.LastName != "" {
		if len(user.LastName) < utils.NameMinLength || len(user.LastName) > utils.NameMaxLength {
			return dto.UserDTO{}, errors.New("last name must be 3-20 characters")
		}
	}

	// Validasi age
	if user.Age != nil && *user.Age <= utils.AgeMin {
		return dto.UserDTO{}, errors.New("age must be greater than 15")
	}

	// Ambil user yang ada
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	// Update field yang diisi
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
