package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

func SubpackageRoutes(router *gin.Engine, ctrl *controllers.SubpackageController, userService *services.UserService) {
	subpackageRoutes := router.Group("/subpackage")
	commonRoutes := subpackageRoutes.Group("", middleware.AllowRoles(userService, models.Photographer, models.Customer))
	{
		commonRoutes.GET("", ctrl.GetAllSubpackages)
		commonRoutes.GET("/:id", ctrl.GetByIdSubpackages)
	}
	photographerRoutes := subpackageRoutes.Group("", middleware.AllowRoles(userService, models.Photographer))
	{

		photographerRoutes.POST("/:packageId", ctrl.CreateSubpackage)
		photographerRoutes.PATCH("/:id", ctrl.UpdateSubpackage)
		photographerRoutes.DELETE("/:id", ctrl.DeleteSubpackage)
	}
}
