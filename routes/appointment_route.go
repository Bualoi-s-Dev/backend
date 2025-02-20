package routes

/*
TODO: Implement the following routes:

Appointment Group (/appointment)
- [x] GET           : get all appointments of user
- [x] GET /:id      : get appointment by id
- [x] POST          : post appointment
- [ ] PATCH /:id    : update appointment
- [ ] PATCH /:id    : photographer accept/reject appointment
- [x] DELETE /:id   : delete appointment
*/

import (
	"os"

	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func AppointmentRoutes(router *gin.Engine, ctrl *controllers.AppointmentController) {
	if os.Getenv("APP_MODE") == "development" {
		appointmentGroup := router.Group("/appointment")

		appointmentGroup.GET("", ctrl.GetAllAppointment)
		appointmentGroup.GET("/:id", ctrl.GetAppointmentById)
		appointmentGroup.POST("", ctrl.CreateAppointment)
		appointmentGroup.PATCH("/:id", ctrl.UpdateAppointment)
		appointmentGroup.PATCH("/status/:id", ctrl.UpdateAppointmentStatus)
		appointmentGroup.DELETE("/:id", ctrl.DeleteAppointment)
	}
}
