package routes

import (
	"gin-user-app/handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, authMiddleware gin.HandlerFunc) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.GET("/verify", authMiddleware, authHandler.VerifyToken) 
	}
}
