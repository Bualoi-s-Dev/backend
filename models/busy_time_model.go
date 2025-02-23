package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BusyTime struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	PhotographerID  primitive.ObjectID `bson:"photographer_id,omitempty" json:"photographerId" ts_type:"string" example:"12345678abcd"`
	Type            BusyTimeType       `bson:"type,omitempty" json:"type" example:"shooting, editing, meeting"`
	StartTime       time.Time          `bson:"start_time,omitempty" json:"startTime" example:"2025-02-23T10:00:00Z"`
	EndTime         time.Time          `bson:"end_time,omitempty" json:"endTime" example:"2025-02-23T12:00:00Z"`
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