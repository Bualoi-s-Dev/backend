package controllers

import (
	"fmt"
	"net/http"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubpackageController struct {
	Service        *services.SubpackageService
	PackageService *services.PackageService
}

func NewSubpackageController(service *services.SubpackageService, packageService *services.PackageService) *SubpackageController {
	return &SubpackageController{Service: service, PackageService: packageService}
}

func (ctrl *SubpackageController) GetAllSubpackages(c *gin.Context) {
	items, err := ctrl.Service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ctrl *SubpackageController) CreateSubpackage(c *gin.Context) {
	var itemRequest dto.SubpackageRequest
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, " + err.Error()})
		return
	}
	if err := ctrl.Service.VerifyStrictRequest(c.Request.Context(), &itemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request model, " + err.Error()})
		return
	}

	packageIdParam := c.Param("packageId")
	packageId, err := primitive.ObjectIDFromHex(packageIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID, " + err.Error()})
		return
	}
	if _, err := ctrl.PackageService.GetById(c.Request.Context(), packageIdParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID, " + err.Error()})
		return
	}
	item := itemRequest.ToModel()
	item.PackageID = packageId
	if err := ctrl.Service.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *SubpackageController) UpdateSubpackage(c *gin.Context) {
	id := c.Param("subpackageId")
	var itemRequest dto.SubpackageRequest
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, " + err.Error()})
		return
	}

	if err := ctrl.Service.Update(c.Request.Context(), id, &itemRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item, " + err.Error()})
		return
	}

	item, err := ctrl.Service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item after updated, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ctrl *SubpackageController) DeleteSubpackage(c *gin.Context) {
	id := c.Param("subpackageId")
	if err := ctrl.Service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Subpackage id %s deleted successfully", id)})
}
