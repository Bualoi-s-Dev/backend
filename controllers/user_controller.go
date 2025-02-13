package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	Service   *services.UserService
	S3Service *services.S3Service
}

func NewUserController(service *services.UserService, s3Service *services.S3Service) *UserController {
	return &UserController{Service: service, S3Service: s3Service}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data, " + err.Error()})
		return
	}

	// Return the profile picture URL
	c.JSON(http.StatusOK, userDb)
}

// GetUserProfileByID godoc
// @Tags User
// @Summary Get a user profile by ID
// @Description Retrieve a user profile from database by user ID
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile/{id} [get]
func (uc *UserController) GetUserProfileByID(c *gin.Context) {
	userID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userDb, err := uc.Service.GetUserByID(c.Request.Context(), objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

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

	var userBody dto.UserRequest
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userBody.Profile != nil {
		if err := uc.S3Service.VerifyBase64(*userBody.Profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile picture, " + err.Error()})
			return
		}
	}

	if userBody.ShowcasePackages != nil {
		isShowcaseValid := uc.Service.VerifyShowcase(c.Request.Context(), user.Packages, *userBody.ShowcasePackages)
		if !isShowcaseValid {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
			return
		}
	}

	// Call the service to update the user's profile, include picture
	newUser, err := uc.Service.UpdateUser(c.Request.Context(), user.ID.Hex(), user.Email, &userBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

// UpdateUserShowcasePackage godoc
// @Tags User
// @Summary Update user showcase packages
// @Description Receive showcase packageIds and put it in user data
// @Param request body UserShowcasePackageRequest true "Update showcase packages Request"
// @Success 200 {object} models.User
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile/showcase [put]
// func (uc *UserController) UpdateUserShowcasePackage(c *gin.Context) {
// 	var req dto.UserShowcasePackageRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}

// 	// Check the owner
// 	user := middleware.GetUserFromContext(c)
// 	isShowcaseValid := uc.Service.VerifyShowcase(c.Request.Context(), user.ShowcasePackages, req.PackageID)
// 	if !isShowcaseValid {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
// 		return
// 	}

// 	user.ShowcasePackages = req.PackageID
// 	_, err := uc.Service.UpdateUser(c.Request.Context(), user.Email, user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update showcase packages, " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// 	return
// }
