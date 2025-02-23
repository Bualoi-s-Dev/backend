package controllers

import (
	"net/http"
	
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusyTimeController struct {
	Service *services.BusyTimeService
}

func NewBusyTimeController(service *services.BusyTimeService) *BusyTimeController {
	return &BusyTimeController{Service: service}
}

// GetBusyTimesByPhotographerId godoc
// @Summary Get busy times by photographer ID
// @Description Get busy times by photographer ID
// @Tags BusyTime
// @Param photographerId path string true "Photographer ID"
// @Success 200 {object} []models.BusyTime
// @Failure 400 {object} string "Bad Request"
// @Router /busytime/photographer/{photographerId} [GET]
func (ctrl *BusyTimeController) GetBusyTimesByPhotographerId(c *gin.Context) {
	photographerIdParam := c.Param("photographerId")
	photographerId, err := primitive.ObjectIDFromHex(photographerIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photographer ID, " + err.Error()})
		return
	}
	items, err := ctrl.Service.GetByPhotographerId(c.Request.Context(), photographerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}