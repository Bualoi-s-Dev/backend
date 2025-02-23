package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func BusyTimeRoutes(router *gin.Engine, ctrl *controllers.BusyTimeController) {
	busyTimeRoutes := router.Group("/busytime")
	{
		busyTimeRoutes.GET("/photographer/:photographerId", ctrl.GetBusyTimesByPhotographerId)
	}
}
