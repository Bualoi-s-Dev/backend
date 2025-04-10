package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

func BusyTimeRoutes(router *gin.Engine, ctrl *controllers.BusyTimeController, userService *services.UserService) {
	busyTimeRoutes := router.Group("/busytime", middleware.AllowRoles(userService, models.Photographer, models.Customer))
	{
		busyTimeRoutes.GET("/photographer/:photographerId", ctrl.GetBusyTimesByPhotographerId)
	}
}
