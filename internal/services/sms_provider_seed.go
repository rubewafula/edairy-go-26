package services

import (
	"log"

	"github.com/rubewafula/edairy-go-26/internal/config"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// EnsureOutliersSMSProvider upserts the default Outliers provider when SMS_OUTLIERS_* env is set.
func EnsureOutliersSMSProvider() error {
	cfg := config.LoadOutliersSMSConfig()
	if !cfg.Configured() {
		return nil
	}

	var existing models.SMSProvider
	err := db.DB.Where("provider_code = ?", cfg.ClientCode).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	provider := models.SMSProvider{
		ProviderCode: cfg.ClientCode,
		ProviderName: "Outliers",
		ApiUrl:       cfg.APIURL,
		ApiKey:       cfg.APIKey,
		SenderID:     cfg.SenderID,
		IsDefault:    "yes",
		Status:       "active",
		Notes:        "Synced from SMS_OUTLIERS_* environment",
	}
	if existing.ID != 0 {
		provider.ID = existing.ID
		if err := db.DB.Model(&existing).Updates(map[string]interface{}{
			"provider_name": provider.ProviderName,
			"api_url":       provider.ApiUrl,
			"api_key":       provider.ApiKey,
			"sender_id":     provider.SenderID,
			"is_default":    "yes",
			"status":        "active",
			"notes":         provider.Notes,
		}).Error; err != nil {
			return err
		}
		log.Printf("[SMS] Updated default Outliers provider %s (sender %s)", cfg.ClientCode, cfg.SenderID)
		return nil
	}

	_ = db.DB.Model(&models.SMSProvider{}).Where("is_default = ?", "yes").Update("is_default", "no")
	if err := db.DB.Create(&provider).Error; err != nil {
		return err
	}
	log.Printf("[SMS] Created default Outliers provider %s (sender %s)", cfg.ClientCode, cfg.SenderID)
	return nil
}
