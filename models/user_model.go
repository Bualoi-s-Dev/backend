package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string `bson:"email" json:"email"`
	Name     string `bson:"name,omitempty" json:"name"`
	Gender   string `bson:"gender,omitempty" json:"gender"`
	Profile  string `bson:"profile,omitempty" json:"profile"`
	Phone    string `bson:"phone,omitempty" json:"phone"`
	Location string `bson:"location,omitempty" json:"location"`
}