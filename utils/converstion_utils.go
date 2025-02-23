package utils

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// StructToBsonMap converts a struct to a BSON update map, removing nil & empty fields.
func StructToBsonMap(req interface{}) (bson.M, error) {
	updates := bson.M{}

	// Handle pointer case by getting the underlying struct
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure the input is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	typeData := val.Type()

	// Iterate over struct fields (starting from index 0 to include all fields)
	for i := 0; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		fieldVal := val.Field(i)
		tag := field.Tag.Get("bson") // Get JSON struct tag
		tag = parseBsonTag(tag)      // Extract the actual field name from the tag

		// Skip if the tag is empty or explicitly ignored
		if tag == "" || tag == "-" {
			continue
		}

		// Handle pointer fields properly
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				continue // Skip nil pointers
			}
			fieldVal = fieldVal.Elem() // Dereference the pointer
		}

		// Check for zero values (nil, empty strings, empty slices, empty maps)
		if isZeroValue(fieldVal) {
			continue
		}

		// Add to updates if it's not a zero value
		updates[tag] = fieldVal.Interface()
	}

	return updates, nil
}

// parseBsonTag extracts the actual field name from a BSON tag
func parseBsonTag(tag string) string {
	if tag == "" || tag == "-" {
		return ""
	}
	parts := strings.Split(tag, ",")
	return parts[0] // Extract only the field name before the comma
}

// isZeroValue checks if a field is a zero value (nil, empty string, empty slice, empty map)
func isZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil() // Ensure we handle nil pointers properly
	default:
		return v.IsZero()
	}
}
