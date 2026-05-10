package utils

import (
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
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
	t, _ := time.Parse("2006-01-02", d)
	return t
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
