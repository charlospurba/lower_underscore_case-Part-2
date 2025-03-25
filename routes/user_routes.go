package routes

import (
	"gin-user-app/handlers"

	"github.com/gin-gonic/gin"
)

// UserRouter mengatur rute pengguna
func UserRouter(r *gin.Engine, userHandler *handlers.UserHandler, authMiddleware gin.HandlerFunc) {
	// Group untuk Protected Routes
	protected := r.Group("/users")
	protected.Use(authMiddleware) 
	{
		protected.GET("/", userHandler.GetUsers)
		protected.GET("/:id", userHandler.GetUserByID)
		protected.POST("/", userHandler.CreateUser)
		protected.PUT("/:id", userHandler.UpdateUser)
		protected.DELETE("/:id", userHandler.DeleteUser)
	}
}
