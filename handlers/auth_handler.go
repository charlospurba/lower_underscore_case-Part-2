package handlers

import (
	"gin-user-app/dto"
	"gin-user-app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler handles authentication-related requests.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

// Login godoc
// @Summary Login to the system
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequestDTO true "Login Request"
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout godoc
// @Summary Logout from the system
// @Description Logout and invalidate the current session
// @Tags auth
// @Success 200 {object} map[string]string "Success"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// VerifyToken godoc
// @Summary Verify JWT token
// @Description Verify the provided JWT token and return user information
// @Tags auth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} dto.UserDTO "User data"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/verify [get]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	// Ambil token dari header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Ambil token dengan format Bearer
	token := authHeader[len("Bearer "):]

	user, err := h.authService.VerifyUser(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Menggunakan DTO UserDTO untuk response
	c.JSON(http.StatusOK, gin.H{
		"user": dto.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	})
}
