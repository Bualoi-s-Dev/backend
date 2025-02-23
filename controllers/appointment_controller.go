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
	if err == services.ErrStatusInvalid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	if err == services.ErrStatusTime {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update status because the appointment has been started"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func getIDFromParam(c *gin.Context) (primitive.ObjectID, error) {
	appointmentId_ := c.Param("id")
	appointmentId, err := primitive.ObjectIDFromHex(appointmentId_)
	if err != nil {
		return primitive.NilObjectID, services.ErrBadRequest
	}
	return appointmentId, nil
}

func getSubpackageIDFromParam(c *gin.Context) (primitive.ObjectID, error) {
	subpackageId_ := c.Param("subpackageId")
	subpackageId, err := primitive.ObjectIDFromHex(subpackageId_)
	if err != nil {
		return primitive.NilObjectID, services.ErrBadRequest
	}
	return subpackageId, nil
}

// GetAllAppointment godoc
// @Tags Appointment
// @Summary Get all available appointments
// @Description Retrieve all available appointments that the user can see from the database
// @Success 200 {array} dto.AppointmentResponse
// @Failure 401 {object} string "Unauthorized"
// @Router /appointment [get]
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

// GetAppointmentById godoc
// @Tags Appointment
// @Summary Get appointment by ID
// @Description Retrieve a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/{id} [get]
func (a *AppointmentController) GetAppointmentById(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	if user.Role == "Guest" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Guest cannot access this endpoint"})
		return
	}

	appointmentId, err := getIDFromParam(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	appointment, err := a.AppointmentService.GetAppointmentById(c.Request.Context(), user, appointmentId)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// GetAppointmentById godoc
// @Tags Appointment
// @Summary Get appointment by ID
// @Description Retrieve a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/{id} [get]
func (a *AppointmentController) CreateAppointment(c *gin.Context) {
	// user
	user := middleware.GetUserFromContext(c)

	if user.Role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can create appointment"})
		return
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	// request

	subpackageId, err := getSubpackageIDFromParam(c)
	if err != nil {
		HandleError(c, err)
	}
	var req dto.AppointmenStrictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	if req.StartTime.Before(time.Now().In(loc)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time must be in the future"})
		return
	}
	//

	appointment, err := a.AppointmentService.CreateOneAppointment(c.Request.Context(), user, subpackageId, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// UpdateAppointment godoc
// @Tags Appointment
// @Summary Update appointment
// @Description Update the time and location of a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/{id} [put]
func (a *AppointmentController) UpdateAppointment(c *gin.Context) {
	// user
	user := middleware.GetUserFromContext(c)

	appointmentId, err := getIDFromParam(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	if user.Role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can update appointment"})
		return
	}
	// TODO: can be editted only while status is "Pending"

	// req
	var req dto.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	updatedAppointment, err := a.AppointmentService.UpdateAppointment(c.Request.Context(), user, appointmentId, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

// UpdateAppointmentStatus godoc
// @Tags Appointment
// @Summary Update appointment status
// @Description Update the status of a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/status/{id} [put]
func (a *AppointmentController) UpdateAppointmentStatus(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	appointmentId, err := getIDFromParam(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req dto.AppointmentUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	if req.Status == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	updatedAppointment, err := a.AppointmentService.UpdateAppointmentStatus(c.Request.Context(), user, appointmentId, &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

// DeleteAppointment godoc
// @Tags Appointment
// @Summary Delete appointment
// @Description Delete a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} string "Appointment was deleted successfully"
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/{id} [delete]
func (a *AppointmentController) DeleteAppointment(c *gin.Context) { // only admin
	user := middleware.GetUserFromContext(c)

	appointmentId, err := getIDFromParam(c)
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
