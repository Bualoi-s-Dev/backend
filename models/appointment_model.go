package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Implement Appointment and BusyTime pair struct
type Appointment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID     primitive.ObjectID `bson:"customer_id" json:"customerId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1238"`
	PhotographerID primitive.ObjectID `bson:"photographer_id" json:"photographerId" ts_type:"string" example:"656e2b5e3f1a324d8b9e1236"`
	PackageID      primitive.ObjectID `bson:"package_id" json:"packageId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1235"`
	SubpackageID   primitive.ObjectID `bson:"sub_package_id" json:"subpackageId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e9999"`
	BusyTimeID     primitive.ObjectID `bson:"busy_time_id" json:"busyTimeId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e7877"`
	Status         AppointmentStatus  `bson:"status" json:"status" binding:"appointment_status" ts_type:"string" example:"pending"`
	Location       string             `bson:"location,omitempty" json:"location" ts_type:"string" example:"Bangkok, Thailand"`
	Price          int                `bson:"price" json:"price" ts_type:"number" example:"1500"`
	// Payment       Payment            `bson:"payment,omitempty" json:"payment,omitempty" example:"{...}"`
}

type AppointmentStatus string

const (
	AppointmentPending   AppointmentStatus = "Pending"
	AppointmentAccepted  AppointmentStatus = "Accepted"
	AppointmentRejected  AppointmentStatus = "Rejected"
	AppointmentCanceled  AppointmentStatus = "Canceled"
	AppointmentCompleted AppointmentStatus = "Completed"
)

var ValidAppointmentStatus = []struct {
	Value  AppointmentStatus
	TSName string
}{
	{AppointmentPending, string(AppointmentPending)},
	{AppointmentAccepted, string(AppointmentAccepted)},
	{AppointmentRejected, string(AppointmentRejected)},
	{AppointmentCanceled, string(AppointmentCanceled)},
	{AppointmentCompleted, string(AppointmentCompleted)},
}
