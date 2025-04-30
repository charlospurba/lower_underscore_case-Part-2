package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	paths := []string{".env", "../.env", "../../.env"}

	loaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan ENV bawaan sistem.")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Konfigurasi GORM
	config := &gorm.Config{}
	if os.Getenv("TEST_ENV") == "true" {
		config.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
		return nil
	}

	DB = db
	return DB
}
