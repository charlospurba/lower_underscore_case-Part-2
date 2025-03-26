package routes

import (
	"gin-user-app/handlers"
	"github.com/gin-gonic/gin"
)

// UserRouter mengatur rute pengguna dengan proteksi middleware
func UserRouter(r *gin.Engine, userHandler *handlers.UserHandler, authMiddleware gin.HandlerFunc) {
	// Group untuk Protected Routes
	protected := r.Group("/users")
	protected.Use(authMiddleware) 

	{
		protected.GET("/", userHandler.GetUsers)         
		protected.GET("/:id", userHandler.GetUserByID)   // GET /users/:id (Butuh Auth)
		protected.POST("/", userHandler.CreateUser)      // POST /users (Butuh Auth)
		protected.PUT("/:id", userHandler.UpdateUser)    // PUT /users/:id (Butuh Auth)
		protected.DELETE("/:id", userHandler.DeleteUser) // DELETE /users/:id (Butuh Auth)
	}
}
