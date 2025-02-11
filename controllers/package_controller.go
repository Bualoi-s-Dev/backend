package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type PackageController struct {
	Service     *services.PackageService
	S3Service   *services.S3Service
	UserService *services.UserService
}

func NewPackageController(service *services.PackageService, s3Service *services.S3Service, userService *services.UserService) *PackageController {
	return &PackageController{Service: service, S3Service: s3Service, UserService: userService}
}

// GetAllPackages godoc
// @Tags Package
// @Summary Get a list of packages
// @Description Retrieve all packages from the database
// @Success 200 {object} []models.Package
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
// @Success 200 {object} models.Package
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
// @Param request body dto.PackageStrictRequest true "Create Package Request"
// @Success 200 {object} models.Package
// @Failure 400 {object} string "Bad Request"
// @Router /package [post]
// @x-order 3
func (ctrl *PackageController) CreateOnePackage(c *gin.Context) {
	var itemInput dto.PackageStrictRequest
	if err := c.ShouldBindJSON(&itemInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	item, err := ctrl.Service.CreateOne(c.Request.Context(), &itemInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item, " + err.Error()})
		return
	}

	user := middleware.GetUserFromContext(c)
	user.Packages = append(user.Packages, item.ID)
	_, err = ctrl.UserService.UpdateUser(c.Request.Context(), user.Email, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user, " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateOnePackage godoc
// @Tags Package
// @Summary Patch a package
// @Param request body dto.PackageRequest true "Replace Package Request"
// @Success 200 {object} models.Package
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [put]
// @x-order 4
func (ctrl *PackageController) ReplaceOnePackage(c *gin.Context) {
	id := c.Param("id")
	// Check the owner, only the owner can update the package
	user := middleware.GetUserFromContext(c)
	isOwner := CheckOwner(user, id)
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	var updates models.PackageRequest
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	item, err := ctrl.Service.ReplaceOne(c.Request.Context(), id, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteOnePackage godoc
// @Tags Package
// @Summary Delete a package
// @Description Delete a package in the database
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Router /package/{id} [delete]
// @x-order 5
func (ctrl *PackageController) DeleteOnePackage(c *gin.Context) {
	id := c.Param("id")
	// Check the owner, only the owner can delete the package
	user := middleware.GetUserFromContext(c)
	isOwner := CheckOwner(user, id)
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	if err := ctrl.Service.DeleteOne(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item, " + err.Error()})
		return
	}

	// Remove the package from the user's packages
	for i, packageId := range user.Packages {
		if packageId.Hex() == id {
			user.Packages = append(user.Packages[:i], user.Packages[i+1:]...)
			break
		}
	}
	_, err := ctrl.UserService.UpdateUser(c.Request.Context(), user.Email, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// Helper function

func CheckOwner(user *models.User, packageId string) bool {
	for _, id := range user.Packages {
		if id.Hex() == packageId {
			return true
		}
	}
	return false
}
