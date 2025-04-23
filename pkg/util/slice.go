package util

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IndexByFieldUnique[T any](slice []T, fieldName string) (map[string]T, error) {
	result := make(map[string]T)

	// Reflect on the slice to ensure it contains structs
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("provided value is not a slice")
	}

	elemType := sliceType.Elem()
	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("elements in the slice are not structs")
	}

	for _, item := range slice {
		// Use reflection to access the field
		fieldVal := reflect.ValueOf(item).FieldByName(fieldName)
		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("no such field: %s in struct", fieldName)
		}

		// Check if the field type is ObjectID and call .Hex() method if so
		var key string
		if fieldVal.Type() == reflect.TypeOf(primitive.ObjectID{}) {
			method := fieldVal.MethodByName("Hex")
			if !method.IsValid() {
				return nil, fmt.Errorf("the Hex method is not valid on field: %s", fieldName)
			}
			key = method.Call(nil)[0].String() // Call .Hex() and use the result as the key
		} else {
			key = fmt.Sprintf("%v", fieldVal.Interface()) // Use the string value of the field as the map key
		}

		// Map the key to the current item, replacing any previous item with the same key
		result[key] = item
	}

	return result, nil
}

func BuildSliceWithExcluded[T comparable](slice []T, excluded []T) []T {
	result := make([]T, 0, len(slice))

	// Create a map of excluded items for faster lookup
	excludedMap := make(map[T]struct{}, len(excluded))
	for _, item := range excluded {
		excludedMap[item] = struct{}{}
	}

	// Iterate over the slice and append items that are not in the excluded map
	for _, item := range slice {
		if _, ok := excludedMap[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

func strSlicecontains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
