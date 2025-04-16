package dto

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// UserDTO digunakan untuk mengirimkan data user yang tidak sensitif
type UserDTO struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       *int      `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserDTO digunakan untuk menangani data yang diterima saat membuat user baru
type CreateUserDTO struct {
	Username  string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Email     string `json:"email" binding:"required,email,gmail"` // Validasi gmail
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"omitempty,min=3,max=20"`
	LastName  string `json:"last_name" binding:"omitempty,min=3,max=20"`
	Age       *int   `json:"age" binding:"omitempty,gt=15"`
}

// UpdateUserDTO digunakan untuk menangani data yang diterima saat memperbarui user
type UpdateUserDTO struct {
	Username  string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Email     string `json:"email" binding:"required,email,gmail"` // Validasi gmail
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"omitempty,min=3,max=20"`
	LastName  string `json:"last_name" binding:"omitempty,min=3,max=20"`
	Age       *int   `json:"age" binding:"omitempty,gt=15"`
}

// RegisterCustomValidations mendaftarkan validator kustom
func RegisterCustomValidations(validate *validator.Validate) {
	validate.RegisterValidation("gmail", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		return strings.HasSuffix(email, "@gmail.com")
	})
}
