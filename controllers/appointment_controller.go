package controllers

import (
	"net/http"
	"time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentController struct {
	AppointmentService *services.AppointmentService
}

func NewAppointmentController(appointmentService *services.AppointmentService) *AppointmentController {
	return &AppointmentController{AppointmentService: appointmentService}
}

func HandleError(c *gin.Context, err error) {
	if err == services.ErrBadRequest {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment id"})
		return
	}
	if err == services.ErrUnauthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authorized to access this appointment"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func GetIDFromParam(c *gin.Context) (primitive.ObjectID, error) {
	appointmentID_ := c.Param("id")
	appointmentId, err := primitive.ObjectIDFromHex(appointmentID_)
	if err != nil {
		return primitive.NilObjectID, services.ErrBadRequest
	}
	return appointmentId, nil
}

func (a *AppointmentController) GetAllAppointment(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	if user.Role == "Guest" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Guest cannot access this endpoint"})
		return
	}

	appointments, err := a.AppointmentService.GetAllAppointment(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

func (a *AppointmentController) GetAppointmentById(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	appointmentId, err := GetIDFromParam(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	appointment, err := a.AppointmentService.GetAppointmentById(c.Request.Context(), appointmentId, user)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (a *AppointmentController) CreateAppointment(c *gin.Context) {
	// user
	user := middleware.GetUserFromContext(c)

	if user.Role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can create appointment"})
		return
	}

	// request
	var req dto.AppointmenStrictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	if req.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time must be in the future"})
		return
	}
	//

	appointment, err := a.AppointmentService.CreateOneAppointment(c.Request.Context(), user, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// only update Appointment time and Location Only
func (a *AppointmentController) UpdateAppointment(c *gin.Context) {
	// user
	user := middleware.GetUserFromContext(c)

	appointmentId, err := GetIDFromParam(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	if user.Role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can update appointment"})
		return
	}

	// req
	var req dto.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	// TODO: maybe update requestBody or something later?
	if req.Status != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status cannot be updated at this endpoint"})
		return
	}

	updatedAppointment, err := a.AppointmentService.UpdateAppointment(c.Request.Context(), user, appointmentId, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

func (a *AppointmentController) DeleteAppointment(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	appointmentId, err := GetIDFromParam(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := a.AppointmentService.DeleteAppointment(c.Request.Context(), appointmentId, user); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment was deleted successfully"})
}
