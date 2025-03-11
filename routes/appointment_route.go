package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/gin-gonic/gin"
)

func AppointmentRoutes(router *gin.Engine, ctrl *controllers.AppointmentController) {
	appointmentGroup := router.Group("/appointment")
	commonRoutes := appointmentGroup.Group("", middleware.AllowRoles(models.Photographer, models.Customer))
	{
		commonRoutes.GET("", ctrl.GetAllAppointment)
		commonRoutes.GET("/:id", ctrl.GetAppointmentById)
	}
	customerRoutes := appointmentGroup.Group("", middleware.AllowRoles(models.Customer))
	{
		customerRoutes.POST("/:subpackageId", ctrl.CreateAppointment)
		customerRoutes.PATCH("/:id", ctrl.UpdateAppointment)
		customerRoutes.PATCH("/status/:id", ctrl.UpdateAppointmentStatus)
	}

	appointmentGroup.DELETE("/:id", ctrl.DeleteAppointment)
}
