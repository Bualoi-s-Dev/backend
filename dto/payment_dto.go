package dto

import "github.com/Bualoi-s-Dev/backend/models"

type PaymentResponse struct {
	Payment     models.Payment    `json:"payment"`
	Appointment AppointmentDetail `json:"appointment"`
}

type PaymentURL struct {
	URL string `json:"url" example:"https://stripe.com"`
}
