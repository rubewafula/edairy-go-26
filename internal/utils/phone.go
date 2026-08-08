package utils

import (
	"fmt"
	"strings"
)

func NormalizeKenyaPhone(override, fallback string) (string, error) {
	phone := strings.TrimSpace(override)
	if phone == "" {
		phone = strings.TrimSpace(fallback)
	}
	if phone == "" {
		return "", fmt.Errorf("phone_number is required")
	}
	phone = strings.TrimPrefix(phone, "+")
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	switch {
	case strings.HasPrefix(digits, "254") && len(digits) == 12:
		// ok
	case strings.HasPrefix(digits, "0") && len(digits) == 10:
		digits = "254" + digits[1:]
	case strings.HasPrefix(digits, "7") && len(digits) == 9:
		digits = "254" + digits
	case strings.HasPrefix(digits, "1") && len(digits) == 9:
		digits = "254" + digits
	default:
		return "", fmt.Errorf("invalid phone number format (expected 07XXXXXXXX or 2547XXXXXXXX)")
	}
	if !strings.HasPrefix(digits, "254") || len(digits) != 12 {
		return "", fmt.Errorf("invalid normalized phone number: %s", digits)
	}
	return digits, nil
}
