// Package utils provides utility functions for the AWX package.
package utils

import "strconv"

// AtoiDefault takes a string and a defaultValue. If the string cannot be converted, defaultValue is returned.
func AtoiDefault(s string, defaultValue *int) *int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return &n
}
