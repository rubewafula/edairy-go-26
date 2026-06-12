package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type SMSService struct {
	notificationService *UINotificationService
}

func NewSMSService() *SMSService {
	return &SMSService{
		notificationService: NewUINotificationService(),
	}
}

// Groups
func (s *SMSService) CreateGroup(req dtos.CreateSMSGroupRequest) (*models.SMSGroup, error) {
	group := &models.SMSGroup{
		Name:        req.Name,
		Description: req.Description,
	}
	err := db.DB.Create(group).Error
	return group, err
}

func (s *SMSService) GetGroup(id string) (*models.SMSGroup, error) {
	var group models.SMSGroup
	err := db.DB.Where("id = ? AND deleted_at IS NULL", id).First(&group).Error
	return &group, err
}

func (s *SMSService) GetGroups(page, limit int) ([]models.SMSGroup, int64, error) {
	var results []models.SMSGroup
	var total int64
	db.DB.Model(&models.SMSGroup{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&results).Error
	return results, total, err
}

func (s *SMSService) UpdateGroup(id string, req dtos.CreateSMSGroupRequest) error {
	var group models.SMSGroup
	if err := db.DB.First(&group, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&group).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}).Error
}

func (s *SMSService) DeleteGroup(id string) error {
	var group models.SMSGroup
	if err := db.DB.First(&group, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&group).Error
}

