package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

// GetUserJWT godoc
// @Tags User
// @Summary Get a user from jwt
// @Description Retrieve a user which matched email from the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /user/me [get]
func (uc *UserController) GetUserJWT(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserProfile(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert user to be struct
	userData, ok := user.(*models.User)
	if !ok || userData.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// Call the service to get the user's profile picture URL
	userDb, err := uc.Service.GetUser(c.Request.Context(), userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	// Return the profile picture URL
	c.JSON(http.StatusOK, userDb)
}

func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert user to be struct
	userData, ok := user.(*models.User)
	if !ok || userData.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	var userBody models.User
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the user's profile picture
	newUser, err := uc.Service.UpdateUser(c.Request.Context(), userData.ID.Hex(), userData.Email, &userBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}
