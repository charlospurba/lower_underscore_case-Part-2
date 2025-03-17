package dto

// LoginRequestDTO digunakan untuk menangani data login yang diterima
type LoginRequestDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
