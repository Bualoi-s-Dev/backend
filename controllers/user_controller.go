package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

// GetUserProfile handles GET /user/me
func (uc *UserController) GetUserProfile(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// convert user to be struct
	userData, ok := user.(*models.User)
	if !ok || userData.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// updated requested data from request body using "ShouldBindJSON"
	var updates models.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// update data
	err := uc.Service.UpdateUser(c.Request.Context(), userData.Email, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}