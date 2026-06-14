package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type SMSCampaignService struct {
	notificationService *UINotificationService
}

func NewSMSCampaignService() *SMSCampaignService {
	return &SMSCampaignService{
		notificationService: NewUINotificationService(),
	}
}

func (s *SMSCampaignService) CreateCampaign(req dtos.CreateSMSCampaignRequest, userID uint64) (*models.SMSCampaign, error) {
	var campaign *models.SMSCampaign

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create entry in sms_messages
		smsMsg := &models.SMSMessage{
			Message: req.Message,
			Status:  "pending",
		}
		if err := tx.Create(smsMsg).Error; err != nil {
			return err
		}

		// 2. Create the Campaign header in 'processing' status
		campaign = &models.SMSCampaign{
			BaseModel:       models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			CampaignName:    req.CampaignName,
			CampaignCode:    req.CampaignCode,
			Description:     req.Description,
			SMSMessageID:    &smsMsg.ID,
			TotalRecipients: 0, // Resolved in background
			Status:          "processing",
		}

		if req.ScheduledAt != "" {
			parsedTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
			if err != nil {
				return err
			}
			campaign.ScheduledAt = &parsedTime
		}

		if err := tx.Create(campaign).Error; err != nil {
			return err
		}

		// 3. Link targeted groups to the campaign
		if len(req.SMSGroupIDs) > 0 {
			var groups []models.SMSGroup
			if err := tx.Where("id IN ?", req.SMSGroupIDs).Find(&groups).Error; err != nil {
				return err
			}
			if err := tx.Model(campaign).Association("SMSGroups").Append(groups); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 3. Trigger heavy lifting in background
	go s.processCampaignInBackground(campaign.ID, userID, req.SMSGroupIDs, req.Message)

	return campaign, nil
}

func (s *SMSCampaignService) GetCampaigns(page, limit int) ([]models.SMSCampaign, int64, error) {
	var campaigns []models.SMSCampaign
	var total int64

	db.DB.Model(&models.SMSCampaign{}).Count(&total)

	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&campaigns).Error
	if err != nil {
		return nil, 0, err
	}
	return campaigns, total, nil
}

func (s *SMSCampaignService) GetCampaign(id string) (*models.SMSCampaign, error) {
	var campaign models.SMSCampaign
	if err := db.DB.First(&campaign, id).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (s *SMSCampaignService) UpdateCampaign(id string, req dtos.UpdateSMSCampaignRequest, userID uint64) error {
	var campaign models.SMSCampaign
	if err := db.DB.First(&campaign, id).Error; err != nil {
		return err
	}

	campaign.UpdatedBy = userID
	campaign.CampaignName = req.CampaignName
	campaign.CampaignCode = req.CampaignCode
	campaign.Description = req.Description
	campaign.Status = req.Status

	if req.ScheduledAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			return fmt.Errorf("invalid scheduled_at format: %w", err)
		}
		campaign.ScheduledAt = &parsedTime
	} else {
		campaign.ScheduledAt = nil // Clear if empty string is provided
	}

	return db.DB.Save(&campaign).Error
}

func (s *SMSCampaignService) DeleteCampaign(id string) error {
	var campaign models.SMSCampaign
	if err := db.DB.First(&campaign, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&campaign).Error
}

