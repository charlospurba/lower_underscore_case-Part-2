package handlers

import (
	"gin-user-app/dto"
	"gin-user-app/models"
	"gin-user-app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UserHandler handles user-related requests.
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.UserDTO "List of users"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var usersDTO []dto.UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, dto.UserDTO{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, usersDTO)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieve user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDTO "User details"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userDTO := dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, userDTO)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param userRequest body dto.CreateUserDTO true "User Request"
// @Success 201 {object} models.User "User created"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserReq dto.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:  createUserReq.Username,
		Password:  createUserReq.Password, // Must hash password before saving
		Email:     createUserReq.Email,
		FirstName: createUserReq.FirstName,
		LastName:  createUserReq.LastName,
		Age:       createUserReq.Age,
	}

	err := h.userService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param userRequest body dto.UpdateUserDTO true "User Update Request"
// @Success 200 {object} models.User "User updated"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updateUserReq dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:        id,
		Username:  updateUserReq.Username,
		Password:  updateUserReq.Password, // Must hash password before saving
		Email:     updateUserReq.Email,
		FirstName: updateUserReq.FirstName,
		LastName:  updateUserReq.LastName,
		Age:       updateUserReq.Age,
	}

	err := h.userService.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
