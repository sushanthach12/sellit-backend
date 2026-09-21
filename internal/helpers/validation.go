package helpers

import "strings"

func CheckIfStringEmpty(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	return false
}

func CheckStringLen(s string, minLen int64, maxLen int64) bool {
	length := int64(len(strings.TrimSpace(s)))
	return length >= minLen && length <= maxLen
}

func CheckIfValidNumber(num float32, required bool) bool {
	if required && num <= 0 {
		return false
	}

	return true
}