func (s *SMSCampaignService) processCampaignInBackground(campaignID uint64, userID uint64, groupIDs []uint64, message string) {
	log.Printf("[SMSCampaignService.processCampaignInBackground] Starting background processing for campaign %d", campaignID)

	var campaign models.SMSCampaign
	if err := db.DB.Preload("SMSGroups").First(&campaign, campaignID).Error; err != nil {
		log.Printf("[SMSCampaignService] Background processing failed: campaign %d not found", campaignID)
		return
	}

	// If groupIDs parameter is empty (e.g. on resume), use the persistent associations
	if len(groupIDs) == 0 {
		for _, g := range campaign.SMSGroups {
			groupIDs = append(groupIDs, g.ID)
		}
	}

	// 1. Resolve Contacts from Groups
	var contacts []models.SMSContact
	err := db.DB.Table("sms_contacts").
		Select("sms_contacts.*").
		Joins("JOIN sms_contact_groups ON sms_contacts.id = sms_contact_groups.sms_contact_id").
		Where("sms_contact_groups.sms_group_id IN ?", groupIDs).
		Group("sms_contacts.id").
		Find(&contacts).Error

	if err != nil || len(contacts) == 0 {
		msg := "Selected groups have no contacts"
		if err != nil {
			msg = "Contact resolution error: " + err.Error()
		}
		log.Printf("[SMSCampaignService.processCampaignInBackground] Campaign %d failed contact resolution: %s", campaignID, msg)
		db.DB.Model(&campaign).Update("status", "failed")
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Campaign Failed",
			Message:          msg,
			NotificationType: "ERROR",
			ReferenceID:      &campaignID,
			ReferenceType:    utils.StringPtr("SMS_CAMPAIGN"),
		})
		return
	}

	log.Printf("[SMSCampaignService.processCampaignInBackground] Resolved %d unique contacts for campaign %d", len(contacts), campaignID)

	// 2. Create Initial Pending Outbox entries
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		for _, contact := range contacts {
			outbox := models.SMSOutbox{
				BaseModel:     models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
				SMSCampaignID: &campaign.ID,
				SMSMessageID:  campaign.SMSMessageID,
				SMSContactID:  &contact.ID,
				PhoneNo:       contact.PhoneNumber,
				Status:        "pending",
			}
			if err := tx.Create(&outbox).Error; err != nil {
				return err
			}
		}
		return tx.Model(&campaign).Update("total_recipients", len(contacts)).Error
	})

	if err != nil {
		log.Printf("[SMSCampaignService.processCampaignInBackground] Campaign %d failed outbox initialization: %v", campaignID, err)
		db.DB.Model(&campaign).Update("status", "failed")
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title: "Campaign Failed", Message: "Failed to initialize outbox records.", NotificationType: "ERROR",
		})
		return
	}

	// 3. Get Default Provider
	var provider models.SMSProvider
	if err := db.DB.Where("is_default = ? AND status = ?", "yes", "active").First(&provider).Error; err != nil {
		log.Printf("[SMSCampaignService.processCampaignInBackground] Campaign %d failed: default provider not found", campaignID)
		db.DB.Model(&campaign).Update("status", "failed")
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title: "Campaign Failed", Message: "Default SMS provider not found.", NotificationType: "ERROR",
		})
		return
	}

	// 4. Check Scheduling
	now := utils.Now()
	if campaign.ScheduledAt != nil && campaign.ScheduledAt.After(now) {
		log.Printf("[SMSCampaignService.processCampaignInBackground] Campaign %d scheduled for %s. Stopping current execution.", campaignID, campaign.ScheduledAt)
		db.DB.Model(&campaign).Update("status", "scheduled")
		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "Campaign Prepared",
			Message:          fmt.Sprintf("Campaign '%s' resolved and scheduled for %s.", campaign.CampaignName, campaign.ScheduledAt.Format("2006-01-02 15:04")),
			NotificationType: "INFO",
		})
		return
	}

	db.DB.Model(&campaign).Update("status", "running")

	// 5. Load Pending Outboxes for dispatch
	var outboxes []models.SMSOutbox
	db.DB.Where("sms_campaign_id = ? AND status = ?", campaignID, "pending").Find(&outboxes)

	log.Printf("[SMSCampaignService.processCampaignInBackground] Dispatching %d outbox entries for campaign %d", len(outboxes), campaignID)

	// 6. Worker Pool Dispatching
	numWorkers := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	jobs := make(chan models.SMSOutbox, len(outboxes))
	var mu sync.Mutex
	var sentCount, failedCount int

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for outbox := range jobs {
				reference := fmt.Sprintf("CMP-%d-%d", campaignID, outbox.ID)
				err := s.dispatchSMS(&provider, outbox.PhoneNo, message, reference)

				mu.Lock()
				if err != nil {
					failedCount++
					db.DB.Model(&outbox).Updates(map[string]interface{}{
						"status":           "failed",
						"response_message": err.Error(),
					})
				} else {
					sentCount++
					now := time.Now()
					db.DB.Model(&outbox).Updates(map[string]interface{}{
						"status":  "sent",
						"sent_at": &now,
					})
				}
				mu.Unlock()
			}
		}()
	}

	for _, o := range outboxes {
		jobs <- o
	}
	close(jobs)
	wg.Wait()

	// 4. Finalize Campaign

	db.DB.Model(&campaign).Updates(map[string]interface{}{
		"total_sent":   sentCount,
		"total_failed": failedCount,
		"status":       "completed",
		"completed_at": &now,
	})

	log.Printf("[SMSCampaignService.processCampaignInBackground] Campaign %d completed. Sent: %d, Failed: %d", campaignID, sentCount, failedCount)

	// 5. Notify UI
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "Campaign Processed",
		Message:          fmt.Sprintf("Campaign '%s' completed. Sent: %d, Failed: %d", campaign.CampaignName, sentCount, failedCount),
		NotificationType: "SUCCESS",
		ReferenceID:      &campaignID,
		ReferenceType:    utils.StringPtr("SMS_CAMPAIGN"),
	})
}

func (s *SMSCampaignService) dispatchSMS(provider *models.SMSProvider, recipient, message, reference string) error {
	payload := map[string]interface{}{
		"message":      message,
		"short_code":   provider.SenderID,
		"client_code":  provider.ProviderCode,
		"key":          provider.ApiKey,
		"recipients":   []string{recipient},
		"reference":    reference,
		"message_type": "BULK",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", provider.ApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("provider returned status: %d", resp.StatusCode)
	}

	return nil
}
