package config

import (
	"os"
	"strconv"
	"strings"
)

type MpesaConfig struct {
	Env              string
	ConsumerKey      string
	ConsumerSecret   string
	B2CShortcode     string
	InitiatorName    string
	InitiatorPassword string
	CertPath         string
	B2CCommandID     string
	B2CResultURL     string
	B2CTimeoutURL    string
	OAuthURL         string
	B2CURL           string
	StatusQueryURL   string
	VerifySSL        bool
	RequestTimeoutSec int
}

func LoadMpesaConfig() MpesaConfig {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("MPESA_ENV")))
	if env == "" {
		env = "sandbox"
	}
	verifySSL := true
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("VERIFY_SSL"))); v == "false" || v == "0" {
		verifySSL = false
	}
	timeout := 60
	if v, err := strconv.Atoi(os.Getenv("MPESA_REQUEST_TIMEOUT_SEC")); err == nil && v > 0 {
		timeout = v
	}
	cmd := strings.TrimSpace(os.Getenv("MPESA_B2C_COMMAND_ID"))
	if cmd == "" {
		cmd = "BusinessPayment"
	}
	resultURL := strings.TrimSpace(os.Getenv("MPESA_B2C_RESULT_URL"))
	if resultURL == "" {
		resultURL = "https://dev.edairy.africa/api/mpesa-b2c-result"
	}
	timeoutURL := strings.TrimSpace(os.Getenv("MPESA_B2C_TIMEOUT_URL"))
	if timeoutURL == "" {
		timeoutURL = "https://dev.edairy.africa/api/mpesa-b2c-timeout"
	}

	oauthURL := strings.TrimSpace(os.Getenv("MPESA_OAUTH_URL"))
	b2cURL := strings.TrimSpace(os.Getenv("MPESA_B2C_URL"))
	statusURL := strings.TrimSpace(os.Getenv("MPESA_STATUS_QUERY_URL"))
	if env == "production" {
		if oauthURL == "" {
			oauthURL = "https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
		}
		if b2cURL == "" {
			b2cURL = "https://api.safaricom.co.ke/mpesa/b2c/v3/paymentrequest"
		}
		if statusURL == "" {
			statusURL = "https://api.safaricom.co.ke/mpesa/transactionstatus/v1/query"
		}
	} else {
		if oauthURL == "" {
			oauthURL = "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
		}
		if b2cURL == "" {
			b2cURL = "https://sandbox.safaricom.co.ke/mpesa/b2c/v3/paymentrequest"
		}
		if statusURL == "" {
			statusURL = "https://sandbox.safaricom.co.ke/mpesa/transactionstatus/v1/query"
		}
	}

	return MpesaConfig{
		Env:               env,
		ConsumerKey:       os.Getenv("MPESA_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("MPESA_CONSUMER_SECRET"),
		B2CShortcode:      os.Getenv("MPESA_B2C_SHORTCODE"),
		InitiatorName:     os.Getenv("MPESA_INITIATOR_NAME"),
		InitiatorPassword: os.Getenv("MPESA_INITIATOR_PASSWORD"),
		CertPath:          os.Getenv("MPESA_CERT_PATH"),
		B2CCommandID:      cmd,
		B2CResultURL:      resultURL,
		B2CTimeoutURL:     timeoutURL,
		OAuthURL:          oauthURL,
		B2CURL:            b2cURL,
		StatusQueryURL:    statusURL,
		VerifySSL:         verifySSL,
		RequestTimeoutSec: timeout,
	}
}

func (c MpesaConfig) Configured() bool {
	return c.ConsumerKey != "" && c.ConsumerSecret != "" && c.B2CShortcode != "" &&
		c.InitiatorName != "" && c.InitiatorPassword != "" && c.CertPath != ""
}
