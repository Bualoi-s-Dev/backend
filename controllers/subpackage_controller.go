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

// GetAllSubpackages godoc
// @Summary Get all subpackages
// @Description Get all subpackages
// @Tags Subpackage
// @Success 200 {object} []dto.SubpackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /subpackage [GET]
func (ctrl *SubpackageController) GetAllSubpackages(c *gin.Context) {
	items, err := ctrl.Service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}

	var responses []dto.SubpackageResponse
	for _, item := range items {
		res, err := ctrl.Service.MappedToSubpackageResponse(c, &item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map items, " + err.Error()})
			return
		}
		responses = append(responses, *res)
	}
	c.JSON(http.StatusOK, responses)
}

// GetByIdSubpackages godoc
// @Summary Get subpackages by ID
// @Description Get subpackages by ID
// @Param id path string true "Subpackage ID"
// @Tags Subpackage
// @Success 200 {object} dto.SubpackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /subpackage/{id} [GET]
func (ctrl *SubpackageController) GetByIdSubpackages(c *gin.Context) {
	id := c.Param("id")
	item, err := ctrl.Service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}

	response, err := ctrl.Service.MappedToSubpackageResponse(c, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map items, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// CreateSubpackage godoc
// @Summary Create a subpackage for a package
// @Description Create a subpackage for a package, require all fields in the request
// @Tags Subpackage
// @Param packageId path string true "Package ID"
// @Param request body dto.SubpackageRequest true "Create Subpackage Request"
// @Success 200 {object} dto.SubpackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /subpackage/{packageId} [POST]
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

	response, err := ctrl.Service.MappedToSubpackageResponse(c, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map items, " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

// UpdateSubpackage godoc
// @Summary Update an existed subpackage
// @Description Update an existed subpackage
// @Tags Subpackage
// @Param id path string true "Subpackage ID"
// @Param request body dto.SubpackageRequest true "Update Subpackage Request"
// @Success 200 {object} dto.SubpackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /subpackage/{id} [PATCH]
func (ctrl *SubpackageController) UpdateSubpackage(c *gin.Context) {
	id := c.Param("id")
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

	response, err := ctrl.Service.MappedToSubpackageResponse(c, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map items, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// DeleteSubpackage godoc
// @Summary Delete an existed subpackage
// @Description Delete an existed subpackage
// @Tags Subpackage
// @Param id path string true "Subpackage ID"
// @Success 200 {object} string "Subpackage id {id} deleted successfully"
// @Failure 400 {object} string "Bad Request"
// @Router /subpackage/{id} [DELETE]
func (ctrl *SubpackageController) DeleteSubpackage(c *gin.Context) {
	id := c.Param("id")
	if _, err := ctrl.Service.GetById(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch item, " + err.Error()})
		return
	}

	if err := ctrl.Service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Subpackage id %s deleted successfully", id)})
}
