package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env") 
	viper.AutomaticEnv()       

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, trying to read from environment variables")
	}

	// Debugging - cek apakah file terbaca
	fmt.Println("Reading .env file...")

	AppConfig = Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
	}

	// Debugging - tampilkan hasil
	fmt.Println("Loaded Config:")
	fmt.Println("DB_HOST:", AppConfig.DBHost)
	fmt.Println("DB_PORT:", AppConfig.DBPort)
	fmt.Println("DB_USER:", AppConfig.DBUser)
	fmt.Println("DB_PASSWORD:", AppConfig.DBPassword)
	fmt.Println("DB_NAME:", AppConfig.DBName)
	fmt.Println("JWT_SECRET:", AppConfig.JWTSecret)

	// Cek apakah JWT_SECRET kosong
	if AppConfig.JWTSecret == "" {
		log.Fatalf("ERROR: JWT_SECRET is missing or empty!")
	}

	// Set env variables (optional)
	os.Setenv("JWT_SECRET", AppConfig.JWTSecret)
}
