package apperrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, message string) {
	var statusCode int

	switch err {
	case ErrBadRequest,
		ErrAppointmentStatusTime,
		ErrAppointmentStatusInvalid,
		ErrTimeOverlapped:
		statusCode = http.StatusBadRequest
	case ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case ErrForbidden:
		statusCode = http.StatusForbidden
	default:
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, gin.H{"error": message + ", " + err.Error()})
}
