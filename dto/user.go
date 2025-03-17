package dto

import "time"

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
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Age       *int   `json:"age"`
}

// UpdateUserDTO digunakan untuk menangani data yang diterima saat memperbarui user
type UpdateUserDTO struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       *int   `json:"age"`
}
