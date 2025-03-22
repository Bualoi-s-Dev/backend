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
// @Param status query string false "Status of appointment"
// @Param availableStartDay query string false "Available start day of appointment"
// @Param availableEndDay query string false "Available end day of appointment"
// @Param name query string false "Name of package title or customer name"
// @Param minPrice query string false "Minimum price of appointment"
// @Param maxPrice query string false "Maximum price of appointment"
// @Success 200 {array} dto.AppointmentResponse
// @Failure 401 {object} string "Unauthorized"
// @Router /appointment [get]
func (a *AppointmentController) GetAllAppointment(c *gin.Context) {
	// Get query parameters for filtering
	filters := map[string]string{
		"status":            c.Query("status"),
		"availableStartDay": c.Query("availableStartDay"),
		"availableEndDay":   c.Query("availableEndDay"),
		"name":              c.Query("name"),
		"minPrice":          c.Query("minPrice"),
		"maxPrice":          c.Query("maxPrice"),
	}

	// Verify date format
	dateFormat := "2006-01-02"
	if filters["availableStartDay"] != "" {
		if _, err := time.Parse(dateFormat, filters["availableStartDay"]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format for availableStartDay. Use YYYY-MM-DD."})
			return
		}
	}
	if filters["availableEndDay"] != "" {
		if _, err := time.Parse(dateFormat, filters["availableEndDay"]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format for availableEndDay. Use YYYY-MM-DD."})
			return
		}
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

// GetAllAppointmentDetail godoc
// @Tags Appointment
// @Summary Get all available appointments with provided details
// @Description Retrieve all available appointments detail that the user can see from the database
// @Success 200 {array} dto.AppointmentDetailResponse
// @Failure 401 {object} string "Unauthorized"
// @Router /appointment/detail [get]
func (a *AppointmentController) GetAllAppointmentDetail(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	appointmentDetails, err := a.AppointmentService.GetAllAppointmentDetail(c, user)
	if err != nil {
		apperrors.HandleError(c, err, "Error while get all appointment detail")
		return
	}
	c.JSON(http.StatusOK, appointmentDetails)
}

// GetAllAppointmentDetailById godoc
// @Tags Appointment
// @Summary Get an appointments with provided details by Id
// @Description Retrieve all available appointments detail that the user can see from the database
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.AppointmentDetailResponse
// @Failure 401 {object} string "Unauthorized"
// @Router /appointment/detail/{id} [get]
func (a *AppointmentController) GetAllAppointmentDetailById(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	appointmentId, err := getIDFromParam(c)
	if err != nil {
		apperrors.HandleError(c, err, "Cannot get appointmentId from param.")
		return
	}

	appointment, err := a.AppointmentService.GetAppointmentById(c.Request.Context(), user, appointmentId)
	appointmentDetail, err := a.AppointmentService.GetAppointmentDetailById(c, user, appointment)
	if err != nil {
		apperrors.HandleError(c, err, "Error while get all appointment detail")
		return
	}
	c.JSON(http.StatusOK, appointmentDetail)
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
// @Success 200 {object} dto.CreateAppointmentResponse
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
	if req.Status == models.AppointmentPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update status to Pending"})
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

	if appointment.Status == models.AppointmentAccepted && req.Status != models.AppointmentCanceled {
		apperrors.HandleError(c, apperrors.ErrBadRequest, "Can change status to Canceled only if status is Accepted")
		return
	}

	// TODO: refactor in route later
	// Customer can't accept appointment
	if user.Role == models.Customer && req.Status == models.AppointmentAccepted {
		apperrors.HandleError(c, apperrors.ErrUnauthorized, "Customer can't accepted an appointment.")
		return
	}

	validStatus := false
	if appointment.Status == models.AppointmentPending && req.Status == models.AppointmentAccepted { // accept pending status
		validStatus = true // from "Pending" to "Accepted" => isValid = true (reserve a photographer busyTime)
	}
	busyTime.IsValid = validStatus

	minimumAvailableCanceledTime := busyTime.StartTime.Add(-24 * time.Hour)
	if req.Status == models.AppointmentCanceled && time.Now().After(minimumAvailableCanceledTime) {
		apperrors.HandleError(c, apperrors.ErrBadRequest, "Cannot cancel an appointment before 24 hours")
		return
	}

	if err := a.BusyTimeService.UpdateValidStatus(c.Request.Context(), busyTime); err != nil {
		apperrors.HandleError(c, err, "(Update Status) Could not update busyTime")
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
