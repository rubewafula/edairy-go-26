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

func (s *SMSCampaignService) GetRecipients(campaignID string) ([]dtos.SMSCampaignRecipientResponse, error) {
	var results []dtos.SMSCampaignRecipientResponse
	err := db.DB.Model(&models.SMSCampaignRecipient{}).Where("sms_campaign_id = ?", campaignID).Scan(&results).Error
	return results, err
}
