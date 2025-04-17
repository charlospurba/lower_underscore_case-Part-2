package dto

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

// CreateUserDTO digunakan untuk membuat user baru
type CreateUserDTO struct {
	Username  string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Email     string `json:"email" binding:"required,email,gmail"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required,min=3,max=20"`
	LastName  string `json:"lastName" binding:"required,min=3,max=20"`
	Age       *int   `json:"age" binding:"required,gt=15"`
}

// UserDTO digunakan untuk respons user
type UserDTO struct {
	ID        int       `json:"id"` // Ubah dari uint ke int
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Age       *int      `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UpdateUserDTO digunakan untuk memperbarui user
type UpdateUserDTO struct {
	Username  string `json:"username" binding:"omitempty,alphanum,min=3,max=20"`
	Email     string `json:"email" binding:"omitempty,email,gmail"`
	Password  string `json:"password" binding:"omitempty,min=8"`
	FirstName string `json:"firstName" binding:"omitempty,min=3,max=20"`
	LastName  string `json:"lastName" binding:"omitempty,min=3,max=20"`
	Age       *int   `json:"age" binding:"omitempty,gt=15"`
}

// RegisterCustomValidations mendaftarkan validator kustom
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("gmail", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		return strings.HasSuffix(email, "@gmail.com")
	})
}
