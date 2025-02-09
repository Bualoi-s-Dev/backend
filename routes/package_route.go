package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func PackageRoutes(router *gin.Engine, ctrl *controllers.PackageController) {
	packageRoutes := router.Group("/package")
	{
		packageRoutes.GET("/", ctrl.GetAllPackages)
		packageRoutes.GET("/:id", ctrl.GetOnePackage)
		packageRoutes.POST("/", ctrl.CreateOnePackage)
		packageRoutes.PUT("/:id", ctrl.ReplaceOnePackage)
		packageRoutes.DELETE("/:id", ctrl.DeleteOnePackage)
	}
}
