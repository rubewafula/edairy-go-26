package services

import (
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

// parseMemberDate accepts 2006-01-02, RFC3339 (2006-01-02T00:00:00Z) and the other common formats.
// Returns the zero time for empty, "0000-00-00..." or unparseable input.
func parseMemberDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "0000") || strings.EqualFold(s, "null") {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	if len(s) >= 10 {
		if t, err := time.Parse("2006-01-02", s[:10]); err == nil {
			return t
		}
	}
	return utils.ParseFlexibleDate(s)
}

// memberZeroDateColumns lists date columns that hold no date, so Save leaves them as they are
// instead of writing '0000-00-00' (rejected by MySQL strict mode).
func memberZeroDateColumns(m *models.Member) []string {
	cols := []string{}
	if m.DateOfBirth.IsZero() {
		cols = append(cols, "dob")
	}
	if m.DateRegistered.IsZero() {
		cols = append(cols, "date_registered")
	}
	return cols
}

// normalizeMemberStatusValue maps a requested status to the stored value ("" = not a valid status).
func normalizeMemberStatusValue(s string) string {
	switch v := strings.ToLower(strings.TrimSpace(s)); v {
	case "pending", "active", "inactive", "suspended", "closed":
		return v
	}
	return ""
}
