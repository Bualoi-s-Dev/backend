package dto

import (
	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingRequest struct {
	Rating	int		`bson:"rating" json:"rating" binding:"required" ts_type:"number" example:"5"`
	Review	*string	`bson:"review,omitempty" json:"review" binding:"omitempty" ts_type:"string" example:"Very good"`
}

type RatingResponse struct {
	ID				primitive.ObjectID	`bson:"_id,omitempty" json:"id" binding:"required" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID		primitive.ObjectID	`bson:"customer_id" json:"customerId" binding:"required" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1238"`
	PhotographerID	primitive.ObjectID	`bson:"photographer_id" json:"photographerId" binding:"required" ts_type:"string" example:"656e2b5e3f1a324d8b9e1236"`
	Rating			int					`bson:"rating" json:"rating" binding:"required" ts_type:"number" example:"5"`
	Review			string				`bson:"review,omitempty" json:"review" binding:"omitempty" ts_type:"string" example:"Very good"`
}

func (item *RatingRequest) ToModel(customerId primitive.ObjectID, photographerId primitive.ObjectID) *models.Rating {
	return &models.Rating{
		CustomerID:		customerId,
		PhotographerID: photographerId,
		Rating:			item.Rating,
		Review:			*item.Review,
	}
}