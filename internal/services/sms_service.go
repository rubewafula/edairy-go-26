package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type SMSService struct{}

func NewSMSService() *SMSService {
	return &SMSService{}
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

func (s *SMSService) GetGroups(page, limit int) ([]models.SMSGroup, int64, error) {
	var results []models.SMSGroup
	var total int64
	db.DB.Model(&models.SMSGroup{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&results).Error
	return results, total, err
}

// Contacts
func (s *SMSService) CreateContact(req dtos.CreateSMSContactRequest) (*models.SMSContact, error) {
	contact := &models.SMSContact{
		SMSGroupID:  req.SMSGroupID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	err := db.DB.Create(contact).Error
	return contact, err
}

func (s *SMSService) GetContactsByGroup(groupID string) ([]models.SMSContact, error) {
	var contacts []models.SMSContact
	err := db.DB.Where("sms_group_id = ?", groupID).Find(&contacts).Error
	return contacts, err
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

// Direct Messaging & Queue
func (s *SMSService) SendSMS(req dtos.SendSMSRequest) (*models.SMSMessage, error) {
	message := &models.SMSMessage{
		Recipient: req.Recipient,
		Message:   req.Message,
		Status:    "pending",
	}
	if req.GroupID != 0 {
		message.SMSGroupID = &req.GroupID
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
