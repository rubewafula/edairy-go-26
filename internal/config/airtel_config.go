package config

import (
	"os"
	"strconv"
	"strings"
)

type AirtelConfig struct {
	Env              string
	ClientID         string
	ClientSecret     string
	Country          string
	Currency         string
	DisbursementPIN  string
	PublicKeyPath    string
	BaseURL          string
	CallbackURL      string
	VerifySSL        bool
	RequestTimeoutSec int
}

func LoadAirtelConfig() AirtelConfig {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("AIRTEL_ENV")))
	if env == "" {
		env = "sandbox"
	}
	country := strings.TrimSpace(os.Getenv("AIRTEL_COUNTRY"))
	if country == "" {
		country = "KE"
	}
	currency := strings.TrimSpace(os.Getenv("AIRTEL_CURRENCY"))
	if currency == "" {
		currency = "KES"
	}
	verifySSL := true
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("VERIFY_SSL"))); v == "false" || v == "0" {
		verifySSL = false
	}
	timeout := 60
	if v, err := strconv.Atoi(os.Getenv("AIRTEL_REQUEST_TIMEOUT_SEC")); err == nil && v > 0 {
		timeout = v
	}
	callback := strings.TrimSpace(os.Getenv("AIRTEL_CALLBACK_URL"))
	if callback == "" {
		callback = "https://dev.edairy.africa/api/airtel-disbursement-callback"
	}
	baseURL := strings.TrimSpace(os.Getenv("AIRTEL_BASE_URL"))
	if baseURL == "" {
		if env == "production" {
			baseURL = "https://openapi.airtel.africa"
		} else {
			baseURL = "https://openapiuat.airtel.africa"
		}
	}

	return AirtelConfig{
		Env:               env,
		ClientID:          os.Getenv("AIRTEL_CLIENT_ID"),
		ClientSecret:      os.Getenv("AIRTEL_CLIENT_SECRET"),
		Country:           country,
		Currency:          currency,
		DisbursementPIN:   os.Getenv("AIRTEL_DISBURSEMENT_PIN"),
		PublicKeyPath:     os.Getenv("AIRTEL_PUBLIC_KEY_PATH"),
		BaseURL:           strings.TrimRight(baseURL, "/"),
		CallbackURL:       callback,
		VerifySSL:         verifySSL,
		RequestTimeoutSec: timeout,
	}
}

func (c AirtelConfig) Configured() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.DisbursementPIN != "" && c.PublicKeyPath != ""
}
