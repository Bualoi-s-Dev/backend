package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func AppointmentRoutes(router *gin.Engine, ctrl *controllers.AppointmentController) {
	appointmentGroup := router.Group("/appointment")

	appointmentGroup.GET("", ctrl.GetAllAppointment)
	appointmentGroup.GET("/:id", ctrl.GetAppointmentById)
	appointmentGroup.POST("/:subpackageId", ctrl.CreateAppointment)
	appointmentGroup.PATCH("/:id", ctrl.UpdateAppointment)
	appointmentGroup.PATCH("/status/:id", ctrl.UpdateAppointmentStatus)
	appointmentGroup.DELETE("/:id", ctrl.DeleteAppointment)
}
