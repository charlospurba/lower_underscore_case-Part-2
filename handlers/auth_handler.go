package handlers

import (
	"gin-user-app/dto"
	"gin-user-app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
// @Security BearerAuth
// @Param loginRequest body dto.LoginRequestDTO true "Login Request"
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
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
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Failed to logout"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Ambil token dari header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	// Ambil token asli (tanpa "Bearer ")
	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}
	tokenString = tokenParts[1]

	// Simpan token ke dalam blacklist
	err := h.authService.Logout(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// VerifyToken godoc
// @Summary Verify JWT token
// @Description Verify the provided JWT token and return user information
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} dto.UserDTO "User data"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/verify [get]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	token := authHeader[len("Bearer "):]

	user, err := h.authService.VerifyUser(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	})
}
