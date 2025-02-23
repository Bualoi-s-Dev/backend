package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func SubpackageRoutes(router *gin.Engine, ctrl *controllers.SubpackageController) {
	subpackageRoutes := router.Group("/subpackage")
	{
		subpackageRoutes.GET("", ctrl.GetAllSubpackages)
		subpackageRoutes.POST("/:packageId", ctrl.CreateSubpackage)
		subpackageRoutes.PATCH("/:subpackageId", ctrl.UpdateSubpackage)
		subpackageRoutes.DELETE("/:subpackageId", ctrl.DeleteSubpackage)
	}
}
