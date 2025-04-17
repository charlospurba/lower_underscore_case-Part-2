package services

import (
	"errors"
	"gin-user-app/dto"
	"gin-user-app/models"
	"gin-user-app/utils"
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

// intPtr adalah helper untuk membuat pointer ke int
func intPtr(i int) *int {
	return &i
}

// TestCreateUser menguji fungsi CreateUser
func TestCreateUser(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Valid user data
	validUser := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       intPtr(20),
	}

	// Mock hash password
	hashedPassword, _ := utils.HashPassword("password123")

	// Test cases
	tests := []struct {
		name          string
		input         dto.CreateUserDTO
		mockSetup     func()
		expectedError string
		expectSuccess bool
	}{
		{
			name:  "Valid user creation",
			input: validUser,
			mockSetup: func() {
				mockRepo.On("FindByUsername", validUser.Username).Return(models.User{}, nil).Once()
				mockRepo.On("Create", mock.AnythingOfType("models.User")).Return(models.User{
					ID:        1,
					Username:  validUser.Username,
					Email:     validUser.Email,
					Password:  hashedPassword,
					FirstName: validUser.FirstName,
					LastName:  validUser.LastName,
					Age:       validUser.Age,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: "",
			expectSuccess: true,
		},
		{
			name: "Username too short",
			input: dto.CreateUserDTO{
				Username:  "ab",
				Email:     "test@gmail.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       intPtr(20),
			},
			mockSetup:     func() {},
			expectedError: "username must be 3-20 characters",
			expectSuccess: false,
		},
		{
			name: "Username with invalid characters",
			input: dto.CreateUserDTO{
				Username:  "test@user",
				Email:     "test@gmail.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       intPtr(20),
			},
			mockSetup:     func() {},
			expectedError: "username must contain only letters and numbers",
			expectSuccess: false,
		},
		{
			name: "Invalid email (not @gmail.com)",
			input: dto.CreateUserDTO{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       intPtr(20),
			},
			mockSetup:     func() {},
			expectedError: "email must be a valid Gmail address (@gmail.com)",
			expectSuccess: false,
		},
		{
			name: "Password too short",
			input: dto.CreateUserDTO{
				Username:  "testuser",
				Email:     "test@gmail.com",
				Password:  "pass",
				FirstName: "Test",
				LastName:  "User",
				Age:       intPtr(20),
			},
			mockSetup:     func() {},
			expectedError: "password must be at least 8 characters",
			expectSuccess: false,
		},
		{
			name: "Age too young",
			input: dto.CreateUserDTO{
				Username:  "testuser",
				Email:     "test@gmail.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				Age:       intPtr(14),
			},
			mockSetup:     func() {},
			expectedError: "age must be greater than 15",
			expectSuccess: false,
		},
		{
			name: "FirstName too short",
			input: dto.CreateUserDTO{
				Username:  "testuser",
				Email:     "test@gmail.com",
				Password:  "password123",
				FirstName: "Te",
				LastName:  "User",
				Age:       intPtr(20),
			},
			mockSetup:     func() {},
			expectedError: "first name must be 3-20 characters",
			expectSuccess: false,
		},
		{
			name: "LastName too short",
			input: dto.CreateUserDTO{
				Username:  "testuser",
				Email:     "test@gmail.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "Us",
				Age:       intPtr(20),
			},
			mockSetup:     func() {},
			expectedError: "last name must be 3-20 characters",
			expectSuccess: false,
		},
		{
			name:  "Username already taken",
			input: validUser,
			mockSetup: func() {
				mockRepo.On("FindByUsername", validUser.Username).Return(models.User{ID: 1}, nil).Once()
			},
			expectedError: "username already taken",
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup()

			// Panggil fungsi CreateUser
			result, err := service.CreateUser(tt.input)

			// Periksa hasil
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Username, result.Username)
				assert.Equal(t, tt.input.Email, result.Email)
				assert.Equal(t, tt.input.FirstName, result.FirstName)
				assert.Equal(t, tt.input.LastName, result.LastName)
				if tt.input.Age != nil {
					assert.Equal(t, *tt.input.Age, *result.Age)
				}
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Equal(t, dto.UserDTO{}, result)
			}

			// Pastikan semua mock dipanggil
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateUser menguji fungsi UpdateUser
func TestUpdateUser(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Existing user data
	existingUser := models.User{
		ID:        1,
		Username:  "olduser",
		Email:     "old@gmail.com",
		Password:  "hashedpassword",
		FirstName: "Old",
		LastName:  "User",
		Age:       intPtr(25),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Valid update data
	validUpdate := dto.UpdateUserDTO{
		Username:  "newuser",
		Email:     "new@gmail.com",
		Password:  "newpassword123",
		FirstName: "New",
		LastName:  "User",
		Age:       intPtr(30),
	}

	// Mock hash password
	hashedPassword, _ := utils.HashPassword("newpassword123")

	// Test cases
	tests := []struct {
		name          string
		id            int
		input         dto.UpdateUserDTO
		mockSetup     func()
		expectedError string
		expectSuccess bool
	}{
		{
			name:  "Valid user update",
			id:    1,
			input: validUpdate,
			mockSetup: func() {
				mockRepo.On("GetByID", 1).Return(existingUser, nil).Once()
				mockRepo.On("Update", mock.AnythingOfType("models.User")).Return(models.User{
					ID:        1,
					Username:  validUpdate.Username,
					Email:     validUpdate.Email,
					Password:  hashedPassword,
					FirstName: validUpdate.FirstName,
					LastName:  validUpdate.LastName,
					Age:       validUpdate.Age,
					CreatedAt: existingUser.CreatedAt,
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: "",
			expectSuccess: true,
		},
		{
			name:  "User not found",
			id:    999,
			input: validUpdate,
			mockSetup: func() {
				mockRepo.On("GetByID", 999).Return(models.User{}, errors.New("user not found")).Once()
			},
			expectedError: "user not found",
			expectSuccess: false,
		},
		{
			name: "Username too short",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "ab",
				Email:     "new@gmail.com",
				Password:  "newpassword123",
				FirstName: "New",
				LastName:  "User",
				Age:       intPtr(30),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "username must be 3-20 characters",
			expectSuccess: false,
		},
		{
			name: "Username with invalid characters",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "new@user",
				Email:     "new@gmail.com",
				Password:  "newpassword123",
				FirstName: "New",
				LastName:  "User",
				Age:       intPtr(30),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "username must contain only letters and numbers",
			expectSuccess: false,
		},
		{
			name: "Invalid email (not @gmail.com)",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "newuser",
				Email:     "new@example.com",
				Password:  "newpassword123",
				FirstName: "New",
				LastName:  "User",
				Age:       intPtr(30),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "email must be a valid Gmail address (@gmail.com)",
			expectSuccess: false,
		},
		{
			name: "Password too short",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "newuser",
				Email:     "new@gmail.com",
				Password:  "pass",
				FirstName: "New",
				LastName:  "User",
				Age:       intPtr(30),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "password must be at least 8 characters",
			expectSuccess: false,
		},
		{
			name: "Age too young",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "newuser",
				Email:     "new@gmail.com",
				Password:  "newpassword123",
				FirstName: "New",
				LastName:  "User",
				Age:       intPtr(14),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "age must be greater than 15",
			expectSuccess: false,
		},
		{
			name: "FirstName too short",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "newuser",
				Email:     "new@gmail.com",
				Password:  "newpassword123",
				FirstName: "Ne",
				LastName:  "User",
				Age:       intPtr(30),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "first name must be 3-20 characters",
			expectSuccess: false,
		},
		{
			name: "LastName too short",
			id:   1,
			input: dto.UpdateUserDTO{
				Username:  "newuser",
				Email:     "new@gmail.com",
				Password:  "newpassword123",
				FirstName: "New",
				LastName:  "Us",
				Age:       intPtr(30),
			},
			mockSetup:     func() {}, // Tidak perlu mock GetByID
			expectedError: "last name must be 3-20 characters",
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup()

			// Panggil fungsi UpdateUser
			result, err := service.UpdateUser(tt.id, tt.input)

			// Periksa hasil
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Username, result.Username)
				assert.Equal(t, tt.input.Email, result.Email)
				assert.Equal(t, tt.input.FirstName, result.FirstName)
				assert.Equal(t, tt.input.LastName, result.LastName)
				if tt.input.Age != nil {
					assert.Equal(t, *tt.input.Age, *result.Age)
				}
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Equal(t, dto.UserDTO{}, result)
			}

			// Pastikan semua mock dipanggil
			mockRepo.AssertExpectations(t)
		})
	}
}
