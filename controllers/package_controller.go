package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type PackageController struct {
	Service *services.PackageService
}

func NewPackageController(service *services.PackageService) *PackageController {
	return &PackageController{Service: service}
}

func (ctrl *PackageController) GetAllPackages(c *gin.Context) {
	items, err := ctrl.Service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ctrl *PackageController) CreateOnePackage(c *gin.Context) {
	var item models.Package
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.Service.CreateOne(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}
	c.JSON(http.StatusCreated, item)
}
