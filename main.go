package main

import (
	"gin-user-app/config"
	"gin-user-app/database"
	"gin-user-app/dto"
	"gin-user-app/handlers"
	"gin-user-app/middleware"
	"gin-user-app/repositories"
	"gin-user-app/routes"
	"gin-user-app/services"

	_ "gin-user-app/docs"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin User App API
// @version         1.0
// @description     API untuk mengelola pengguna dengan Gin Framework.
// @host           localhost:8080
// @BasePath       /
func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	db := database.InitDB()

	// Inisialisasi repository
	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Inisialisasi service
	authService := services.NewAuthService(authRepo, config.AppConfig.JWTSecret)
	userService := services.NewUserService(userRepo)

	// Inisialisasi handler
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Setup router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Daftarkan validator kustom ke validator yang digunakan Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		dto.RegisterCustomValidations(v)
	}

	// Tambahkan endpoint Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	routes.AuthRoutes(r, authHandler, middleware.AuthMiddleware(config.AppConfig.JWTSecret))
	routes.UserRouter(r, userHandler, middleware.AuthMiddleware(config.AppConfig.JWTSecret))

	// Jalankan server
	r.Run(":8080")
}
