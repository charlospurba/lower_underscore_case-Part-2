package services

import (
	"gin-user-app/database"
	"gin-user-app/dto"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	os.Setenv("TEST_ENV", "true")
	db := database.InitDB()
	if db == nil {
		t.Fatalf("Failed to initialize test database")
	}

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// TestCreateUser_ValidInputIntegration menguji pembuatan pengguna dengan input yang valid.
func TestCreateUser_ValidInputIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)
	assert.Equal(t, createUserDTO.Username, user.Username)
	assert.Equal(t, createUserDTO.Email, user.Email)
	assert.Equal(t, createUserDTO.FirstName, user.FirstName)
	assert.Equal(t, createUserDTO.LastName, user.LastName)
	assert.Equal(t, createUserDTO.Age, user.Age)
	assert.NotZero(t, user.ID)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)

	tx.Rollback()
}

// TestCreateUserIntegration_ShortUsername menguji pembuatan pengguna dengan username terlalu pendek.
func TestCreateUserIntegration_ShortUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "ab",
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

	tx.Rollback()
}

// TestCreateUserIntegration_NonGmailEmail menguji pembuatan pengguna dengan email yang tidak berakhiran @gmail.com.
func TestCreateUserIntegration_NonGmailEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
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

	tx.Rollback()
}

// TestCreateUserIntegration_ShortPassword menguji pembuatan pengguna dengan password terlalu pendek.
func TestCreateUserIntegration_ShortPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "pass", // Too short
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "password must be at least 8 characters", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)

	tx.Rollback()
}

// TestCreateUserIntegration_YoungAge menguji pembuatan pengguna dengan usia kurang dari atau sama dengan 15.
func TestCreateUserIntegration_YoungAge(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 14 // Too young

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

	tx.Rollback()
}

// TestCreateUserIntegration_DuplicateUsername menguji pembuatan pengguna dengan username yang sudah digunakan.
func TestCreateUserIntegration_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO1 := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test1@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}
	_, err := service.CreateUser(createUserDTO1)
	assert.NoError(t, err)

	createUserDTO2 := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test2@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}

	user, err := service.CreateUser(createUserDTO2)
	assert.Error(t, err)
	assert.Equal(t, "username already taken", err.Error())
	assert.Equal(t, dto.UserDTO{}, user)

	tx.Rollback()
}

// TestUpdateUser_ValidInputIntegration menguji pembaruan pengguna dengan input yang valid.
func TestUpdateUser_ValidInputIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}
	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)

	updateUserDTO := dto.UpdateUserDTO{
		Username:  "updateduser",
		Email:     "updatedemail@gmail.com",
		FirstName: "UpdatedTest",
		LastName:  "UpdatedUser",
		Age:       &age,
	}

	updatedUser, err := service.UpdateUser(user.ID, updateUserDTO)
	assert.NoError(t, err)

	assert.Equal(t, updateUserDTO.Username, updatedUser.Username)
	assert.Equal(t, updateUserDTO.Email, updatedUser.Email)
	assert.Equal(t, updateUserDTO.FirstName, updatedUser.FirstName)
	assert.Equal(t, updateUserDTO.LastName, updatedUser.LastName)
	assert.Equal(t, *updateUserDTO.Age, *updatedUser.Age)
	assert.NotZero(t, updatedUser.ID)
	assert.WithinDuration(t, time.Now(), updatedUser.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), updatedUser.UpdatedAt, time.Second)

	tx.Rollback()
}

// TestUpdateUser_InvalidEmailIntegration menguji pembaruan pengguna dengan email yang tidak valid.
func TestUpdateUser_InvalidEmailIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}
	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)

	updateUserDTO := dto.UpdateUserDTO{
		Username:  "updateduser",
		Email:     "invalid-email",
		FirstName: "UpdatedTest",
		LastName:  "UpdatedUser",
		Age:       &age,
	}

	updatedUser, err := service.UpdateUser(user.ID, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "invalid email format", err.Error())
	assert.Equal(t, dto.UserDTO{}, updatedUser)

	tx.Rollback()
}

// TestUpdateUserIntegration_ShortUsername menguji pembaruan pengguna dengan username terlalu pendek.
func TestUpdateUserIntegration_ShortUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}
	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)

	updateUserDTO := dto.UpdateUserDTO{
		Username:  "ab",
		Email:     "updatedemail@gmail.com",
		FirstName: "UpdatedTest",
		LastName:  "UpdatedUser",
		Age:       &age,
	}

	updatedUser, err := service.UpdateUser(user.ID, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "username must be 3-20 characters", err.Error())
	assert.Equal(t, dto.UserDTO{}, updatedUser)

	tx.Rollback()
}

// TestUpdateUserIntegration_NonGmailEmail menguji pembaruan pengguna dengan email yang tidak berakhiran @gmail.com.
func TestUpdateUserIntegration_NonGmailEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}
	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)

	updateUserDTO := dto.UpdateUserDTO{
		Username:  "updateduser",
		Email:     "test@yahoo.com", // Non-Gmail
		FirstName: "UpdatedTest",
		LastName:  "UpdatedUser",
		Age:       &age,
	}

	updatedUser, err := service.UpdateUser(user.ID, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "email must be a valid Gmail address (@gmail.com)", err.Error())
	assert.Equal(t, dto.UserDTO{}, updatedUser)

	tx.Rollback()
}

// TestUpdateUser_NotFoundIntegration menguji pembaruan pengguna yang tidak ada di database.
func TestUpdateUser_NotFoundIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&models.User{})

	tx := db.Begin()
	userRepo := repositories.NewUserRepository(tx)
	service := NewUserService(userRepo)
	age := 20

	createUserDTO := dto.CreateUserDTO{
		Username:  "testuser",
		Email:     "test@gmail.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Age:       &age,
	}
	user, err := service.CreateUser(createUserDTO)
	assert.NoError(t, err)

	nonExistentUserID := user.ID + 1
	updateUserDTO := dto.UpdateUserDTO{
		Username:  "updateduser",
		Email:     "updatedemail@gmail.com",
		FirstName: "UpdatedTest",
		LastName:  "UpdatedUser",
		Age:       &age,
	}

	updatedUser, err := service.UpdateUser(nonExistentUserID, updateUserDTO)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, dto.UserDTO{}, updatedUser)

	tx.Rollback()
}
