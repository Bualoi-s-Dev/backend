package validators

import (
	"time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/go-playground/validator/v10"
)

// ValidatePackageType checks if the PackageType is valid
func ValidatePackageType(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(models.PackageType)

	// Check if the value exists in the validPackageTypes slice
	for _, validType := range models.ValidPackageTypes {
		if value == validType.Value {
			return true
		}
	}

	return false
}

// ValidateBankname checks if the BankName is valid
func ValidateBankName(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(models.BankName)

	// Check if the value exists in the validBankNames slice
	for _, validBank := range models.ValidBankNames {
		if value == validBank.Value {
			return true
		}
	}

	return false
}

// ValidateBankname checks if the UserRole is valid
func ValidateUserRole(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(models.UserRole)

	// Check if the value exists in the validBankNames slice
	for _, validRole := range models.ValidUserRoles {
		if value == validRole.Value {
			return true
		}
	}

	return false
}

// ValidateAppointmentStatus check if AppointmentStatus is valid
func ValidateAppointmentStatus(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(models.AppointmentStatus)

	// Check if the value exists in the validBankNames slice
	for _, validStatus := range models.ValidAppointmentStatus {
		if value == validStatus.Value {
			return true
		}
	}
	return false
}

func ValidateDayNames(fl validator.FieldLevel) bool {
	field := fl.Field().Interface().([]models.DayName)

	for _, day := range field {
		isValid := false
		for _, validDay := range models.ValidDayNames {
			if day == validDay.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			return false
		}
	}
	return true
}

// Validate time format
func ValidateTime(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	_, err := time.Parse("15:04", value)
	return err == nil
}

// Validate date format
func ValidateDate(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok || value == "" {
		return true // Allow empty values (handled separately by isInf rule)
	}
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

// ValidateBusyTimeType checks if the BusyTimeType is valid
func ValidateBusyTimeType(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(models.BusyTimeType)

	// Check if the value exists in the validBankNames slice
	for _, validType := range models.ValidBusyTimeTypes {
		if value == validType.Value {
			return true
		}
	}

	return false
}

// Custom validation function
func IsInfRule(fl validator.FieldLevel) bool {
	req, ok := fl.Parent().Interface().(dto.SubpackageRequest)
	if !ok {
		return false
	}

	// If IsInf is true, start & end days can be empty
	if req.IsInf != nil && *req.IsInf {
		return true
	}

	// If IsInf is false, both fields must be non-empty
	if req.IsInf != nil && !*req.IsInf {
		return req.AvailableStartDay != nil && *req.AvailableStartDay != "" &&
			req.AvailableEndDay != nil && *req.AvailableEndDay != ""
	}

	return true
}

// RegisterCustomValidators registers custom validators to the validator
// Add more validators here
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("package_type", ValidatePackageType)
	v.RegisterValidation("bank_name", ValidateBankName)
	v.RegisterValidation("user_role", ValidateUserRole)
	v.RegisterValidation("appointment_status", ValidateAppointmentStatus)
	v.RegisterValidation("day_names", ValidateDayNames)
	v.RegisterValidation("time_format", ValidateTime)
	v.RegisterValidation("date_format", ValidateDate)
	v.RegisterValidation("busy_time_type", ValidateBusyTimeType)
	v.RegisterValidation("isInf_rule", IsInfRule)
}
