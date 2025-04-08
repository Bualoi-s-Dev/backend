package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
// @Param packageId query string false "Package ID of subpackage"
// @Param type query string false "Type of subpackage"
// @Param availableStartTime query string false "Available start time of subpackage"
// @Param availableEndTime query string false "Available end time of subpackage"
// @Param availableStartDay query string false "Available start day of subpackage"
// @Param availableEndDay query string false "Available end day of subpackage"
// @Param repeatedDay query []string false "Repeated day of subpackage"
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Limit number of items per page, default is 10"
// @Tags Subpackage
// @Success 200 {object} []dto.FilteredSubpackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /subpackage [GET]
func (ctrl *SubpackageController) GetAllSubpackages(c *gin.Context) {
	// Get query parameters for filtering
	filters := map[string]string{
		"packageId":          c.Query("packageId"),
		"type":               c.Query("type"),
		"availableStartTime": c.Query("availableStartTime"),
		"availableEndTime":   c.Query("availableEndTime"),
		"availableStartDay":  c.Query("availableStartDay"),
		"availableEndDay":    c.Query("availableEndDay"),
		"repeatedDay":        c.Query("repeatedDay"),
	}

	// Verify date format
	dateFormat := "2006-01-02"
	if filters["availableStartDay"] != "" {
		if _, err := time.Parse(dateFormat, filters["availableStartDay"]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format for availableStartDay. Use YYYY-MM-DD."})
			return
		}
	}
	if filters["availableEndDay"] != "" {
		if _, err := time.Parse(dateFormat, filters["availableEndDay"]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format for availableEndDay. Use YYYY-MM-DD."})
			return
		}
	}

	// Verify time format (HH:MM)
	timeFormat := "15:03"
	if filters["availableStartTime"] != "" {
		if _, err := time.Parse(timeFormat, filters["availableStartTime"]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format for availableStartTime. Use HH:MM."})
			return
		}
	}
	if filters["availableEndTime"] != "" {
		if _, err := time.Parse(timeFormat, filters["availableEndTime"]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format for availableEndTime. Use HH:MM."})
			return
		}
	}

	// Verify repeatedDay format
	validDays := map[string]bool{
		"SUN": true, "MON": true, "TUE": true, "WED": true, "THU": true, "FRI": true, "SAT": true,
	}
	if filters["repeatedDay"] != "" {
		days := strings.Split(filters["repeatedDay"], ",")
		for _, day := range days {
			if !validDays[strings.TrimSpace(day)] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repeatedDay value. Use a comma-separated list of SUN, MON, TUE, WED, THU, FRI, SAT."})
				return
			}
		}
	}

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Retrieve and filter subpackages with pagination
	subpackages, totalCount, err := ctrl.Service.GetFilteredSubpackages(c.Request.Context(), filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}

	if len(subpackages) == 0 {
		subpackages = []dto.SubpackageResponse{}
	}

	// Calculate max pages
	maxPage := (totalCount + limit - 1) / limit // Equivalent to ceil(totalCount / limit)

	response := dto.FilteredSubpackageResponse{
		Subpackages: subpackages,
		Pagination: dto.Pagination{
			Page:    page,
			Limit:   limit,
			MaxPage: maxPage,
			Total:   totalCount,
		},
	}

	c.JSON(http.StatusOK, response)
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
	// Min price
	if *itemRequest.Price < 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 20"})
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
