package controllers

import (
	"net/http"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	AppointmentService *services.AppointmentService
}

func NewAppointmentController(appointmentService *services.AppointmentService) *AppointmentController {
	return &AppointmentController{AppointmentService: appointmentService}
}

func (a *AppointmentController) GetAllAppointment(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	appointments, err := a.AppointmentService.GetAllAppointment(c, string(ugdser.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

func (a *AppointmentController) GetAppointmentById(c *gin.Context) {
	id := c.Param("id")

	appointment, err := a.AppointmentService.GetAppointmentById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func (a *AppointmentController) CreateAppointment(c *gin.Context) {
	var appointmentRequest dto.AppointmentRequest
	if err := c.ShouldBindJSON(&appointmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.AppointmentService.CreateAppointment(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated)
}

func (a *AppointmentController) UpdateAppointment(c *gin.Context) {
	id := c.Param("id")

	var appointmentRequest dto.AppointmentRequest
	if err := c.ShouldBindJSON(&appointmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.AppointmentService.UpdateAppointment(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK)
}

func (a *AppointmentController) DeleteAppointment(c *gin.Context) {
	id := c.Param("id")

	err := a.AppointmentService.DeleteAppointment(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}
