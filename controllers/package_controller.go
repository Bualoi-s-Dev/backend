package controllers

import (
	"fmt"
	"net/http"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
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
// @Param id path string true "Package ID"
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
	if err := ctrl.S3Service.VerifyMultipleBase64(itemInput.Photos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request Image, " + err.Error()})
	}

	user := middleware.GetUserFromContext(c)
	item, err := ctrl.Service.CreateOne(c.Request.Context(), &itemInput, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item, " + err.Error()})
		return
	}

	userReq := dto.UpdateUserPackageRequest{Packages: append(user.Packages, item.ID)}
	err = ctrl.UserService.UpdateOwnerPackage(c.Request.Context(), user.ID, userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user, " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateOnePackage godoc
// @Tags Package
// @Summary Patch a package
// @Param id path string true "Package ID"
// @Param request body dto.PackageRequest true "Replace Package Request"
// @Success 200 {object} models.Package
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
	isOwner := ctrl.Service.CheckOwner(user, id)
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
	c.JSON(http.StatusOK, item)
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
	isOwner := ctrl.Service.CheckOwner(user, id)
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own the package"})
		return
	}

	if err := ctrl.Service.DeleteOne(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item, " + err.Error()})
		return
	}

	userReq := dto.UpdateUserPackageRequest{}
	// Remove the package from the user's packages
	fmt.Println("id :", id)
	fmt.Println("user.Packages :", user.Packages)
	for i, packageId := range user.Packages {
		fmt.Println("packageId.Hex() :", packageId.Hex())
		if packageId.Hex() == id {
			userReq.Packages = append(user.Packages[:i], user.Packages[i+1:]...)
			break
		}
	}
	// Remove the package from the user's showcase packages
	fmt.Println("user.ShowcasePackages :", user.ShowcasePackages)
	for i, packageId := range user.ShowcasePackages {
		fmt.Println("packageId.Hex() :", packageId.Hex())
		if packageId.Hex() == id {
			userReq.ShowcasePackages = append(user.ShowcasePackages[:i], user.ShowcasePackages[i+1:]...)
			break
		}
	}
	err = ctrl.UserService.UpdateOwnerPackage(c.Request.Context(), user.ID, userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user, " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
