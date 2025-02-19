package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
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

func (a *AppointmentController) GetAllAppointment(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	if user.Role == "Guest" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Guest cannot access this endpoint"})
		return
	}

	appointments, err := a.AppointmentService.GetAllAppointment(c, user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

func (a *AppointmentController) GetAppointmentById(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	appointmentID_ := c.Param("id")
	appointmentID, err := primitive.ObjectIDFromHex(appointmentID_)
	if err != nil {
		fmt.Println("Invalid appointmentID:", err)
		return
	}

	appointment, err := a.AppointmentService.GetAppointmentById(c, appointmentID, user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func (a *AppointmentController) CreateAppointment(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	if user.Role != "Customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only customer can create appointment"})
		return
	}

	var itemInput dto.AppointmenStrictRequest
	if err := c.ShouldBindJSON(&itemInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, " + err.Error()})
		return
	}

	if itemInput.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time must be in the future"})
		return
	}

	// TODO: Verify the appointment later

	subpackage, err := a.AppointmentService.FindSubpackageByID(c, itemInput.SubPackageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pkg, err := a.AppointmentService.FindPackageByID(c, subpackage.PackageID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	appointment := models.Appointment{
		ID:             primitive.NewObjectID(),
		CustomerID:     user.ID,
		PhotographerID: pkg.OwnerID,
		StartTime:      itemInput.StartTime,
		EndTime:        itemInput.StartTime.Add(time.Duration(subpackage.Duration) * time.Minute),
		SubPackageID:   itemInput.SubPackageID,
		Status:         "Pending",
		Location:       itemInput.Location,
	}

	item, err := a.AppointmentService.CreateOneAppointment(c, &appointment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// func (a *AppointmentController) UpdateAppointment(c *gin.Context) {
// 	id := c.Param("id")

// 	var appointmentRequest dto.AppointmentRequest
// 	if err := c.ShouldBindJSON(&appointmentRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := a.AppointmentService.UpdateAppointment(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, appointment)
// }

// func (a *AppointmentController) DeleteAppointment(c *gin.Context) {
// 	id := c.Param("id")

// 	err := a.AppointmentService.DeleteAppointment(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
// }
