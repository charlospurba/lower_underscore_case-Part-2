package handlers

import (
	"net/http"
	"strconv"

	"gin-user-app/dto"
	"gin-user-app/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler menangani request terkait user
type UserHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) SetValidator(validate *validator.Validate) {
	h.validate = validate
}

// GetUsers mendapatkan semua user
// @Summary Get all users
// @Description Retrieve a list of users
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} dto.UserDTO
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID mendapatkan user berdasarkan ID
// @Summary Get a user by ID
// @Description Retrieve user details by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser membuat user baru
// @Summary Create a new user
// @Description Register a new user with a Gmail email address
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param user body dto.CreateUserDTO true "User data (email must be @gmail.com)"
// @Success 201 {object} dto.UserDTO
// @Failure 400 {object} map[string][]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var userReq dto.CreateUserDTO

	// Bind JSON dan tangani jika gagal
	if err := c.ShouldBindJSON(&userReq); err != nil {
		// Tangani error validasi
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var errors []string
			for _, e := range validationErrs {
				switch e.Tag() {
				case "gmail":
					errors = append(errors, "email must be a valid Gmail address (@gmail.com)")
				case "email":
					errors = append(errors, "invalid email format")
				case "required":
					errors = append(errors, e.Field()+" is required")
				case "alphanum":
					errors = append(errors, e.Field()+" must contain only letters and numbers")
				case "min":
					errors = append(errors, e.Field()+" must be at least "+e.Param()+" characters")
				case "max":
					errors = append(errors, e.Field()+" must not exceed "+e.Param()+" characters")
				case "gt":
					errors = append(errors, e.Field()+" must be greater than "+e.Param())
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Panggil service untuk buat user
	user, err := h.userService.CreateUser(userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser memperbarui data user berdasarkan ID
// @Summary Update a user
// @Description Update user details with a Gmail email address
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserDTO true "Updated user data (email must be @gmail.com)"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} map[string][]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Cek apakah user ada sebelum diupdate
	_, err = h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userReq dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&userReq); err != nil {
		// Tangani error validasi
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var errors []string
			for _, e := range validationErrs {
				switch e.Tag() {
				case "gmail":
					errors = append(errors, "email must be a valid Gmail address (@gmail.com)")
				case "email":
					errors = append(errors, "invalid email format")
				case "required":
					errors = append(errors, e.Field()+" is required")
				case "alphanum":
					errors = append(errors, e.Field()+" must contain only letters and numbers")
				case "min":
					errors = append(errors, e.Field()+" must be at least "+e.Param()+" characters")
				case "max":
					errors = append(errors, e.Field()+" must not exceed "+e.Param()+" characters")
				case "gt":
					errors = append(errors, e.Field()+" must be greater than "+e.Param())
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.userService.UpdateUser(id, userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser menghapus user berdasarkan ID
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Cek apakah user ada sebelum dihapus
	_, err = h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
