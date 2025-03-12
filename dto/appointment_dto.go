package dto

import (
	"time"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentRequest struct { // cannot patch subpackage
	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time" ts_type:"string" example:"2025-02-18T10:00:00Z"`
	Location  *string    `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`
}

type AppointmentUpdateStatusRequest struct {
	Status models.AppointmentStatus `bson:"status,omitempty" json:"status" binding:"appointment_status" example:"pending"` // "pending", "accepted", "rejected", "completed"
}

type AppointmentStrictRequest struct {
	StartTime time.Time `bson:"start_time" json:"start_time" ts_type:"string" example:"2025-02-18T10:00:00Z"`
	// Status    models.AppointmentStatus `bson:"status" json:"status" example:"Pending" binding:"appointment_status"` // "pending", "accepted", "rejected", "completed"
	Location string `bson:"location" json:"location" example:"Bangkok, Thailand"`
}

type AppointmentResponse struct {
	ID             string `bson:"_id,omitempty" json:"id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID     string `bson:"user_id" json:"customer_id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e5678"`
	PhotographerID string `bson:"photographer_id" json:"photographer_id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e91011"`
	PackageID      string `bson:"package_id" json:"package_id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e6969"`
	SubPackageID   string `bson:"sub_package_id" json:"sub_package_id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e4205"`

	StartTime time.Time `bson:"start_time" json:"start_time" ts_type:"string" example:"2025-02-18T10:00:00Z"`
	EndTime   time.Time `bson:"end_time" json:"end_time" ts_type:"string" example:"2025-02-18T12:00:00Z"`
	Status    string    `bson:"status" json:"status" example:"Pending"` // "pending", "accepted", "rejected", "completed"
	Location  string    `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`
	// Payment       Payment   `bson:"payment,omitempty" json:"payment,omitempty" example:"{...}"`
}

func (req *AppointmentStrictRequest) ToModel(user *models.User, pkg *models.Package, subpackage *models.Subpackage, busyTime *models.BusyTime) *models.Appointment {
	return &models.Appointment{
		ID:             primitive.NewObjectID(),
		CustomerID:     user.ID,
		PhotographerID: pkg.OwnerID,
		PackageID:      pkg.ID,
		SubpackageID:   subpackage.ID,
		BusyTimeID:     busyTime.ID,
		Status:         "Pending",
		Location:       req.Location,
		Price:          subpackage.Price,
	}
}