// Contacts
func (s *SMSService) CreateContact(req dtos.CreateSMSContactRequest) (*models.SMSContact, error) {
	contact := &models.SMSContact{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	err := db.DB.Create(contact).Error
	return contact, err
}

func (s *SMSService) GetContacts(page, limit int) ([]models.SMSContact, int64, error) {
	var results []models.SMSContact
	var total int64
	db.DB.Model(&models.SMSContact{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&results).Error
	return results, total, err
}

func (s *SMSService) GetContactsByGroup(groupID string) ([]models.SMSContact, error) {
	var contacts []models.SMSContact
	err := db.DB.Where("sms_group_id = ?", groupID).Find(&contacts).Error
	return contacts, err
}

func (s *SMSService) GetContact(id string) (*models.SMSContact, error) {
	var contact models.SMSContact
	err := db.DB.Where("id = ? AND deleted_at IS NULL", id).First(&contact).Error
	return &contact, err
}

func (s *SMSService) UpdateContact(id string, req dtos.CreateSMSContactRequest) error {
	var contact models.SMSContact
	if err := db.DB.First(&contact, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&contact).Updates(map[string]interface{}{
		"name":         req.Name,
		"phone_number": req.PhoneNumber,
	}).Error
}

func (s *SMSService) DeleteContact(id string) error {
	var contact models.SMSContact
	if err := db.DB.First(&contact, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&contact).Error
}

// Providers
func (s *SMSService) CreateProvider(req dtos.CreateSMSProviderRequest, userID uint64) (*models.SMSProvider, error) {
	provider := &models.SMSProvider{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		ProviderCode: req.ProviderCode,
		ProviderName: req.ProviderName,
		ApiUrl:       req.ApiUrl,
		ApiKey:       req.ApiKey,
		ApiSecret:    req.ApiSecret,
		SenderID:     req.SenderID,
		Username:     req.Username,
		Password:     req.Password,
		IsDefault:    req.IsDefault,
		Status:       req.Status,
		Notes:        req.Notes,
	}
	if err := db.DB.Create(provider).Error; err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *SMSService) GetProviders() ([]dtos.SMSProviderResponse, error) {
	var results []dtos.SMSProviderResponse
	err := db.DB.Model(&models.SMSProvider{}).Scan(&results).Error
	return results, err
}

func (s *SMSService) GetProvider(id string) (*dtos.SMSProviderResponse, error) {
	var result dtos.SMSProviderResponse
	err := db.DB.Model(&models.SMSProvider{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SMSService) UpdateProvider(id string, req dtos.CreateSMSProviderRequest, userID uint64) error {
	var provider models.SMSProvider
	if err := db.DB.First(&provider, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"provider_code": req.ProviderCode,
		"provider_name": req.ProviderName,
		"api_url":       req.ApiUrl,
		"api_key":       req.ApiKey,
		"api_secret":    req.ApiSecret,
		"sender_id":     req.SenderID,
		"username":      req.Username,
		"password":      req.Password,
		"is_default":    req.IsDefault,
		"status":        req.Status,
		"notes":         req.Notes,
		"updated_by":    userID,
	}
	return db.DB.Model(&provider).Updates(updates).Error
}

func (s *SMSService) DeleteProvider(id string) error {
	var provider models.SMSProvider
	if err := db.DB.First(&provider, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&provider).Error
}

// Templates
func (s *SMSService) CreateTemplate(req dtos.CreateSMSTemplateRequest, userID uint64) (*models.SMSTemplate, error) {
	template := &models.SMSTemplate{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		TemplateCode: req.TemplateCode,
		TemplateName: req.TemplateName,
		ModuleName:   req.ModuleName,
		Message:      req.Message,
		Variables:    req.Variables,
		Status:       req.Status,
	}
	err := db.DB.Create(template).Error
	return template, err
}

func (s *SMSService) GetTemplates() ([]dtos.SMSTemplateResponse, error) {
	var results []dtos.SMSTemplateResponse
	err := db.DB.Model(&models.SMSTemplate{}).Scan(&results).Error
	return results, err
}

func (s *SMSService) GetTemplate(id string) (*dtos.SMSTemplateResponse, error) {
	var result dtos.SMSTemplateResponse
	err := db.DB.Model(&models.SMSTemplate{}).Where("id = ? AND deleted_at IS NULL", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SMSService) UpdateTemplate(id string, req dtos.CreateSMSTemplateRequest, userID uint64) error {
	var template models.SMSTemplate
	if err := db.DB.First(&template, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"template_code": req.TemplateCode,
		"template_name": req.TemplateName,
		"module_name":   req.ModuleName,
		"message":       req.Message,
		"variables":     req.Variables,
		"status":        req.Status,
		"updated_by":    userID,
	}
	return db.DB.Model(&template).Updates(updates).Error
}

func (s *SMSService) DeleteTemplate(id string, userID uint64) error {
	var template models.SMSTemplate
	if err := db.DB.First(&template, id).Error; err != nil {
		return err
	}
	// Perform a soft delete and update updated_by
	return db.DB.Model(&template).Update("updated_by", userID).Delete(&template).Error
}

// Direct Messaging & Queue
func (s *SMSService) SendSMS(req dtos.SendSMSRequest) (*models.SMSMessage, error) {
	message := &models.SMSMessage{
		Message: req.Message,
		Status:  "pending",
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		queueItem := &models.SMSQueue{
			SMSMessageID: message.ID,
			Processed:    "no",
		}
		return tx.Create(queueItem).Error
	})

	return message, err
}

func (s *SMSService) GetQueue(page, limit int) ([]dtos.SMSQueueResponse, int64, error) {
	var results []dtos.SMSQueueResponse
	var total int64
	db.DB.Model(&models.SMSQueue{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Model(&models.SMSQueue{}).Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *SMSService) GetMessages(page, limit int) ([]dtos.SMSMessageResponse, int64, error) {
	var results []dtos.SMSMessageResponse
	var total int64
	db.DB.Model(&models.SMSMessage{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Model(&models.SMSMessage{}).Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *SMSService) GetMessage(id string) (*dtos.SMSMessageResponse, error) {
	var result dtos.SMSMessageResponse
	err := db.DB.Model(&models.SMSMessage{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (s *SMSService) CreateMessage(req dtos.SendSMSRequest) (*models.SMSMessage, error) {
	// Leveraging SendSMS to handle both record creation and queuing
	return s.SendSMS(req)
}

func (s *SMSService) UpdateMessage(id string, req dtos.SendSMSRequest) error {
	var message models.SMSMessage
	if err := db.DB.First(&message, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"recipient": req.Recipient,
		"message":   req.Message,
	}
	if req.GroupID != 0 {
		updates["sms_group_id"] = req.GroupID
	}

	return db.DB.Model(&message).Updates(updates).Error
}

func (s *SMSService) DeleteMessage(id string) error {
	var message models.SMSMessage
	if err := db.DB.First(&message, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&message).Error
}

// SMS Outbox related methods

func (s *SMSService) GetSMSOutboxesByCampaign(campaignID string) ([]dtos.SMSOutboxResponse, error) {
	var results []dtos.SMSOutboxResponse
	query := `
		SELECT 
			so.*, 
			sc.campaign_name, 
			sm.message AS message_text, 
			scon.name AS contact_name
		FROM sms_outboxes so
		LEFT JOIN sms_campaigns sc ON so.sms_campaign_id = sc.id
		LEFT JOIN sms_messages sm ON so.sms_message_id = sm.id
		LEFT JOIN sms_contacts scon ON so.sms_contact_id = scon.id
		WHERE so.sms_campaign_id = ? AND so.deleted_at IS NULL
		ORDER BY so.id DESC
	`
	err := db.DB.Raw(query, campaignID).Scan(&results).Error
	return results, err
}

func (s *SMSService) CreateSMSOutbox(req dtos.CreateSMSOutboxRequest, userID uint64) (*models.SMSOutbox, error) {
	// 1. Create or Resolve SMS Contact
	var contact models.SMSContact
	if err := db.DB.Where("phone_number = ?", req.PhoneNo).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			contact = models.SMSContact{
				Name:        req.PhoneNo, // Default name to phone number
				PhoneNumber: req.PhoneNo,
			}
			if err := db.DB.Create(&contact).Error; err != nil {
				return nil, fmt.Errorf("failed to auto-create contact: %w", err)
			}
		} else {
			return nil, err
		}
	}

	// 2. Create entry in sms_messages with status pending
	smsMsg := &models.SMSMessage{
		Message: req.Message,
		Status:  "pending",
	}
	if err := db.DB.Create(smsMsg).Error; err != nil {
		return nil, fmt.Errorf("failed to create sms message record: %w", err)
	}

	// 3. Create the outbox entry with associated IDs
	outbox := &models.SMSOutbox{
		BaseModel:    models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		PhoneNo:      req.PhoneNo,
		SMSMessageID: &smsMsg.ID,
		SMSContactID: &contact.ID,
		Status:       "pending",
	}

	if err := db.DB.Create(outbox).Error; err != nil {
		return nil, err
	}

	// 4. Query default SMS provider
	var provider models.SMSProvider
	if err := db.DB.Where("is_default = ? AND status = ?", "yes", "active").First(&provider).Error; err != nil {
		s.handleOutboxFailure(outbox, "Default SMS provider not found or inactive", userID)
		return outbox, nil
	}

	// 5. Dispatch the actual SMS
	reference := fmt.Sprintf("%d%d", outbox.ID, time.Now().Unix())
	if err := s.dispatchSMS(&provider, outbox.PhoneNo, req.Message, reference); err != nil {
		s.handleOutboxFailure(outbox, err.Error(), userID)
	} else {
		// 6. Update status on success
		now := utils.Now()
		db.DB.Model(outbox).Updates(map[string]interface{}{
			"status":  "sent",
			"sent_at": &now,
		})
		db.DB.Model(smsMsg).Update("status", "sent")

		s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
			Title:            "SMS Sent",
			Message:          fmt.Sprintf("Message to %s dispatched successfully via %s", outbox.PhoneNo, provider.ProviderName),
			NotificationType: "SUCCESS",
		})
	}

	return outbox, nil
}

func (s *SMSService) GetAllSMSOutboxes(page, limit int) ([]dtos.SMSOutboxResponse, int64, error) {
	var results []dtos.SMSOutboxResponse
	var total int64

	db.DB.Model(&models.SMSOutbox{}).Count(&total)

	offset := (page - 1) * limit
	query := `
		SELECT 
			so.*, 
			sc.campaign_name, 
			sm.message AS message_text, 
			scon.name AS contact_name
		FROM sms_outboxes so
		LEFT JOIN sms_campaigns sc ON so.sms_campaign_id = sc.id
		LEFT JOIN sms_messages sm ON so.sms_message_id = sm.id
		LEFT JOIN sms_contacts scon ON so.sms_contact_id = scon.id
		WHERE so.deleted_at IS NULL
		ORDER BY so.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}
	return results, total, err
}

func (s *SMSService) GetSMSOutbox(id string) (*dtos.SMSOutboxResponse, error) {
	var result dtos.SMSOutboxResponse
	query := `
		SELECT 
			so.*, 
			sc.campaign_name, 
			sm.message AS message_text, 
			scon.name AS contact_name
		FROM sms_outboxes so
		LEFT JOIN sms_campaigns sc ON so.sms_campaign_id = sc.id
		LEFT JOIN sms_messages sm ON so.sms_message_id = sm.id
		LEFT JOIN sms_contacts scon ON so.sms_contact_id = scon.id
		WHERE so.id = ? AND so.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *SMSService) UpdateSMSOutbox(id string, req dtos.UpdateSMSOutboxRequest, userID uint64) error {
	var outbox models.SMSOutbox
	if err := db.DB.First(&outbox, id).Error; err != nil {
		return err
	}

	outbox.UpdatedBy = userID
	if req.Status != nil {
		outbox.Status = *req.Status
	}
	if req.SentAt != nil {
		outbox.SentAt = req.SentAt
	}
	if req.DeliveredAt != nil {
		outbox.DeliveredAt = req.DeliveredAt
	}

	return db.DB.Save(&outbox).Error
}

func (s *SMSService) DeleteSMSOutbox(id string) error {

	var outbox models.SMSOutbox
	if err := db.DB.First(&outbox, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&outbox).Error
}

func (s *SMSService) dispatchSMS(provider *models.SMSProvider, recipient, message, reference string) error {
	payload := map[string]interface{}{
		"message":      message,
		"short_code":   provider.SenderID,
		"link_id":      nil,
		"call_back":    "",
		"client_code":  provider.ProviderCode,
		"key":          provider.ApiKey,
		"recipients":   []string{recipient},
		"reference":    reference,
		"message_type": "BULK",
		"bulk":         1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal sms payload: %w", err)
	}

	req, err := http.NewRequest("POST", provider.ApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("provider connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sms provider error status: %d", resp.StatusCode)
	}

	return nil
}

func (s *SMSService) handleOutboxFailure(outbox *models.SMSOutbox, errMsg string, userID uint64) {
	db.DB.Model(outbox).Update("status", "failed")
	s.notificationService.CreateNotification(userID, dtos.CreateUINotificationRequest{
		Title:            "SMS Dispatch Failed",
		Message:          fmt.Sprintf("Failed to send SMS to %s: %s", outbox.PhoneNo, errMsg),
		NotificationType: "ERROR",
	})
}

// In-App Configurations
func (s *SMSService) CreateInAppConfig(req dtos.CreateSMSInAppConfigurationRequest, userID uint64) (*dtos.SMSInAppConfigurationResponse, error) {
	config := &models.SMSInAppConfiguration{
		BaseModel:           models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		ActivityCode:        req.ActivityCode,
		ActivityDescription: req.ActivityDescription,
		IsEnabled:           req.IsEnabled,
	}
	if err := db.DB.Create(config).Error; err != nil {
		return nil, err
	}
	return s.toSMSInAppConfigurationResponse(config), nil
}

func (s *SMSService) GetInAppConfigs(page, limit int) ([]dtos.SMSInAppConfigurationResponse, int64, error) {
	var configs []models.SMSInAppConfiguration
	var total int64
	db.DB.Model(&models.SMSInAppConfiguration{}).Count(&total)
	offset := (page - 1) * limit
	if err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	var results []dtos.SMSInAppConfigurationResponse
	for _, config := range configs {
		results = append(results, *s.toSMSInAppConfigurationResponse(&config))
	}
	return results, total, nil
}

func (s *SMSService) GetInAppConfig(id string) (*dtos.SMSInAppConfigurationResponse, error) {
	var config models.SMSInAppConfiguration
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", id).First(&config).Error; err != nil {
		return nil, err
	}
	return s.toSMSInAppConfigurationResponse(&config), nil
}

func (s *SMSService) UpdateInAppConfig(id string, req dtos.UpdateSMSInAppConfigurationRequest, userID uint64) error {
	var config models.SMSInAppConfiguration
	if err := db.DB.First(&config, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"activity_description": req.ActivityDescription,
		"is_enabled":           req.IsEnabled,
		"updated_by":           userID,
	}
	return db.DB.Model(&config).Updates(updates).Error
}

func (s *SMSService) DeleteInAppConfig(id string, userID uint64) error {
	var config models.SMSInAppConfiguration
	if err := db.DB.First(&config, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&config).Update("updated_by", userID).Delete(&config).Error
}

// toSMSInAppConfigurationResponse converts a models.SMSInAppConfiguration to a dtos.SMSInAppConfigurationResponse
func (s *SMSService) toSMSInAppConfigurationResponse(config *models.SMSInAppConfiguration) *dtos.SMSInAppConfigurationResponse {
	if config == nil {
		return nil
	}
	return &dtos.SMSInAppConfigurationResponse{
		ID:                  config.ID,
		ActivityCode:        config.ActivityCode,
		ActivityDescription: config.ActivityDescription,
		IsEnabled:           config.IsEnabled,
		CreatedAt:           config.CreatedAt,
		UpdatedAt:           config.UpdatedAt,
	}
}
