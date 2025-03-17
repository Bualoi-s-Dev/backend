package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingRequest struct {
	Rating	int		`bson:"rating" json:"rating" binding:"required" ts_type:"number" example:"5"`
	Review	*string	`bson:"review,omitempty" json:"review" binding:"omitempty" ts_type:"string" example:"Very good"`
}

func (item *RatingRequest) ToModel(customerId primitive.ObjectID, photographerId primitive.ObjectID) *models.Rating {
	return &models.Rating{
		CustomerID:		customerId,
		PhotographerID: photographerId,
		Rating:			item.Rating,
		Review:			*item.Review,
	}
}