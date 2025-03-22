package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1234"`
	CustomerID     primitive.ObjectID `bson:"customer_id" json:"customerId" ts_type:"string" example:"656e2b5e3f1a3c4d8b9e1238"`
	PhotographerID primitive.ObjectID `bson:"photographer_id" json:"photographerId" ts_type:"string" example:"656e2b5e3f1a324d8b9e1236"`
	Rating         int                `bson:"rating" json:"rating" ts_type:"number" example:"5"`
	Review         string             `bson:"review,omitempty" json:"review" ts_type:"string" example:"Very good"`
	CreatedTime    time.Time          `bson:"created_time,omitempty" json:"createdTime" ts_type:"Date" example:"2025-02-23T12:00:00Z"`
}
