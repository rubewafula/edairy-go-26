package config

import (
	"os"
	"strconv"
	"strings"
)

type AstraConfig struct {
	RemoteURL          string
	AuthURL            string
	MoneyTransferURL   string
	AuthIdentity       string
	TenantPassword     string
	VerifySSL          bool
	WithdrawalCallback string
	ConnectTimeoutSec  int
	RequestTimeoutSec  int
}

func LoadAstraConfig() AstraConfig {
	verifySSL := true
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("VERIFY_SSL"))); v == "false" || v == "0" {
		verifySSL = false
	}
	connectTimeout := 10
	if v, err := strconv.Atoi(os.Getenv("ASTRA_CONNECT_TIMEOUT_SEC")); err == nil && v > 0 {
		connectTimeout = v
	}
	requestTimeout := 60
	if v, err := strconv.Atoi(os.Getenv("ASTRA_REQUEST_TIMEOUT_SEC")); err == nil && v > 0 {
		requestTimeout = v
	}
	callback := strings.TrimSpace(os.Getenv("WITHDRAWAL_CALLBACK_URL"))
	if callback == "" {
		callback = "https://dev.edairy.africa/api/withdrawal-callback"
	}
	return AstraConfig{
		RemoteURL:          strings.TrimRight(os.Getenv("ASTRA_API_BASE_URL"), "/"),
		AuthURL:            strings.TrimRight(os.Getenv("ASTRA_AUTHENTICATION_API_BASE_URL"), "/"),
		MoneyTransferURL:   strings.TrimRight(os.Getenv("MONEY_TRANSFER_BASE_URL"), "/"),
		AuthIdentity:       os.Getenv("AUTH_IDENTITY"),
		TenantPassword:     os.Getenv("ASTRA_TENANT_PASSWORD"),
		VerifySSL:          verifySSL,
		WithdrawalCallback: callback,
		ConnectTimeoutSec:  connectTimeout,
		RequestTimeoutSec:  requestTimeout,
	}
}

func (c AstraConfig) Configured() bool {
	return c.RemoteURL != "" && c.AuthURL != "" && c.AuthIdentity != "" && c.TenantPassword != ""
}
