package config

import (
	"os"
	"strconv"
	"strings"
)

func resolvedValue(values map[string]string, key string) string {
	if values == nil {
		return ""
	}
	return strings.TrimSpace(values[key])
}

func verifySSLFromResolved(values map[string]string) bool {
	if v := resolvedValue(values, "verify_ssl"); v != "" {
		return !(strings.EqualFold(v, "false") || v == "0")
	}
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("VERIFY_SSL"))); v == "false" || v == "0" {
		return false
	}
	return true
}

func requestTimeoutFromResolved(values map[string]string, envKey string, defaultSec int) int {
	if v := resolvedValue(values, "request_timeout_sec"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if v, err := strconv.Atoi(os.Getenv(envKey)); err == nil && v > 0 {
		return v
	}
	return defaultSec
}

// BuildMpesaConfig builds M-Pesa config from resolved channel values (DB + env merged).
func BuildMpesaConfig(values map[string]string) MpesaConfig {
	env := strings.ToLower(resolvedValue(values, "env"))
	if env == "" {
		env = "sandbox"
	}
	cmd := resolvedValue(values, "b2c_command_id")
	if cmd == "" {
		cmd = "BusinessPayment"
	}
	resultURL := resolvedValue(values, "b2c_result_url")
	if resultURL == "" {
		resultURL = "https://dev.edairy.africa/api/mpesa-b2c-result"
	}
	timeoutURL := resolvedValue(values, "b2c_timeout_url")
	if timeoutURL == "" {
		timeoutURL = "https://dev.edairy.africa/api/mpesa-b2c-timeout"
	}
	oauthURL := resolvedValue(values, "oauth_url")
	b2cURL := resolvedValue(values, "b2c_url")
	statusURL := resolvedValue(values, "status_query_url")
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
		ConsumerKey:       resolvedValue(values, "consumer_key"),
		ConsumerSecret:    resolvedValue(values, "consumer_secret"),
		B2CShortcode:      resolvedValue(values, "b2c_shortcode"),
		InitiatorName:     resolvedValue(values, "initiator_name"),
		InitiatorPassword: resolvedValue(values, "initiator_password"),
		CertPath:          resolvedValue(values, "cert_path"),
		B2CCommandID:      cmd,
		B2CResultURL:      resultURL,
		B2CTimeoutURL:     timeoutURL,
		OAuthURL:          oauthURL,
		B2CURL:            b2cURL,
		StatusQueryURL:    statusURL,
		VerifySSL:         verifySSLFromResolved(values),
		RequestTimeoutSec: requestTimeoutFromResolved(values, "MPESA_REQUEST_TIMEOUT_SEC", 60),
	}
}

// BuildAirtelConfig builds Airtel config from resolved channel values.
func BuildAirtelConfig(values map[string]string) AirtelConfig {
	env := strings.ToLower(resolvedValue(values, "env"))
	if env == "" {
		env = "sandbox"
	}
	country := resolvedValue(values, "country")
	if country == "" {
		country = "KE"
	}
	currency := resolvedValue(values, "currency")
	if currency == "" {
		currency = "KES"
	}
	callback := resolvedValue(values, "callback_url")
	if callback == "" {
		callback = "https://dev.edairy.africa/api/airtel-disbursement-callback"
	}
	baseURL := resolvedValue(values, "base_url")
	if baseURL == "" {
		if env == "production" {
			baseURL = "https://openapi.airtel.africa"
		} else {
			baseURL = "https://openapiuat.airtel.africa"
		}
	}

	return AirtelConfig{
		Env:               env,
		ClientID:          resolvedValue(values, "client_id"),
		ClientSecret:      resolvedValue(values, "client_secret"),
		Country:           country,
		Currency:          currency,
		DisbursementPIN:   resolvedValue(values, "disbursement_pin"),
		PublicKeyPath:     resolvedValue(values, "public_key_path"),
		BaseURL:           strings.TrimRight(baseURL, "/"),
		CallbackURL:       callback,
		VerifySSL:         verifySSLFromResolved(values),
		RequestTimeoutSec: requestTimeoutFromResolved(values, "AIRTEL_REQUEST_TIMEOUT_SEC", 60),
	}
}

// BuildJengaConfig builds Jenga/Equitel config from resolved channel values.
func BuildJengaConfig(values map[string]string) JengaConfig {
	env := strings.ToLower(resolvedValue(values, "env"))
	if env == "" {
		env = "sandbox"
	}
	callback := resolvedValue(values, "callback_url")
	if callback == "" {
		callback = "https://dev.edairy.africa/api/jenga-mobile-callback"
	}
	baseURL := resolvedValue(values, "base_url")
	if baseURL == "" {
		if env == "production" {
			baseURL = "https://api.finserve.africa"
		} else {
			baseURL = "https://uat.finserve.africa"
		}
	}
	country := resolvedValue(values, "source_country_code")
	if country == "" {
		country = "KE"
	}

	return JengaConfig{
		Env:               env,
		APIKey:            resolvedValue(values, "api_key"),
		MerchantCode:      resolvedValue(values, "merchant_code"),
		ConsumerSecret:    resolvedValue(values, "consumer_secret"),
		PrivateKeyPath:    resolvedValue(values, "private_key_path"),
		SourceAccount:     resolvedValue(values, "source_account"),
		SourceAccountName: resolvedValue(values, "source_account_name"),
		SourceCountryCode: country,
		CallbackURL:       callback,
		BaseURL:           strings.TrimRight(baseURL, "/"),
		VerifySSL:         verifySSLFromResolved(values),
		RequestTimeoutSec: requestTimeoutFromResolved(values, "JENGA_REQUEST_TIMEOUT_SEC", 60),
	}
}

// BuildAstraConfig builds DTB/Astra config from resolved channel values.
func BuildAstraConfig(values map[string]string) AstraConfig {
	callback := resolvedValue(values, "withdrawal_callback_url")
	if callback == "" {
		callback = "https://dev.edairy.africa/api/withdrawal-callback"
	}
	connectTimeout := 10
	if v, err := strconv.Atoi(os.Getenv("ASTRA_CONNECT_TIMEOUT_SEC")); err == nil && v > 0 {
		connectTimeout = v
	}

	return AstraConfig{
		RemoteURL:          strings.TrimRight(resolvedValue(values, "remote_url"), "/"),
		AuthURL:            strings.TrimRight(resolvedValue(values, "auth_url"), "/"),
		MoneyTransferURL:   strings.TrimRight(resolvedValue(values, "money_transfer_url"), "/"),
		AuthIdentity:       resolvedValue(values, "auth_identity"),
		TenantPassword:     resolvedValue(values, "tenant_password"),
		VerifySSL:          verifySSLFromResolved(values),
		WithdrawalCallback: callback,
		ConnectTimeoutSec:  connectTimeout,
		RequestTimeoutSec:  requestTimeoutFromResolved(values, "ASTRA_REQUEST_TIMEOUT_SEC", 60),
	}
}
