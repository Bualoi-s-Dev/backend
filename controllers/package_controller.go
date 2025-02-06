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

// GetAllPackages godoc
// @Tags Package
// @Summary Get a list of packages
// @Description Retrieve all packages from the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package [get]
// @x-order 1
func (ctrl *PackageController) GetAllPackages(c *gin.Context) {
	items, err := ctrl.Service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetOnePackage godoc
// @Tags Package
// @Summary Get a packages by id
// @Description Retrieve a packages which matched id from the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [get]
// @x-order 2
func (ctrl *PackageController) GetOnePackage(c *gin.Context) {
	id := c.Param("id")
	item, err := ctrl.Service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateOnePackage godoc
// @Tags Package
// @Summary Create a package
// @Description Create a package in the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package [post]
// @x-order 3
func (ctrl *PackageController) CreateOnePackage(c *gin.Context) {
	var item models.Package
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}
	if err := ctrl.Service.CreateOne(c.Request.Context(), &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item, " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateOnePackage godoc
// @Tags Package
// @Summary Patch a package
// @Description Update a package in some field to the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [patch]
// @x-order 4
func (ctrl *PackageController) UpdateOnePackage(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	if err := ctrl.Service.UpdateOne(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// DeleteOnePackage godoc
// @Tags Package
// @Summary Delete a package
// @Description Delete a package in the database
// @Success 200 {array} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [delete]
// @x-order 5
func (ctrl *PackageController) DeleteOnePackage(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.Service.DeleteOne(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
