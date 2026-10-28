package config

import "os"

const DefaultOutliersSMSAPIURL = "https://sms.outliers.co.ke/api/v2/send-sms"

// OutliersSMSConfig maps to Outliers POST /api/v2/send-sms JSON body fields.
type OutliersSMSConfig struct {
	ClientCode string // client_code
	APIKey     string // key
	SenderID   string // short_code (e.g. E-Dairy)
	APIURL     string
}

func LoadOutliersSMSConfig() OutliersSMSConfig {
	url := os.Getenv("SMS_OUTLIERS_API_URL")
	if url == "" {
		url = DefaultOutliersSMSAPIURL
	}
	return OutliersSMSConfig{
		ClientCode: os.Getenv("SMS_OUTLIERS_CLIENT_CODE"),
		APIKey:     os.Getenv("SMS_OUTLIERS_API_KEY"),
		SenderID:   os.Getenv("SMS_OUTLIERS_SENDER_ID"),
		APIURL:     url,
	}
}

func (c OutliersSMSConfig) Configured() bool {
	return c.ClientCode != "" && c.APIKey != "" && c.SenderID != ""
}
