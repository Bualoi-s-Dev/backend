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

// GetUserProfile godoc
// @Tags User
// @Summary Get a user from jwt
// @Description Retrieve a user which matched email from the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /user/me [get]
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