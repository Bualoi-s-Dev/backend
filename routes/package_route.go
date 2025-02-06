package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func PackageRoutes(router *gin.Engine, ctrl *controllers.PackageController) {
	router.GET("/package", ctrl.GetAllPackages)
	router.GET("/package/:id", ctrl.GetOnePackage)
	router.POST("/package", ctrl.CreateOnePackage)
	router.PATCH("/package/:id", ctrl.UpdateOnePackage)
	router.DELETE("/package/:id", ctrl.DeleteOnePackage)
}
