package models

import (
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

// RegisterCustomValidators registers custom validators to the validator
// Add more validators here
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("package_type", ValidatePackageType)
	v.RegisterValidation("bank_name", ValidateBankName)
}
