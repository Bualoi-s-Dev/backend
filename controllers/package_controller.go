package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageController struct {
	Service           *services.PackageService
	S3Service         *services.S3Service
	UserService       *services.UserService
	SubpackageService *services.SubpackageService
}

func NewPackageController(service *services.PackageService, s3Service *services.S3Service, userService *services.UserService, subpackageService *services.SubpackageService) *PackageController {
	return &PackageController{Service: service, S3Service: s3Service, UserService: userService, SubpackageService: subpackageService}
}

// GetAllPackages godoc
// @Tags Package
// @Summary Get a list of packages
// @Description Retrieve all packages from the database
// @Param search query string false "Filter by package title or owner (prefix filter)"
// @Param type query string false "Filter by package type(prefix filter)"
// @Param minPrice query int false "minimum price"
// @Param maxPrice query int false "maximum price"
// @Success 200 {object} []dto.PackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /package [get]
// @x-order 1
func (ctrl *PackageController) GetAllPackages(c *gin.Context) {
	items, err := ctrl.Service.GetAll(c.Request.Context())
	searchString, _ := c.GetQuery("search")
	searchType_, _ := c.GetQuery("type")
	searchType := models.PackageType(searchType_)

	minPrice_, hasMinPrice := c.GetQuery("minPrice")
	maxPrice_, hasMaxPrice := c.GetQuery("maxPrice")
	var minPrice, maxPrice int

	if hasMinPrice {
		minPrice, err = strconv.Atoi(minPrice_)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
			return
		}
	} else {
		minPrice = 0
	}

	if hasMaxPrice {
		maxPrice, err = strconv.Atoi(maxPrice_)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
			return
		}
	} else {
		maxPrice = math.MaxInt32
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}

	res := []dto.PackageResponse{}
	for _, item := range items {
		IsFiltered, err := ctrl.Service.FilterPackage(c.Request.Context(), &item, searchString, searchType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter item, " + err.Error()})
			return
		}
		if !IsFiltered {
			continue
		}

		mappedItem, err := ctrl.Service.MappedToPackageResponse(c.Request.Context(), &item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map item, " + err.Error()})
			return
		}

		if hasMaxPrice || hasMinPrice {
			IsFiltered, err = ctrl.Service.FilterPrice(c.Request.Context(), mappedItem, minPrice, maxPrice)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter price, " + err.Error()})
				return
			}
			if !IsFiltered {
				continue
			}
		}

		res = append(res, *mappedItem)
	}
	c.JSON(http.StatusOK, res)
}

// GetRecommendedPackages godoc
// @Tags Package
// @Summary Get a list of recommended packages
// @Description Retrieve recommended packages from the database by size
// @Param size query int false "Recommended package size"
// @Success 200 {object} []dto.PackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /package/recommend [get]
// @x-order 6
func (ctrl *PackageController) GetRecommendedPackages(c *gin.Context) {
	size_ := c.Query("size")
	size, err := strconv.Atoi(size_)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	items, err := ctrl.Service.GetAllRecommended(c.Request.Context(), size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}

	res := []dto.PackageResponse{}
	for _, item := range items {
		mappedItem, err := ctrl.Service.MappedToPackageResponse(c.Request.Context(), &item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map item, " + err.Error()})
			return
		}
		res = append(res, *mappedItem)
	}

	c.JSON(http.StatusOK, res)
}

// GetOnePackage godoc
// @Tags Package
// @Summary Get a packages by id
// @Description Retrieve a packages which matched id from the database
// @Param id path string true "Package ID"
// @Success 200 {object} dto.PackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [get]
// @x-order 2
func (ctrl *PackageController) GetOnePackage(c *gin.Context) {
	id := c.Param("id")
	item, err := ctrl.Service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch item, " + err.Error()})
		return
	}

	res, err := ctrl.Service.MappedToPackageResponse(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// CreateOnePackage godoc
// @Tags Package
// @Summary Create a package
// @Description Create a package in the database
// @Param request body dto.PackageRequest true "Create Package Request"
// @Success 200 {object} dto.PackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /package [post]
// @x-order 3
func (ctrl *PackageController) CreateOnePackage(c *gin.Context) {
	var itemInput dto.PackageRequest
	if err := c.ShouldBindJSON(&itemInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}
	if err := ctrl.Service.VerifyStrictRequest(c.Request.Context(), &itemInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}
	if err := ctrl.S3Service.VerifyMultipleBase64(*itemInput.Photos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request Image, " + err.Error()})
	}

	user := middleware.GetUserFromContext(c)
	item, err := ctrl.Service.CreateOne(c.Request.Context(), &itemInput, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item, " + err.Error()})
		return
	}

	res, err := ctrl.Service.MappedToPackageResponse(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map item, " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// UpdateOnePackage godoc
// @Tags Package
// @Summary Patch a package
// @Param id path string true "Package ID"
// @Param request body dto.PackageRequest true "Replace Package Request"
// @Success 200 {object} dto.PackageResponse
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [patch]
// @x-order 4
func (ctrl *PackageController) UpdateOnePackage(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.Service.CheckPackageExist(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID, " + err.Error()})
		return
	}
	// Check the owner, only the owner can update the package
	user := middleware.GetUserFromContext(c)
	isOwner, err := ctrl.Service.CheckOwner(c.Request.Context(), user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check owner, " + err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	var updates dto.PackageRequest
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}
	if updates.Photos != nil && len(*updates.Photos) > 0 {
		if err := ctrl.S3Service.VerifyMultipleBase64(*updates.Photos); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request Image, " + err.Error()})
		}
	}

	item, err := ctrl.Service.UpdateOne(c.Request.Context(), id, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item, " + err.Error()})
		return
	}

	res, err := ctrl.Service.MappedToPackageResponse(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteOnePackage godoc
// @Tags Package
// @Summary Delete a package
// @Description Delete a package in the database
// @Param id path string true "Package ID"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [delete]
// @x-order 5
func (ctrl *PackageController) DeleteOnePackage(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.Service.CheckPackageExist(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID, " + err.Error()})
		return
	}
	// Check the owner, only the owner can delete the package
	user := middleware.GetUserFromContext(c)
	isOwner, err := ctrl.Service.CheckOwner(c.Request.Context(), user, id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to check owner, " + err.Error()})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	isPackageDeletable := true
	var errSubpackage models.Subpackage
	hexId, _ := primitive.ObjectIDFromHex(id)
	subpackages, err := ctrl.SubpackageService.GetByPackageId(c.Request.Context(), hexId)
	fmt.Println("Subpackages: ", subpackages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subpackages, " + err.Error()})
		return
	}
	for _, subpackage := range subpackages {
		isSubpackageDeletable, err := ctrl.SubpackageService.IsSubpackageDeletable(c.Request.Context(), subpackage.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check subpackage deletable, " + err.Error()})
			return
		}
		if !isSubpackageDeletable {
			isPackageDeletable = false
			errSubpackage = subpackage
			break
		}
	}
	if !isPackageDeletable {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Package cannot be deleted because subpackage %s is not deletable", errSubpackage.Title)})
		return
	}

	if err := ctrl.Service.DeleteOne(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item, " + err.Error()})
		return
	}
	for _, subpackage := range subpackages {
		err := ctrl.SubpackageService.Delete(c.Request.Context(), subpackage.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subpackage, " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
