package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// ValidatePackageType checks if the PackageType is valid
func ValidatePackageType(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(PackageType)

	// Check if the value exists in the validPackageTypes slice
	for _, validType := range ValidPackageTypes {
		if value == validType.Value {
			return true
		}
	}

	return false
}

// ValidateBankname checks if the BankName is valid
func ValidateBankName(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(BankName)

	// Check if the value exists in the validBankNames slice
	for _, validBank := range ValidBankNames {
		if value == validBank.Value {
			return true
		}
	}

	return false
}

// ValidateBankname checks if the UserRole is valid
func ValidateUserRole(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(UserRole)

	// Check if the value exists in the validBankNames slice
	for _, validRole := range ValidUserRoles {
		if value == validRole.Value {
			return true
		}
	}

	return false
}

func ValidateDayNames(fl validator.FieldLevel) bool {
	field := fl.Field().Interface().([]DayName)

	for _, day := range field {
		isValid := false
		for _, validDay := range ValidDayNames {
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
	value := fl.Field().Interface().(string)
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

// ValidateBusyTimeType checks if the BusyTimeType is valid
func ValidateBusyTimeType(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(BusyTimeType)

	// Check if the value exists in the validBankNames slice
	for _, validType := range ValidBusyTimeTypes {
		if value == validType.Value {
			return true
		}
	}

	return false
}

// RegisterCustomValidators registers custom validators to the validator
// Add more validators here
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("package_type", ValidatePackageType)
	v.RegisterValidation("bank_name", ValidateBankName)
	v.RegisterValidation("user_role", ValidateUserRole)
	v.RegisterValidation("day_names", ValidateDayNames)
	v.RegisterValidation("time_format", ValidateTime)
	v.RegisterValidation("date_format", ValidateDate)
	v.RegisterValidation("busy_time_type", ValidateBusyTimeType)
}
