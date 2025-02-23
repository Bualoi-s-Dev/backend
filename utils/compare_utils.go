package utils

import (
	"fmt"
	"reflect"
)

// CompareStructsExcept compares two struct values and returns an error if any field differs,
// except for the specified excluded fields.
func CompareStructsExcept(expected interface{}, actual interface{}, excludeFields []string) error {
	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)

	if expectedVal.Kind() == reflect.Ptr {
		expectedVal = expectedVal.Elem()
	}
	if actualVal.Kind() == reflect.Ptr {
		actualVal = actualVal.Elem()
	}

	if expectedVal.Kind() != reflect.Struct || actualVal.Kind() != reflect.Struct {
		return fmt.Errorf("both expected and actual must be structs")
	}

	expectedType := expectedVal.Type()
	excludeMap := make(map[string]struct{})
	for _, field := range excludeFields {
		excludeMap[field] = struct{}{}
	}

	for i := 0; i < expectedVal.NumField(); i++ {
		fieldName := expectedType.Field(i).Name

		// Skip excluded fields
		if _, found := excludeMap[fieldName]; found {
			continue
		}

		expectedField := expectedVal.Field(i)
		actualField := actualVal.Field(i)

		// Handle slice comparison separately
		if expectedField.Kind() == reflect.Slice {
			if !compareSlices(expectedField, actualField) {
				return fmt.Errorf("mismatch in field %s: expected %v, got %v", fieldName, expectedField.Interface(), actualField.Interface())
			}
		} else {

			if !reflect.DeepEqual(expectedField.Interface(), actualField.Interface()) {
				return fmt.Errorf("mismatch in field %s: expected %v, got %v", fieldName, expectedField.Interface(), actualField.Interface())
			}
		}
	}
	return nil
}

// Helper function to compare slices
func compareSlices(a, b reflect.Value) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !reflect.DeepEqual(a.Index(i).Interface(), b.Index(i).Interface()) {
			return false
		}
	}
	return true
}
