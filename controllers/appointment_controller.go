package controllers

import (
	"net/http"
	"time"

	"github.com/Bualoi-s-Dev/backend/apperrors"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Change some error raising to forbidden (instead of Unauthorized)

type AppointmentController struct {
	AppointmentService *services.AppointmentService
	BusyTimeService    *services.BusyTimeService
}

func NewAppointmentController(appointmentService *services.AppointmentService, busyTimeService *services.BusyTimeService) *AppointmentController {
	return &AppointmentController{
		AppointmentService: appointmentService,
		BusyTimeService:    busyTimeService,
	}
}

func getIDFromParam(c *gin.Context) (primitive.ObjectID, error) {
	appointmentId_ := c.Param("id")
	appointmentId, err := primitive.ObjectIDFromHex(appointmentId_)
	if err != nil {
		return primitive.NilObjectID, apperrors.ErrBadRequest
	}
	return appointmentId, nil
}

func getSubpackageIDFromParam(c *gin.Context) (primitive.ObjectID, error) {
	subpackageId_ := c.Param("subpackageId")
	subpackageId, err := primitive.ObjectIDFromHex(subpackageId_)
	if err != nil {
		return primitive.NilObjectID, apperrors.ErrBadRequest
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
	role := middleware.GetUserRoleFromContext(c)

	if role == "Guest" {
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
	role := middleware.GetUserRoleFromContext(c)

	if role == "Guest" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Guest cannot access this endpoint"})
		return
	}

	appointmentId, err := getIDFromParam(c)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get appointmentId from param.")
		return
	}

	appointment, err := a.AppointmentService.GetAppointmentById(c.Request.Context(), user, appointmentId)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get the appointment from this id")
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
	role := middleware.GetUserRoleFromContext(c)

	if role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can create appointment"})
		return
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	// request

	subpackageId, err := getSubpackageIDFromParam(c)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get subpackageId from param")
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

	busyTimeType := models.TypePhotographer
	isValid := false
	busyTimeReq := &dto.BusyTimeRequest{
		Type:      &busyTimeType,
		StartTime: &req.StartTime,
		IsValid:   &isValid,
	}

	busyTime, err := a.BusyTimeService.CreateFromSubpackage(c.Request.Context(), busyTimeReq, subpackageId)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot create BusyTime")
	}

	appointment, err := a.AppointmentService.CreateOneAppointment(c.Request.Context(), user, subpackageId, busyTime, &req)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot create this appointment")
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
	// FIXME: I Forgot but something(s) must be fixed
	// user
	user := middleware.GetUserFromContext(c)
	role := middleware.GetUserRoleFromContext(c)

	appointmentId, err := getIDFromParam(c)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get appointmentId from param")
		return
	}

	if role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can update appointment properties"})
		return
	}
	appointment, err := a.AppointmentService.GetAppointmentById(c.Request.Context(), user, appointmentId)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get the appointment from this id")
		return
	}
	if appointment.Status != "Pending" {
		apperrors.HandleError(c, apperrors.ErrBadRequest, "Cannot update an appointment if its status is not Pending")
		return
	}

	// req
	var req dto.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	updatedAppointment, err := a.AppointmentService.UpdateAppointment(c.Request.Context(), user, appointmentId, &req)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot update this appointment")
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
	role := middleware.GetUserRoleFromContext(c)
	appointmentId, err := getIDFromParam(c)

	if err != nil {
		apperrors.HandleError(c, err, "Cannot get appointmentId from param.")
		return
	}
	if role == models.Guest {
		apperrors.HandleError(c, apperrors.ErrUnauthorized, "Guest cannot update any appointment")
		return
	}

	var req dto.AppointmentUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	// cannot update any status to complete, it is done via AutoUpdate
	if req.Status == models.AppointmentCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update status to Completed directly"})
		return
	}

	appointment, err := a.AppointmentService.GetAppointmentById(c.Request.Context(), user, appointmentId)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get the appointment from this appointmentId")
	}

	busyTime, err := a.BusyTimeService.GetById(c.Request.Context(), appointment.BusyTimeID.Hex())
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get the busyTime from this busyTimeId")
		return
	}

	// cannot change if it an terminal status
	if appointment.Status == models.AppointmentCanceled || appointment.Status == models.AppointmentCompleted || appointment.Status == models.AppointmentRejected {
		apperrors.HandleError(c, apperrors.ErrAppointmentStatusInvalid, "Cannot change terminated status")
		return
	}

	if appointment.Status == models.AppointmentAccepted && req.Status == models.AppointmentCanceled { // cannot canceled when appointment has begun
		// If change from accepted to canceled => must change isValid of busy time to false

		// TODO: Change this to function later { <--- all maybe create a **BusyTimeService.Update**
		// Note this just update only isValid
		if err := a.BusyTimeService.Delete(c, busyTime.ID.Hex()); err != nil {
			apperrors.HandleError(c, err, "(Update Status) Could not delete appointment before update")
		}

		busyTimeIsValid := false
		busyTimeReq := &dto.BusyTimeRequest{
			Type:      &busyTime.Type,
			StartTime: &busyTime.StartTime,
			IsValid:   &busyTimeIsValid,
		}

		if err := a.BusyTimeService.Create(c, busyTimeReq, appointment.PhotographerID); err != nil {
			apperrors.HandleError(c, err, "(Update Status) Could not re-create appointment")
			return
		}
		// } ENDTODO:

	}

	// TODO: if change from "pending" -> "accepted" (only phtoographer) => must change isValid of busy time to true
	if appointment.Status == models.AppointmentPending && req.Status == models.AppointmentAccepted {
		// check availability
		if role != models.Photographer {
			apperrors.HandleError(c, apperrors.ErrForbidden, "customer cannot change ")
		}
		// TODO: Change this to function later { <--- all maybe create a **BusyTimeService.Update**
		// Note this just update only isValid
		if err := a.BusyTimeService.Delete(c, busyTime.ID.Hex()); err != nil {
			apperrors.HandleError(c, err, "(Update Status) Could not delete appointment before update")
		}

		busyTimeIsValid := true
		busyTimeReq := &dto.BusyTimeRequest{
			Type:      &busyTime.Type,
			StartTime: &busyTime.StartTime,
			IsValid:   &busyTimeIsValid,
		}

		if err := a.BusyTimeService.Create(c, busyTimeReq, appointment.PhotographerID); err != nil {
			apperrors.HandleError(c, err, "(Update Status) Could not re-create appointment")
			return
		}
		// } ENDTODO:

	}

	updatedAppointment, err := a.AppointmentService.UpdateAppointmentStatus(c.Request.Context(), user, appointment, &req)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot update this appointment status")
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
		apperrors.HandleError(c, err, "Cannot get appointmentId from param.")
		return
	}

	if err := a.AppointmentService.DeleteAppointment(c.Request.Context(), appointmentId, user); err != nil {
		apperrors.HandleError(c, err, "Cannot delete this appointment.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment was deleted successfully"})
}
