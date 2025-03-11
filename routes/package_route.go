package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/gin-gonic/gin"
)

func PackageRoutes(router *gin.Engine, ctrl *controllers.PackageController) {
	packageRoutes := router.Group("/package")
	commonRoutes := packageRoutes.Group("", middleware.AllowRoles(models.Photographer, models.Customer))
	{
		commonRoutes.GET("", ctrl.GetAllPackages)
		commonRoutes.GET("/:id", ctrl.GetOnePackage)
	}

	photographerRoutes := packageRoutes.Group("", middleware.AllowRoles(models.Photographer))
	{

		photographerRoutes.POST("", ctrl.CreateOnePackage)
		photographerRoutes.PATCH("/:id", ctrl.UpdateOnePackage)
		photographerRoutes.DELETE("/:id", ctrl.DeleteOnePackage)
	}
}
