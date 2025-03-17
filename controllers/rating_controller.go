package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/apperrors"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingController struct {
	Service *services.RatingService
}

func NewRatingController(service *services.RatingService) *RatingController {
	return &RatingController{Service: service}
}

// GetAllRatingsFromPhotographer godoc
// @Summary Get all ratings from a photographer
// @Description Retrieve all ratings received by a specific photographer
// @Tags Rating
// @Param photographerId path string true "Photographer ID"
// @Success 200 {object} []models.Rating
// @Failure 400 {object} string "Bad Request"
// @Router /user/{photographerId}/rating [GET]
func (ctrl *RatingController) GetAllRatingsFromPhotographer(c *gin.Context) {
	photographerId := c.Param("photographerId")

	// Convert photographerId to ObjectID
	photographerObjectID, err := primitive.ObjectIDFromHex(photographerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photographer ID"})
		return
	}

	// Fetch ratings for the photographer
	items, err := ctrl.Service.GetByPhotographerId(c.Request.Context(), photographerObjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ratings, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateRating godoc
// @Summary Create a rating for a photographer
// @Description Customers can create a rating for a photographer
// @Tags Rating
// @Param photographerId path string true "Photographer ID"
// @Param request body dto.RatingRequest true "Create Rating Request"
// @Success 201 {object} models.Rating
// @Failure 403 {object} string "Forbidden"
// @Failure 400 {object} string "Bad Request"
// @Router /user/{photographerId}/rating [POST]
func (ctrl *RatingController) CreateRating(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	photographerId := c.Param("photographerId")
	photographerObjectID, err := primitive.ObjectIDFromHex(photographerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photographer ID"})
		return
	}

	var ratingRequest dto.RatingRequest
	if err := c.ShouldBindJSON(&ratingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, " + err.Error()})
		return
	}

	item, err := ctrl.Service.CreateOneFromCustomer(c.Request.Context(), &ratingRequest, user.ID, photographerObjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating, " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetRatingById godoc
// @Summary Get a specific rating by ID
// @Description Retrieve a rating by ID for a specific photographer
// @Tags Rating
// @Param photographerId path string true "Photographer ID"
// @Param ratingId path string true "Rating ID"
// @Success 200 {object} models.Rating
// @Failure 400 {object} string "Bad Request"
// @Router /user/{photographerId}/rating/{ratingId} [GET]
func (ctrl *RatingController) GetRatingById(c *gin.Context) {
	ratingId := c.Param("ratingId")
	ratingObjectID, err := primitive.ObjectIDFromHex(ratingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	item, err := ctrl.Service.GetById(c.Request.Context(), ratingObjectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateRating godoc
// @Summary Update an existing rating
// @Description Update an existing rating by ID
// @Tags Rating
// @Param photographerId path string true "Photographer ID"
// @Param ratingId path string true "Rating ID"
// @Param request body dto.RatingRequest true "Update Rating Request"
// @Success 200 {object} models.Rating
// @Failure 403 {object} string "Forbidden"
// @Failure 400 {object} string "Bad Request"
// @Router /user/{photographerId}/rating/{ratingId} [PATCH]
func (ctrl *RatingController) UpdateRating(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	ratingId := c.Param("ratingId")
	ratingObjectID, err := primitive.ObjectIDFromHex(ratingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// Bind request JSON to DTO
	var ratingRequest dto.RatingRequest
	if err := c.ShouldBindJSON(&ratingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	
	err = ctrl.Service.UpdateOne(c.Request.Context(), user.ID, ratingObjectID, &ratingRequest)
	if err != nil {
		if err == apperrors.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own review"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rating, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating updated successfully"})
}

// DeleteRating godoc
// @Summary Delete an existing rating
// @Description Delete a rating by ID
// @Tags Rating
// @Param photographerId path string true "Photographer ID"
// @Param ratingId path string true "Rating ID"
// @Success 200 {object} string "Rating id {ratingId} deleted successfully"
// @Failure 403 {object} string "Forbidden"
// @Failure 400 {object} string "Bad Request"
// @Router /user/{photographerId}/rating/{ratingId} [DELETE]
func (ctrl *RatingController) DeleteRating(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	ratingId := c.Param("ratingId")
	ratingObjectID, err := primitive.ObjectIDFromHex(ratingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// Call service to delete the rating
	err = ctrl.Service.DeleteOne(c.Request.Context(), user.ID, ratingObjectID)
	if err != nil {
		if err == apperrors.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own review"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}
