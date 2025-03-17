package routes

import (
	"gin-user-app/handlers"
	"github.com/gin-gonic/gin"
)

// UserRoutes mengatur rute user dengan middleware autentikasi
func UserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, authMiddleware gin.HandlerFunc) {
	userGroup := router.Group("/users")
	userGroup.Use(authMiddleware) // Gunakan middleware auth

	userGroup.GET("/", userHandler.GetUsers)
	userGroup.GET("/:id", userHandler.GetUserByID)
	userGroup.POST("/", userHandler.CreateUser)
	userGroup.PUT("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)
}
