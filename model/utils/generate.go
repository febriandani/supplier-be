package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

func GenerateTraceID() string {
	return uuid.New().String()
}

func FormatPhoneNumber(phone string) string {
	// Remove any leading spaces or plus sign
	phone = strings.TrimLeft(phone, " +")

	// Check if the phone number starts with "62"
	if strings.HasPrefix(phone, "62") {
		return phone
	}

	// Check if the phone number starts with "+62"
	if strings.HasPrefix(phone, "+62") {
		// Remove the "+" character and return the number
		return strings.TrimPrefix(phone, "+")
	}

	if strings.Contains(phone, "0") {
		// Replace the "0" with "62" as the prefix and return the number
		return fmt.Sprintf("62%s", strings.Replace(phone, "0", "", 1))
	}

	// If the phone number doesn't start with "62" or "+62", add "62" as the prefix
	return fmt.Sprintf("62%s", phone)
}

// IsValidEmail checks if the given email is in a valid format
func IsValidEmail(email string) bool {
	// Basic regex for email validation
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsPhoneNumber(input string) bool {
	// We'll assume phone numbers are digits only (or start with +)
	cleaned := strings.TrimSpace(input)
	cleaned = strings.TrimPrefix(cleaned, "+")

	for _, r := range cleaned {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(cleaned) >= 9 // could be flexible
}
