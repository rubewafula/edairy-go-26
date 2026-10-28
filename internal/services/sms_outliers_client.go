package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/config"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

// DispatchOutliersV2 sends SMS using Outliers API v2 payload shape.
func DispatchOutliersV2(provider *models.SMSProvider, recipient, message, reference string) error {
	if provider == nil {
		return fmt.Errorf("sms provider is nil")
	}
	apiURL := provider.ApiUrl
	if apiURL == "" {
		apiURL = config.DefaultOutliersSMSAPIURL
	}
	payload := map[string]interface{}{
		"client_code": provider.ProviderCode,
		"key":         provider.ApiKey,
		"short_code":  provider.SenderID,
		"message":     message,
		"recipients":  []string{recipient},
		"reference":   reference,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal sms payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("create sms request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("provider connection error: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sms provider HTTP %d: %s", resp.StatusCode, string(raw))
	}
	return nil
}
