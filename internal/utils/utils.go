package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NormalizePhone(phone string) string {
	re := regexp.MustCompile(`^(?:0|254)?([17]\d{8})$`)
	m := re.FindStringSubmatch(phone)

	if len(m) == 2 {
		return "254" + m[1]
	}
	return phone
}

func ParseDate(d string) time.Time {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04", // Matches "2026-05-22T03:52"
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, d); err == nil {
			return t
		}
	}
	return time.Time{}
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = e.Tag()
	}

	return errors
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Uint64ToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}

// GormDeletedAtToPtr converts gorm.DeletedAt to *time.Time for JSON omitempty
func GormDeletedAtToPtr(deletedAt gorm.DeletedAt) *time.Time {
	if deletedAt.Valid {
		return &deletedAt.Time
	}
	return nil
}

// GenerateUniqueFilename creates a unique filename using a UUID and the original extension.
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return uuid.New().String() + ext
}

// SaveBase64ToFile decodes a base64 string, saves it to disk, and returns the relative path.
// The uploadDir should be relative to the application's root.
func SaveBase64ToFile(base64Content string, originalFilename string, uploadDir string) (string, error) {
	if base64Content == "" {
		return "", nil
	}

	// Decode base64 string
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %w", err)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate a unique filename
	uniqueFilename := GenerateUniqueFilename(originalFilename)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return "/" + filePath, nil // Return the relative URL/path
}

// DeleteFile deletes a file from the file system given its path.
func DeleteFile(filePath string) error {
	// Remove leading slash if present, as filepath.Join expects relative paths
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File does not exist, nothing to delete
	}
	return os.Remove(filePath)
}
