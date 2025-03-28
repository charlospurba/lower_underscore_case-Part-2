package models

import "gorm.io/gorm"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BlacklistedToken struct {
	gorm.Model
	Token string `gorm:"uniqueIndex"`
}
