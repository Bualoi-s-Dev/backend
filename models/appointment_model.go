package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID     primitive.ObjectID `bson:"customer_id" json:"customer_id" example:"656e2b5e3f1a3c4d8b9e1238"`
	PhotographerID primitive.ObjectID `bson:"photographer_id" json:"photographer_id" example:"656e2b5e3f1a324d8b9e1236"`
	PackageID      primitive.ObjectID `bson:"package_id" json:"package_id" example:"656e2b5e3f1a3c4d8b9e1235"`
	SubPackageID   primitive.ObjectID `bson:"sub_package_id" json:"sub_package_id" example:"sub_package_001"`
	StartTime      time.Time          `bson:"start_time" json:"start_time" example:"2025-02-18T10:00:00Z"`
	EndTime        time.Time          `bson:"end_time" json:"end_time" example:"2025-02-18T12:00:00Z"`
	Status         AppointmentStatus  `bson:"status" json:"status" example:"pending"` // "pending", "accepted", "rejected", "completed"
	Location       string             `bson:"location,omitempty" json:"location" example:"Bangkok, Thailand"`
	// Payment       Payment            `bson:"payment,omitempty" json:"payment,omitempty" example:"{...}"`
}

type AppointmentStatus string

const (
	Pending   AppointmentStatus = "Pending"
	Accepted  AppointmentStatus = "Accepted"
	Rejected  AppointmentStatus = "Rejected"
	Completed AppointmentStatus = "Completed"
)

var ValidAppointmentStatus = []struct {
	Value  AppointmentStatus
	TSName string
}{
	{Pending, string(Pending)},
	{Accepted, string(Accepted)},
	{Rejected, string(Rejected)},
	{Completed, string(Completed)},
}
