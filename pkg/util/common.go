package util

import (
	"fmt"
	"math/rand"
	"strings"
)

// ExistsInSlide check if item exists in slice
func ExistsInSlide(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// UniqueSlide remove duplicate items in slice
func RemoveDuplicatetInSlide(input []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueSlide := []string{}

	for _, item := range input {
		// Check if the item is unique
		if !uniqueMap[item] {
			uniqueSlide = append(uniqueSlide, item)
			uniqueMap[item] = true
		}
	}

	return uniqueSlide
}

func IsDuplicateItemInSlice[T comparable](slice []T) bool {
	uniqueMap := make(map[T]bool)

	for _, item := range slice {
		if uniqueMap[item] {
			return true
		}
		uniqueMap[item] = true
	}

	return false
}

// IsValidEmail
func IsValidEmail(email string) bool {
	// Check if the email has the correct format
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	return true
}

// randome from range min to max
func RandomInt(min, max int) int {
	// return min + rand.Intn(max-min)
	if max <= min {
		fmt.Println("DEBUG: Invalid range in RandomInt:", min, max)
		return min // hoặc một giá trị mặc định
	}
	return min + rand.Intn(max-min)
}
