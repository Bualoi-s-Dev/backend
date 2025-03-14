package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Bualoi-s-Dev/backend/apperrors"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	// Get query parameters for filtering
	filters := map[string]string{
		"status":            c.Query("status"),
		"availableStartDay": c.Query("availableStartDay"),
		"availableEndDay":   c.Query("availableEndDay"),
	}

	user := middleware.GetUserFromContext(c)

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Retrieve and filter subpackages with pagination
	appointments, err := a.AppointmentService.GetFilteredAppointments(c.Request.Context(), user, filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (a *AppointmentController) GetAllAppointmentDetail(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	appointmentDetails, err := a.AppointmentService.GetAllAppointmentDetail(c, user)
	if err != nil {
		apperrors.HandleError(c, err, "Error while get all appointment detail")
		return
	}
	c.JSON(http.StatusOK, appointmentDetails)
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
// @Summary Create appointment
// @Description Create a new appointment from a specific subpackage
// @Param subpackageId path string true "Subpackage ID"
// @Body {AppointmenStrictRequest} request body "Create Appointment Request"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/{subpackageId} [post]
func (a *AppointmentController) CreateAppointment(c *gin.Context) {
	// user
	user := middleware.GetUserFromContext(c)

	loc, _ := time.LoadLocation("Asia/Bangkok")
	// request

	subpackageId, err := getSubpackageIDFromParam(c)
	if err != nil {
		fmt.Println("Cannot get Subpackage ID From Param")
		apperrors.HandleError(c, err, "Cannot get subpackageId from param")
		return
	}
	var req dto.AppointmentStrictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	if req.StartTime.Before(time.Now().In(loc)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time must be in the future"})
		return
	}

	busyTimeType := models.TypeAppointment
	busyTimeReq := &dto.BusyTimeStrictRequest{
		Type:      busyTimeType,
		StartTime: req.StartTime,
		IsValid:   false,
	}

	busyTime, err := a.BusyTimeService.CreateFromSubpackage(c.Request.Context(), busyTimeReq, subpackageId)
	if err != nil {
		fmt.Println("Cannot Create Busy Time from Subpackage!!")
		apperrors.HandleError(c, err, "Cannot create BusyTime")
		return
	}

	appointment, err := a.AppointmentService.CreateOneAppointment(c.Request.Context(), user, subpackageId, busyTime, &req)
	if err != nil {
		fmt.Println("Cannot create Appointment!!!")
		apperrors.HandleError(c, err, "Cannot create this appointment")
		return
	}

	c.JSON(http.StatusCreated, bson.M{"appointment": appointment, "busyTime": busyTime})
}

/*
// UpdateAppointment godoc
// @Tags Appointment
// @Summary Update appointment
// @Description Update the time and location of a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/{id} [patch]
func (a *AppointmentController) UpdateAppointment(c *gin.Context) {
	// user
	user := middleware.GetUserFromContext(c)

	appointmentId, err := getIDFromParam(c)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get appointmentId from param")
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

	busyTime, err := a.BusyTimeService.GetById(c.Request.Context(), appointment.BusyTimeID.Hex())
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get the busyTime from this busyTimeId")
		return
	}
	if busyTime.IsValid { // status is pending and current busyTime is valid => impossible
		apperrors.HandleError(c, apperrors.ErrInternalServer, "busyTime cannot be valid if current appointment status is pending")
		return
	}

	// req
	var req dto.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	// update time => need to update busy time too
	if req.StartTime != nil {
		if err := a.BusyTimeService.Delete(c, busyTime.ID.Hex()); err != nil {
			apperrors.HandleError(c, err, "(Update Status) Could not delete appointment before update")
		}

		// calculate endTime from busyTime and startTime
		loc, _ := time.LoadLocation("Asia/Bangkok")
		currentTime := time.Now().In(loc)
		if req.StartTime.Before(currentTime) {
			apperrors.HandleError(c, apperrors.ErrBadRequest, "Cannot channge appointment start time before current Time: "+currentTime.String())
			return
		}

		// calculate duration Endtime - StartTime (from appointment)
		duration := busyTime.EndTime.Sub(busyTime.StartTime)
		endTime := req.StartTime.Add(duration)

		busyTimeReq := &dto.BusyTimeRequest{
			Type:      &busyTime.Type,
			StartTime: req.StartTime,
			EndTime:   &endTime,
			IsValid:   &busyTime.IsValid, // false
		}

		if err := a.BusyTimeService.Create(c, busyTimeReq, appointment.PhotographerID); err != nil {
			// Error becuase the new startTime is conflict with the photgrapher's BusyTime
			busyTimeReq.StartTime = &busyTime.StartTime
			busyTimeReq.EndTime = &busyTime.EndTime
			if err := a.BusyTimeService.Create(c, busyTimeReq, appointment.PhotographerID); err != nil {
				apperrors.HandleError(c, apperrors.ErrInternalServer, "(On update status) Can't recreate the same busyTime after deleted") // shouldn't happen
				return
			}
			apperrors.HandleError(c, err, "(Update Status) Could not re-create appointment")
			return
		}
		// } ENDTODO:

	}

	updatedAppointment, err := a.AppointmentService.UpdateAppointment(c.Request.Context(), user, appointment, &req)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot update this appointment")
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}
*/

// UpdateAppointmentStatus godoc
// @Tags Appointment
// @Summary Update appointment status
// @Description Update the status of a specific appointment by its ID
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentResponse
// @Failure 400 {object} string "Invalid appointment id"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /appointment/status/{id} [patch]
func (a *AppointmentController) UpdateAppointmentStatus(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	appointmentId, err := getIDFromParam(c)

	if err != nil {
		apperrors.HandleError(c, err, "Cannot get appointmentId from param.")
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
		return
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

	// TODO: refactor in route later
	if user.Role == models.Customer && req.Status == models.AppointmentAccepted {
		apperrors.HandleError(c, apperrors.ErrUnauthorized, "Customer can't accepted an appointment.")
		return
	}

	var validStatus bool
	if appointment.Status == models.AppointmentAccepted && req.Status == models.AppointmentCanceled {
		validStatus = false
	} else if appointment.Status == models.AppointmentPending && req.Status == models.AppointmentAccepted {
		validStatus = true // from "Pending" to "Accepted" => isValid = true (reserve a photographer busyTime)
	}

	oldID := busyTime.ID
	if err := a.BusyTimeService.Delete(c, oldID.Hex()); err != nil {
		apperrors.HandleError(c, err, "(Update Status) Could not delete appointment before update")
	}

	busyTimeReq := &dto.BusyTimeStrictRequest{
		Type:      busyTime.Type,
		StartTime: busyTime.StartTime,
		IsValid:   validStatus,
	}

	if err := a.BusyTimeService.CreateForUpdate(c, busyTimeReq, oldID, appointment.PhotographerID); err != nil {
		apperrors.HandleError(c, err, "(Update Status) Could not re-create appointment")
		return
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
