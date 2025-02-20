package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentRequest struct { // cannot patch subpackage
	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time" example:"2025-02-18T10:00:00Z"`
	EndTime   *time.Time `bson:"-" json:"-" example:"2025-02-18T12:00:00Z"`        // FIXME: remove this field (or not)?
	Status    *string    `bson:"status,omitempty" json:"status" example:"pending"` // "pending", "accepted", "rejected", "completed"
	Location  *string    `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`
}

type AppointmenStrictRequest struct {
	SubPackageID primitive.ObjectID `bson:"sub_package_id" json:"sub_package_id" example:"507f1f77bcf86cd799439011"`

	StartTime time.Time `bson:"start_time" json:"start_time" example:"2025-02-18T10:00:00Z"`
	Status    string    `bson:"status" json:"status" example:"pending"` // "pending", "accepted", "rejected", "completed"
	Location  string    `bson:"location" json:"location" example:"Bangkok, Thailand"`
}

type AppointmentResponse struct {
	ID             string `bson:"_id,omitempty" json:"id,omitempty" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID     string `bson:"user_id" json:"customer_id" example:"customer_123"`
	PhotographerID string `bson:"photographer_id" json:"photographer_id" example:"photographer_456"`
	PackageID      string `bson:"package_id" json:"package_id" example:"package_789"`
	SubPackageID   string `bson:"sub_package_id" json:"sub_package_id" example:"sub_package_001"`

	StartTime time.Time `bson:"start_time" json:"start_time" example:"2025-02-18T10:00:00Z"`
	EndTime   time.Time `bson:"end_time" json:"end_time" example:"2025-02-18T12:00:00Z"`
	Status    string    `bson:"status" json:"status" example:"pending"` // "pending", "accepted", "rejected", "completed"
	Location  string    `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`
	// Payment       Payment   `bson:"payment,omitempty" json:"payment,omitempty" example:"{...}"`
}
