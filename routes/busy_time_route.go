package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/gin-gonic/gin"
)

func BusyTimeRoutes(router *gin.Engine, ctrl *controllers.BusyTimeController) {
	busyTimeRoutes := router.Group("/busytime", middleware.AllowRoles(models.Photographer, models.Customer))
	{
		busyTimeRoutes.GET("/photographer/:photographerId", ctrl.GetBusyTimesByPhotographerId)
	}
}
