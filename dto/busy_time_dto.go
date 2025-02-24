package dto

import (
	"time"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusyTimeRequest struct {
	Type      *models.BusyTimeType `bson:"type" json:"type" binding:"required" example:"PHOTOGRAPHER"`
	StartTime *time.Time           `bson:"start_time" json:"startTime" binding:"required" example:"2025-02-23T10:00:00Z"`
	EndTime   *time.Time           `bson:"end_time" json:"endTime" binding:"required" example:"2025-02-23T12:00:00Z"` // TODO: remove json body later <- this only gen from start time + subpackage duration
	IsValid   *bool                `bson:"is_valid" json:"isValid" example:"true"`
}

func (item *BusyTimeRequest) ToModel(photographerID primitive.ObjectID) *models.BusyTime {
	return &models.BusyTime{
		PhotographerID: photographerID,
		Type:           *item.Type,
		StartTime:      *item.StartTime,
		EndTime:        *item.EndTime,
		IsValid:        *item.IsValid,
	}
}
