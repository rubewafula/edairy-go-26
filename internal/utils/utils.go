package utils

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GetUniqueRouteIDs extracts a slice of unique non-zero route IDs from a transporter-to-route map.
func GetUniqueRouteIDs(m map[uint64]uint64) []uint64 {
	uniqueMap := make(map[uint64]struct{})
	for _, routeID := range m {
		if routeID != 0 {
			uniqueMap[routeID] = struct{}{}
		}
	}

	result := make([]uint64, 0, len(uniqueMap))
	for routeID := range uniqueMap {
		result = append(result, routeID)
	}
	return result
}

// ParseDate parses a string into a time.Time object.
func ParseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

// ParseFlexibleDate attempts to parse a date string using various common formats.
func ParseFlexibleDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	formats := []string{
		"02 Jan 2006", // 01 Jan 2024
		"2006-01-02",  // ISO
		"02/01/2006",  // DD/MM/YYYY
		"01/02/2006",  // MM/DD/YYYY
		"02-01-2006",  // DD-MM-YYYY
		"2006/01/02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return t
		}
	}
	return time.Time{}
}

// ParseDatePtr parses a string into a time.Time pointer, returning nil if empty or invalid.
func ParseDatePtr(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

// Uint64ToString converts a uint64 to a string.
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// FormatValidationError returns a simple string representation of validation errors.
func FormatValidationError(err error) string {
	return err.Error()
}

// HashPassword hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// NormalizePhone removes non-digit characters and ensures a common prefix (e.g., +254).
func NormalizePhone(phone string) string {
	digitsOnly := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
	if strings.HasPrefix(digitsOnly, "0") {
		return "+254" + digitsOnly[1:]
	}
	return "+" + digitsOnly
}

// ParseFloat safely converts a string to float64, returning 0 if empty or invalid.
func ParseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}
