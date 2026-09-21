package helpers

import "strings"

func CheckIfStringEmpty(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	return false
}

func CheckIfValidNumber(num float32, required bool) bool {
	if required && num == 0 {
		return false
	}

	return true
}
