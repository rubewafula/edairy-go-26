package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type SMSCampaignService struct{}

func NewSMSCampaignService() *SMSCampaignService {
	return &SMSCampaignService{}
}

func (s *SMSCampaignService) CreateCampaign(req dtos.CreateSMSCampaignRequest, userID uint64) (*models.SMSCampaign, error) {
	campaign := &models.SMSCampaign{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		CampaignCode: req.CampaignCode,
		CampaignName: req.CampaignName,
		Description:  req.Description,
		Message:      req.Message,
		SMSGroupID:   &req.SMSGroupID,
		Status:       "draft",
		SiteID:       &req.SiteID,
	}

	if req.ScheduledAt != "" {
		t := utils.ParseDate(req.ScheduledAt)
		campaign.ScheduledAt = &t
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(campaign).Error; err != nil {
			return err
		}

		// Pull contacts from group and add to recipients
		var contacts []models.SMSContact
		if err := tx.Where("sms_group_id = ?", req.SMSGroupID).Find(&contacts).Error; err != nil {
			return err
		}

		for _, contact := range contacts {
			recipient := &models.SMSCampaignRecipient{
				SMSCampaignID: campaign.ID,
				RecipientName: contact.Name,
				PhoneNo:       contact.PhoneNumber,
				Status:        "pending",
			}
			if err := tx.Create(recipient).Error; err != nil {
				return err
			}
		}

		return tx.Model(campaign).Update("total_recipients", len(contacts)).Error
	})

	return campaign, err
}

func (s *SMSCampaignService) GetCampaigns(page, limit int) ([]dtos.SMSCampaignResponse, int64, error) {
	var results []dtos.SMSCampaignResponse
	var total int64
	db.DB.Model(&models.SMSCampaign{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT sc.*, sg.name as group_name
		FROM sms_campaigns sc
		LEFT JOIN sms_groups sg ON sc.sms_group_id = sg.id
		WHERE sc.deleted_at IS NULL
		ORDER BY sc.id DESC LIMIT ? OFFSET ?`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SMSCampaignService) GetCampaign(id string) (*dtos.SMSCampaignResponse, error) {
	var result dtos.SMSCampaignResponse
	query := `
		SELECT sc.*, sg.name as group_name
		FROM sms_campaigns sc
		LEFT JOIN sms_groups sg ON sc.sms_group_id = sg.id
		WHERE sc.id = ? AND sc.deleted_at IS NULL
		LIMIT 1`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *SMSCampaignService) UpdateCampaign(id string, req dtos.UpdateSMSCampaignRequest, userID uint64) error {
	var campaign models.SMSCampaign
	if err := db.DB.First(&campaign, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"campaign_code": req.CampaignCode,
		"campaign_name": req.CampaignName,
		"description":   req.Description,
		"message":       req.Message,
		"status":        req.Status,
		"updated_by":    userID,
	}

	if req.SMSGroupID != 0 {
		updates["sms_group_id"] = req.SMSGroupID
	}
	if req.SiteID != 0 {
		updates["site_id"] = req.SiteID
	}
	if req.ScheduledAt != "" {
		t := utils.ParseDate(req.ScheduledAt)
		updates["scheduled_at"] = &t
	}

	return db.DB.Model(&campaign).Updates(updates).Error
}

func (s *SMSCampaignService) DeleteCampaign(id string) error {
	return db.DB.Delete(&models.SMSCampaign{}, id).Error
}

func (s *SMSCampaignService) GetSMSCampaignRecipientsByCampaign(campaignID string) ([]dtos.SMSCampaignRecipientResponse, error) {
	var results []dtos.SMSCampaignRecipientResponse
	err := db.DB.Model(&models.SMSCampaignRecipient{}).Where("sms_campaign_id = ?", campaignID).Scan(&results).Error
	return results, err
}

func (s *SMSCampaignService) CreateSMSCampaignRecipient(req dtos.CreateSMSCampaignRecipientRequest) (*models.SMSCampaignRecipient, error) {
	status := req.Status
	if status == "" {
		status = "pending" // Default status
	}

	recipient := &models.SMSCampaignRecipient{
		SMSCampaignID: req.SMSCampaignID,
		RecipientName: req.RecipientName,
		PhoneNo:       req.PhoneNo,
		Status:        status,
	}

	if err := db.DB.Create(recipient).Error; err != nil {
		return nil, err
	}
	return recipient, nil
}

func (s *SMSCampaignService) GetAllSMSCampaignRecipients(page, limit int) ([]dtos.SMSCampaignRecipientResponse, int64, error) {
	var results []dtos.SMSCampaignRecipientResponse
	var total int64
	db.DB.Model(&models.SMSCampaignRecipient{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			scr.id, scr.sms_campaign_id, sc.campaign_name, scr.recipient_name, scr.phone_no, scr.status,
			scr.created_at, scr.updated_at
		FROM sms_campaign_recipients scr
		LEFT JOIN sms_campaigns sc ON scr.sms_campaign_id = sc.id
		WHERE scr.deleted_at IS NULL
		ORDER BY scr.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SMSCampaignService) GetSMSCampaignRecipient(id string) (*dtos.SMSCampaignRecipientResponse, error) {
	var result dtos.SMSCampaignRecipientResponse
	query := `
		SELECT 
			scr.id, scr.sms_campaign_id, sc.campaign_name, scr.recipient_name, scr.phone_no, scr.status,
			scr.created_at, scr.updated_at
		FROM sms_campaign_recipients scr
		LEFT JOIN sms_campaigns sc ON scr.sms_campaign_id = sc.id
		WHERE scr.id = ? AND scr.deleted_at IS NULL
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

func (s *SMSCampaignService) UpdateSMSCampaignRecipient(id string, req dtos.UpdateSMSCampaignRecipientRequest) error {
	return db.DB.Model(&models.SMSCampaignRecipient{}).Where("id = ?", id).Updates(map[string]interface{}{
		"recipient_name": req.RecipientName,
		"phone_no":       req.PhoneNo,
		"status":         req.Status,
	}).Error
}

func (s *SMSCampaignService) DeleteSMSCampaignRecipient(id string) error {
	return db.DB.Delete(&models.SMSCampaignRecipient{}, id).Error
}
