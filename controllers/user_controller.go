package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Description Retrieve a user from firebase jwt
// @Success 200 {object} models.User
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

// GetUserProfile godoc
// @Tags User
// @Summary Get a user from database
// @Description Retrieve a user from database
// @Success 200 {object} models.User
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile [get]
func (uc *UserController) GetUserProfile(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user := middleware.GetUserFromContext(c)

	// Call the service to get the user's profile picture URL
	userDb, err := uc.Service.GetUser(c.Request.Context(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	// Return the profile picture URL
	c.JSON(http.StatusOK, userDb)
}

// UpdateUserProfile godoc
// @Tags User
// @Summary Update user data
// @Description Receive a user data and update it, the profile is base64 and return in url
// @Param request body models.User true "Update User Request"
// @Success 200 {object} models.User
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile [put]
func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user := middleware.GetUserFromContext(c)

	var userBody models.User
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userBody.ID = user.ID

	isShowcaseValid := uc.Service.VerifyShowcase(c.Request.Context(), user.ShowcasePackages, userBody.ShowcasePackages)
	if !isShowcaseValid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	// Call the service to update the user's profile, include picture
	newUser, err := uc.Service.UpdateUserWithNewImage(c.Request.Context(), user.ID.Hex(), user.Email, &userBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

type UserShowcasePackageRequest struct {
	PackageID []primitive.ObjectID `json:"packageIds" binding:"required" example:"12345678abcd,12345678abcd"`
}

// UpdateUserShowcasePackage godoc
// @Tags User
// @Summary Update user showcase packages
// @Description Receive showcase packageIds and put it in user data
// @Param request body UserShowcasePackageRequest true "Update showcase packages Request"
// @Success 200 {object} models.User
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile/showcase [put]
func (uc *UserController) UpdateUserShowcasePackage(c *gin.Context) {
	var req UserShowcasePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Check the owner
	user := middleware.GetUserFromContext(c)
	isShowcaseValid := uc.Service.VerifyShowcase(c.Request.Context(), user.ShowcasePackages, req.PackageID)
	if !isShowcaseValid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	user.ShowcasePackages = req.PackageID
	_, err := uc.Service.UpdateUser(c.Request.Context(), user.Email, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update showcase packages, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
	return
}
