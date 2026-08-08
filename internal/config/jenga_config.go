package config

import (
	"os"
	"strconv"
	"strings"
)

type JengaConfig struct {
	Env               string
	APIKey            string
	MerchantCode      string
	ConsumerSecret    string
	PrivateKeyPath    string
	SourceAccount     string
	SourceAccountName string
	SourceCountryCode string
	CallbackURL       string
	BaseURL           string
	VerifySSL         bool
	RequestTimeoutSec int
}

func LoadJengaConfig() JengaConfig {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("JENGA_ENV")))
	if env == "" {
		env = "sandbox"
	}
	verifySSL := true
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("VERIFY_SSL"))); v == "false" || v == "0" {
		verifySSL = false
	}
	timeout := 60
	if v, err := strconv.Atoi(os.Getenv("JENGA_REQUEST_TIMEOUT_SEC")); err == nil && v > 0 {
		timeout = v
	}
	callback := strings.TrimSpace(os.Getenv("JENGA_CALLBACK_URL"))
	if callback == "" {
		callback = "https://dev.edairy.africa/api/jenga-mobile-callback"
	}
	baseURL := strings.TrimSpace(os.Getenv("JENGA_BASE_URL"))
	if baseURL == "" {
		if env == "production" {
			baseURL = "https://api.finserve.africa"
		} else {
			baseURL = "https://uat.finserve.africa"
		}
	}
	country := strings.TrimSpace(os.Getenv("JENGA_SOURCE_COUNTRY_CODE"))
	if country == "" {
		country = "KE"
	}

	return JengaConfig{
		Env:               env,
		APIKey:            os.Getenv("JENGA_API_KEY"),
		MerchantCode:      os.Getenv("JENGA_MERCHANT_CODE"),
		ConsumerSecret:    os.Getenv("JENGA_CONSUMER_SECRET"),
		PrivateKeyPath:    os.Getenv("JENGA_PRIVATE_KEY_PATH"),
		SourceAccount:     os.Getenv("JENGA_SOURCE_ACCOUNT"),
		SourceAccountName: os.Getenv("JENGA_SOURCE_ACCOUNT_NAME"),
		SourceCountryCode: country,
		CallbackURL:       callback,
		BaseURL:           strings.TrimRight(baseURL, "/"),
		VerifySSL:         verifySSL,
		RequestTimeoutSec: timeout,
	}
}

func (c JengaConfig) Configured() bool {
	return c.APIKey != "" && c.ConsumerSecret != "" && c.PrivateKeyPath != "" && c.SourceAccount != ""
}
