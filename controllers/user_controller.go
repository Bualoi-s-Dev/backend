package controllers

import (
	"net/http"
	"strconv"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	Service         *services.UserService
	S3Service       *services.S3Service
	BusyTimeService *services.BusyTimeService
}

func NewUserController(service *services.UserService, s3Service *services.S3Service, busyTimeService *services.BusyTimeService) *UserController {
	return &UserController{Service: service, S3Service: s3Service, BusyTimeService: busyTimeService}
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
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile [get]
func (uc *UserController) GetUserProfile(c *gin.Context) {
	// Retrieve user from context (set by FirebaseAuthMiddleware)
	user := middleware.GetUserFromContext(c)

	// Call the service to get the user's profile picture URL
	userDb, err := uc.Service.GetUserByEmail(c.Request.Context(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data, " + err.Error()})
		return
	}

	// Return the profile picture URL
	c.JSON(http.StatusOK, userDb)
}

// GetUserPhogographer godoc
// @Tags User
// @Summary Get user photographers from database
// @Description Retrieve user photographers from database
// @Param name query string false "Photographer name"
// @Param location query string false "Photographer location"
// @Param minPrice query string false "Minimum price"
// @Param maxPrice query string false "Maximum price"
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Limit number of items per page, default is 10"
// @Param type query string false "Package type"
// @Success 200 {object} []dto.FilteredUserPhotographerResponse
// @Failure 400 {object} string "Bad Request"
// @Router /user/photographers [get]
func (uc UserController) GetPhotographers(c *gin.Context) {
	// Get query parameters for filtering
	filters := map[string]string{
		"name":     c.Query("name"),
		"location": c.Query("location"),
		"minPrice": c.Query("minPrice"),
		"maxPrice": c.Query("maxPrice"),
		"type":     c.Query("type"),
	}

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Call the service to get the user's profile picture URL
	photographers, totalCount, err := uc.Service.GetFilteredPhotographers(c.Request.Context(), filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photographers, " + err.Error()})
		return
	}

	//Return the profile picture URL
	// Calculate max pages
	maxPage := (totalCount + limit - 1) / limit // Equivalent to ceil(totalCount / limit)

	// Construct response
	response := dto.FilteredUserPhotographerResponse{
		Photographers: photographers,
		Pagination: dto.Pagination{
			Page:    page,
			Limit:   limit,
			MaxPage: maxPage,
			Total:   totalCount,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetUserProfileByID godoc
// @Tags User
// @Summary Get a user profile by ID
// @Description Retrieve a user profile from database by user ID
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
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
// @Param request body dto.UserRequest true "Update User Request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} string "Bad Request"
// @Router /user/profile [patch]
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
		isShowcaseValid, err := uc.Service.VerifyShowcase(c.Request.Context(), user.ID, *userBody.ShowcasePackages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify showcase, " + err.Error()})
			return
		}
		if !isShowcaseValid {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
			return
		}
	}

	// Call the service to update the user's profile, include picture
	newUser, err := uc.Service.UpdateUser(c.Request.Context(), user.ID, user.Email, &userBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

// CreateUserBusyTime godoc
// @Summary Create a busy time for the authenticated user
// @Description Create a busy time entry using user ID from the JWT
// @Tags User
// @Param request body dto.BusyTimeStrictRequest true "Create BusyTime Request"
// @Success 201 {object} dto.BusyTimeStrictRequest
// @Failure 400 {object} string "Bad Request"
// @Router /user/busytime [post]
func (uc *UserController) CreateUserBusyTime(c *gin.Context) {
	var busyTimeRequest dto.BusyTimeStrictRequest
	if err := c.ShouldBindJSON(&busyTimeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, " + err.Error()})
		return
	}

	// Extract user from context using existing middleware function
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user not found"})
		return
	}

	// Call BusyTimeService with extracted user ID
	res, err := uc.BusyTimeService.CreateFromUser(c.Request.Context(), &busyTimeRequest, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create busy time, " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// DeleteUserBusyTime godoc
// @Summary Delete a photographer busy time
// @Description Delele a busy time entry using busy time ID, require to be the owner of the busy time
// @Tags User
// @Param busyTimeId path string true "BusyTime ID"
// @Success 200 {object} string "Success"
// @Failure 400 {object} string "Bad Request"
// @Router /user/busytime/{busyTimeId} [delete]
func (uc *UserController) DeleteUserBusyTime(c *gin.Context) {
	id := c.Param("busyTimeId")
	busyTime, err := uc.BusyTimeService.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found busy time ID, " + err.Error()})
		return
	}
	if busyTime.Type == models.TypeAppointment {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete appointment type BusyTime"})
		return
	}

	user := middleware.GetUserFromContext(c)
	if user.ID != busyTime.PhotographerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the busy time"})
		return
	}

	if err := uc.BusyTimeService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete busy time, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Busy time deleted"})
}
