package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id" ts_type:"string" example:"12345678abcd"`
	AppointmentID primitive.ObjectID  `bson:"appointment_id" json:"appointmentId" ts_type:"string" example:"12345678abcd"`
	Customer      CustomerPayment     `bson:"customer" json:"customer"`
	Photographer  PhotographerPayment `bson:"photographer" json:"photographer"`
}

type CustomerPayment struct {
	Status          PaymentStatus `bson:"status" json:"status" binding:"omitempty,payment_status" example:"Paid"`
	CheckoutID      *string       `bson:"checkout_id" json:"checkoutId" ts_type:"string" example:"12345678abcd"`
	PaymentIntentID *string       `bson:"payment_intent_id" json:"paymentIntentId" ts_type:"string" example:"12345678abcd"`
}

type PhotographerPayment struct {
	Status               PaymentStatus `bson:"status" json:"status" binding:"omitempty,payment_status" example:"Paid"`
	BalanceTransactionID *string       `bson:"balance_transaction_id" json:"balanceTransactionId" ts_type:"string" example:"12345678abcd"`
}

type PaymentStatus string

const (
	Unpaid    PaymentStatus = "Unpaid"
	InTransit PaymentStatus = "InTransit"
	Paid      PaymentStatus = "Paid"
)
