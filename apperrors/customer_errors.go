package apperrors

import "errors"

var (
	ErrBadRequest     = errors.New("Invalid request data")
	ErrInternalServer = errors.New("Internal server error")
	ErrUnauthorized   = errors.New("Unauthorized")

	// appointment
	ErrAppointmentStatusInvalid = errors.New("Invalid appointment status to update")
	ErrAppointmentStatusTime    = errors.New("Invalid status time to update")
)
