package main

import (
	"gin-user-app/config"
	"gin-user-app/database"
	"gin-user-app/dto"
	"gin-user-app/handlers"
	"gin-user-app/middleware"
	"gin-user-app/models"
	"gin-user-app/repositories"
	"gin-user-app/routes"
	"gin-user-app/services"
	"log"

	_ "gin-user-app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm" // Tambahkan import ini
)

// @title           Gin User App API
// @version         1.0
// @description     API untuk mengelola pengguna dengan Gin Framework.
// @host           localhost:8080
// @BasePath       /
func main() {
	// Load configuration
	log.Println("Reading .env file...")
	config.LoadConfig()

	// Initialize database
	db := database.InitDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// Run database migrations
	log.Println("Running database migrations...")
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations completed")

	createDummyUser(db)

	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)

	authService := services.NewAuthService(authRepo, config.AppConfig.JWTSecret)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.AuthRoutes(r, authHandler, middleware.AuthMiddleware(config.AppConfig.JWTSecret))
	routes.UserRouter(r, userHandler, middleware.AuthMiddleware(config.AppConfig.JWTSecret))

	log.Println("Starting server on :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// createDummyUser creates a dummy user if it doesn't already exist.
func createDummyUser(db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	age := 25
	dummyUserDTO := dto.CreateUserDTO{
		Username:  "dummyuser",
		Email:     "dummyuser@gmail.com",
		Password:  "dummypassword123",
		FirstName: "Dummy",
		LastName:  "User",
		Age:       &age,
	}

	existingUser, err := userRepo.FindByUsername(dummyUserDTO.Username)
	if err == nil && existingUser.ID != 0 {
		log.Println("Dummy user already exists:", dummyUserDTO.Username)
		return
	}

	_, err = userService.CreateUser(dummyUserDTO)
	if err != nil {
		log.Println("Failed to create dummy user:", err)
		return
	}
	log.Println("Dummy user created successfully:", dummyUserDTO.Username)
}
