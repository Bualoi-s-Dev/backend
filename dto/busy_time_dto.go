package dto

import (
	"time"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusyTimeStrictRequest struct {
	Name      string              `bson:"name,omitempty" json:"name" binding:"required" example:"Vacation"`
	Type      models.BusyTimeType `bson:"type" json:"type" binding:"required" example:"Photographer"`
	StartTime time.Time           `bson:"start_time" json:"startTime" binding:"required" ts_type:"string" example:"2025-02-23T10:00:00Z"`
	EndTime   time.Time           `bson:"end_time" json:"endTime" binding:"required" ts_type:"string" example:"2025-02-23T12:00:00Z"`
	IsValid   bool                `bson:"is_valid" json:"isValid" binding:"required" ts_type:"boolean"  example:"true"`
}

func (item *BusyTimeStrictRequest) ToModel(photographerID primitive.ObjectID) *models.BusyTime {
	return &models.BusyTime{
		ID:             primitive.NewObjectID(),
		PhotographerID: photographerID,
		Name:           item.Name,
		Type:           item.Type,
		StartTime:      item.StartTime,
		EndTime:        item.EndTime,
		IsValid:        item.IsValid,
	}
}
