package dto

import (
	"time"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentRequest struct { // cannot patch subpackage
	StartTime *time.Time `bson:"start_time,omitempty" json:"startTime" ts_type:"string" example:"2025-02-18T10:00:00Z"`
	Location  *string    `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`
}

type AppointmentUpdateStatusRequest struct {
	Status models.AppointmentStatus `bson:"status,omitempty" json:"status" binding:"appointment_status" example:"pending"` // "pending", "accepted", "rejected", "completed"
}

type AppointmentStrictRequest struct {
	StartTime time.Time `bson:"start_time" json:"startTime" ts_type:"string" example:"2025-02-18T10:00:00Z"`
	// Status    models.AppointmentStatus `bson:"status" json:"status" example:"Pending" binding:"appointment_status"` // "pending", "accepted", "rejected", "completed"
	Location string `bson:"location" json:"location" example:"Bangkok, Thailand"`
}

type AppointmentResponse struct {
	ID             primitive.ObjectID       `bson:"_id,omitempty" json:"id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID     primitive.ObjectID       `bson:"customer_id" json:"customerId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1238"`
	PhotographerID primitive.ObjectID       `bson:"photographer_id" json:"photographerId" ts_type:"string" example:"656e2b5e3f1a324d8b9e1236"`
	Package        models.Package           `bson:"package" json:"package" ts_type:"models.Package" example:"{}"`
	Subpackage     models.Subpackage        `bson:"sub_package" json:"subpackage" ts_type:"models.Subpackage" example:"{}"`
	BusyTimeID     primitive.ObjectID       `bson:"busy_time_id" json:"busyTimeId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e7877"`
	Status         models.AppointmentStatus `bson:"status" json:"status" binding:"appointment_status" ts_type:"string" example:"pending"`
	Location       string                   `bson:"location,omitempty" json:"location" ts_type:"string" example:"Bangkok, Thailand"`
	Price          int                      `bson:"price" json:"price" ts_type:"number" example:"1500"`
}

type AppointmentDetail struct {
	ID               primitive.ObjectID       `bson:"_id" json:"id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	Package          models.Package           `bson:"package" json:"package" ts_type:"models.Package" example:"{}"`
	Subpackage       models.Subpackage        `bson:"subpackage" json:"subpackage" ts_type:"model.Subpackage" example:"{}"`
	CustomerID       primitive.ObjectID       `bson:"customer_id" json:"customerId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	PhotographerID   primitive.ObjectID       `bson:"photographer_id" json:"photographerId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	PackageName      string                   `bson:"package_name" json:"packageName" ts_type:"string" example:"Wedding Package"`
	SubpackageName   string                   `bson:"subpackage_name" json:"subpackageName" ts_type:"string" example:"Wedding Subpackage"`
	CustomerName     string                   `bson:"customer_name" json:"customerName" ts_type:"string" example:"John Doe"`
	PhotographerName string                   `bson:"photographer_name" json:"photographerName" ts_type:"string" example:"Jane Smith"`
	Price            int                      `bson:"price" json:"price" ts_type:"number" example:"1500"`
	StartTime        time.Time                `bson:"start_time" json:"startTime" ts_type:"string" example:"2023-10-01T10:00:00Z"`
	EndTime          time.Time                `bson:"end_time" json:"endTime" ts_type:"string" example:"2023-10-01T12:00:00Z"`
	Status           models.AppointmentStatus `bson:"status" json:"status" ts_type:"string" example:"Pending"`
	Location         string                   `bson:"location" json:"location" ts_type:"string" example:"Bangkok, Thailand"`
}

type CreateAppointmentResponse struct {
	Appointment AppointmentResponse `bson:"appointment" json:"appointment" ts_type:"AppointmentResponse"`
	BusyTime    models.BusyTime     `bson:"busy_time" json:"busyTime" ts_type:"BusyTime"`
}

func (req *AppointmentStrictRequest) ToModel(user *models.User, pkg *models.Package, subpackage *models.Subpackage, busyTime *models.BusyTime) *models.Appointment {
	return &models.Appointment{
		ID:             primitive.NewObjectID(),
		CustomerID:     user.ID,
		PhotographerID: pkg.OwnerID,
		Package:        *pkg,
		Subpackage:     *subpackage,
		BusyTimeID:     busyTime.ID,
		Status:         "Pending",
		Location:       req.Location,
		Price:          subpackage.Price,
	}
}

type FilteredAppointmentResponse struct {
	Appointments []AppointmentResponse `json:"appointments"`
	Pagination   Pagination            `json:"pagination"`
}
