package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Age       *int           `json:"age"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
