package services

import (
	"errors"
	"gin-user-app/dto"
	"gin-user-app/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUser_ValidInput(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	// Mock repository behavior
	mockRepo.On("FindByUsername", createUserDTO.Username).Return(models.User{}, nil)
	mockRepo.On("Create", mock.Anything).Return(models.User{
		ID:        1,
		Username:  createUserDTO.Username,
		Email:     createUserDTO.Email,
		FirstName: createUserDTO.FirstName,
		LastName:  createUserDTO.LastName,
		Age:       createUserDTO.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)
	assert.Equal(t, createUserDTO.Username, user.Username)
	assert.Equal(t, createUserDTO.Email, user.Email)
	assert.Equal(t, createUserDTO.FirstName, user.FirstName)
	assert.Equal(t, createUserDTO.LastName, user.LastName)
	assert.Equal(t, createUserDTO.Age, user.Age)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "ab", // Terlalu pendek
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "username must be 3-20 characters", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertNotCalled(t, "FindByUsername")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_NonGmailEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@yahoo.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "email must be a valid Gmail address (@gmail.com)", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertNotCalled(t, "FindByUsername")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_ShortPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "pass", // Terlalu pendek
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "password must be at least 8 characters", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertNotCalled(t, "FindByUsername")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_InvalidAge(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 14 // Terlalu muda

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "age must be greater than 15", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertNotCalled(t, "FindByUsername")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_UsernameTaken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	// Mock username sudah digunakan
	mockRepo.On("FindByUsername", createUserDTO.Username).Return(models.User{ID: 1}, nil)

	user, err := service.CreateUser(createUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "username already taken", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUpdateUser_ValidInput(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	age := 25

	updateUserDTO := dto.UpdateUserDTO{
		Username:  "newuser",
		Email:     "new@gmail.com",
		FirstName: "New",
		LastName:  "User",
		Age:       &age,
	}

	// Mock existing user
	existingUser := models.User{
		ID:        1,
		Username:  "olduser",
		Email:     "old@gmail.com",
		FirstName: "Old",
		LastName:  "User",
		Age:       &age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", 1).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything).Return(models.User{
		ID:        1,
		Username:  updateUserDTO.Username,
		Email:     updateUserDTO.Email,
		FirstName: updateUserDTO.FirstName,
		LastName:  updateUserDTO.LastName,
		Age:       updateUserDTO.Age,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: time.Now(),
	}, nil)

	user, err := service.UpdateUser(1, updateUserDTO)
	assert.NoError(t, err)
	assert.Equal(t, updateUserDTO.Username, user.Username)
	assert.Equal(t, updateUserDTO.Email, user.Email)
	assert.Equal(t, updateUserDTO.FirstName, user.FirstName)
	assert.Equal(t, updateUserDTO.LastName, user.LastName)
	assert.Equal(t, updateUserDTO.Age, user.Age)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_InvalidUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	updateUserDTO := dto.UpdateUserDTO{
		Username: "ab", // Terlalu pendek
	}

	user, err := service.UpdateUser(1, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "username must be 3-20 characters", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertNotCalled(t, "GetByID") // Pastikan GetByID tidak dipanggil
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUser_NonGmailEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	updateUserDTO := dto.UpdateUserDTO{
		Email: "test@yahoo.com",
	}

	user, err := service.UpdateUser(1, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "email must be a valid Gmail address (@gmail.com)", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertNotCalled(t, "GetByID") // Pastikan GetByID tidak dipanggil
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	updateUserDTO := dto.UpdateUserDTO{
		Username: "newuser",
	}

	// Mock user not found
	mockRepo.On("GetByID", 1).Return(models.User{}, errors.New("user not found"))

	user, err := service.UpdateUser(1, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}
