package config

import "strings"

const (
	ConfigSourceDB    = "db"
	ConfigSourceEnv   = "env"
	ConfigSourceUnset = "unset"
)

type DisbursementConfigKeyDef struct {
	Key       string
	EnvVar    string
	Label     string
	InputType string // text, password, path, url
	IsSecret  bool
}

var providerConfigTemplates = map[string][]DisbursementConfigKeyDef{
	"MPESA": {
		{Key: "env", EnvVar: "MPESA_ENV", Label: "Environment", InputType: "text"},
		{Key: "consumer_key", EnvVar: "MPESA_CONSUMER_KEY", Label: "Consumer Key", InputType: "password", IsSecret: true},
		{Key: "consumer_secret", EnvVar: "MPESA_CONSUMER_SECRET", Label: "Consumer Secret", InputType: "password", IsSecret: true},
		{Key: "b2c_shortcode", EnvVar: "MPESA_B2C_SHORTCODE", Label: "B2C Shortcode", InputType: "text"},
		{Key: "initiator_name", EnvVar: "MPESA_INITIATOR_NAME", Label: "Initiator Name", InputType: "text"},
		{Key: "initiator_password", EnvVar: "MPESA_INITIATOR_PASSWORD", Label: "Initiator Password", InputType: "password", IsSecret: true},
		{Key: "cert_path", EnvVar: "MPESA_CERT_PATH", Label: "Certificate Path", InputType: "path"},
		{Key: "b2c_command_id", EnvVar: "MPESA_B2C_COMMAND_ID", Label: "B2C Command ID", InputType: "text"},
		{Key: "b2c_result_url", EnvVar: "MPESA_B2C_RESULT_URL", Label: "B2C Result URL", InputType: "url"},
		{Key: "b2c_timeout_url", EnvVar: "MPESA_B2C_TIMEOUT_URL", Label: "B2C Timeout URL", InputType: "url"},
		{Key: "oauth_url", EnvVar: "MPESA_OAUTH_URL", Label: "OAuth URL", InputType: "url"},
		{Key: "b2c_url", EnvVar: "MPESA_B2C_URL", Label: "B2C URL", InputType: "url"},
		{Key: "status_query_url", EnvVar: "MPESA_STATUS_QUERY_URL", Label: "Status Query URL", InputType: "url"},
	},
	"AIRTEL": {
		{Key: "env", EnvVar: "AIRTEL_ENV", Label: "Environment", InputType: "text"},
		{Key: "client_id", EnvVar: "AIRTEL_CLIENT_ID", Label: "Client ID", InputType: "password", IsSecret: true},
		{Key: "client_secret", EnvVar: "AIRTEL_CLIENT_SECRET", Label: "Client Secret", InputType: "password", IsSecret: true},
		{Key: "disbursement_pin", EnvVar: "AIRTEL_DISBURSEMENT_PIN", Label: "Disbursement PIN", InputType: "password", IsSecret: true},
		{Key: "public_key_path", EnvVar: "AIRTEL_PUBLIC_KEY_PATH", Label: "Public Key Path", InputType: "path"},
		{Key: "country", EnvVar: "AIRTEL_COUNTRY", Label: "Country Code", InputType: "text"},
		{Key: "currency", EnvVar: "AIRTEL_CURRENCY", Label: "Currency", InputType: "text"},
		{Key: "base_url", EnvVar: "AIRTEL_BASE_URL", Label: "Base URL", InputType: "url"},
		{Key: "callback_url", EnvVar: "AIRTEL_CALLBACK_URL", Label: "Callback URL", InputType: "url"},
	},
	"EQUITY_JENGA": {
		{Key: "env", EnvVar: "JENGA_ENV", Label: "Environment", InputType: "text"},
		{Key: "api_key", EnvVar: "JENGA_API_KEY", Label: "API Key", InputType: "password", IsSecret: true},
		{Key: "merchant_code", EnvVar: "JENGA_MERCHANT_CODE", Label: "Merchant Code", InputType: "text"},
		{Key: "consumer_secret", EnvVar: "JENGA_CONSUMER_SECRET", Label: "Consumer Secret", InputType: "password", IsSecret: true},
		{Key: "private_key_path", EnvVar: "JENGA_PRIVATE_KEY_PATH", Label: "Private Key Path", InputType: "path"},
		{Key: "source_account", EnvVar: "JENGA_SOURCE_ACCOUNT", Label: "Source Account", InputType: "text"},
		{Key: "source_account_name", EnvVar: "JENGA_SOURCE_ACCOUNT_NAME", Label: "Source Account Name", InputType: "text"},
		{Key: "source_country_code", EnvVar: "JENGA_SOURCE_COUNTRY_CODE", Label: "Source Country Code", InputType: "text"},
		{Key: "callback_url", EnvVar: "JENGA_CALLBACK_URL", Label: "Callback URL", InputType: "url"},
		{Key: "base_url", EnvVar: "JENGA_BASE_URL", Label: "Base URL", InputType: "url"},
	},
	"DTB": {
		{Key: "remote_url", EnvVar: "ASTRA_API_BASE_URL", Label: "API Base URL", InputType: "url"},
		{Key: "auth_url", EnvVar: "ASTRA_AUTHENTICATION_API_BASE_URL", Label: "Auth API Base URL", InputType: "url"},
		{Key: "money_transfer_url", EnvVar: "MONEY_TRANSFER_BASE_URL", Label: "Money Transfer Base URL", InputType: "url"},
		{Key: "auth_identity", EnvVar: "AUTH_IDENTITY", Label: "Auth Identity", InputType: "text"},
		{Key: "tenant_password", EnvVar: "ASTRA_TENANT_PASSWORD", Label: "Tenant Password", InputType: "password", IsSecret: true},
		{Key: "withdrawal_callback_url", EnvVar: "WITHDRAWAL_CALLBACK_URL", Label: "Withdrawal Callback URL", InputType: "url"},
	},
}

