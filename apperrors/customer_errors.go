package apperrors

import "errors"

var (
	ErrBadRequest     = errors.New("Invalid request data")
	ErrInternalServer = errors.New("Internal server error")
	ErrUnauthorized   = errors.New("Unauthorized")
	ErrForbidden      = errors.New("Permission Denied")
)

// Appointment
var (
	ErrAppointmentStatusInvalid = errors.New("Invalid appointment status to update, Cannot Update Canceled or Completed")
	ErrAppointmentStatusTime    = errors.New("Invalid status time to update")
)

// BusyTime
var (
	ErrTimeOverlapped = errors.New("Time overlap while reserving")
)

// Rating
var (
	ErrCustomerRatingMismatched		= errors.New("Customer does not own this rating")
	ErrPhotographerRatingMismatched	= errors.New("Photographer does not own this rating")	
)
