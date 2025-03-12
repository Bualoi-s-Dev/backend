package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusyTime struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	PhotographerID primitive.ObjectID `bson:"photographer_id,omitempty" json:"photographerId" ts_type:"string" example:"12345678abcd"`
	Type           BusyTimeType       `bson:"type,omitempty" json:"type" binding:"omitempty,busy_time_type" example:"Photographer"`
	StartTime      time.Time          `bson:"start_time,omitempty" json:"startTime" ts_type:"string" example:"2025-02-23T10:00:00Z"`
	EndTime        time.Time          `bson:"end_time,omitempty" json:"endTime" ts_type:"string" example:"2025-02-23T12:00:00Z"`
	IsValid        bool               `bson:"is_valid" json:"isValid" ts_type:"boolean" example:"true"`
}

type BusyTimeType string

const (
	TypePhotographer BusyTimeType = "Photographer"
	TypeAppointment  BusyTimeType = "Appointment"
)

var ValidBusyTimeTypes = []struct {
	Value  BusyTimeType
	TSName string
}{
	{TypePhotographer, string(TypePhotographer)},
	{TypeAppointment, string(TypeAppointment)},
}