// Mobile channels always use their API template even if provider label differs (e.g. EQUITY vs EQUITY_JENGA).
var channelProviderTemplate = map[string]string{
	"MOBILE_MPESA":     "MPESA",
	"MOBILE_AIRTEL":    "AIRTEL",
	"MOBILE_EQUITEL":   "EQUITY_JENGA",
	"MOBILE_DTB_MPESA": "DTB",
}

func GetProviderConfigTemplate(provider string) []DisbursementConfigKeyDef {
	p := strings.ToUpper(strings.TrimSpace(provider))
	if defs, ok := providerConfigTemplates[p]; ok {
		out := make([]DisbursementConfigKeyDef, len(defs))
		copy(out, defs)
		return out
	}
	return nil
}

func GetProviderConfigTemplateForChannel(channelCode, provider string) []DisbursementConfigKeyDef {
	code := strings.ToUpper(strings.TrimSpace(channelCode))
	if templateKey, ok := channelProviderTemplate[code]; ok {
		if defs := providerConfigTemplates[templateKey]; len(defs) > 0 {
			out := make([]DisbursementConfigKeyDef, len(defs))
			copy(out, defs)
			return out
		}
	}
	return GetProviderConfigTemplate(provider)
}

func LookupConfigKeyDef(provider, key string) (DisbursementConfigKeyDef, bool) {
	return LookupConfigKeyDefForChannel("", provider, key)
}

func LookupConfigKeyDefForChannel(channelCode, provider, key string) (DisbursementConfigKeyDef, bool) {
	key = strings.ToLower(strings.TrimSpace(key))
	for _, def := range GetProviderConfigTemplateForChannel(channelCode, provider) {
		if def.Key == key {
			return def, true
		}
	}
	return DisbursementConfigKeyDef{}, false
}

func MaskSecretValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(value)-4) + value[len(value)-4:]
}
