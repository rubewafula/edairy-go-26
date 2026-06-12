package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type SMSCampaignService struct{}

func NewSMSCampaignService() *SMSCampaignService {
	return &SMSCampaignService{}
}

func (s *SMSCampaignService) CreateCampaign(req dtos.CreateSMSCampaignRequest, userID uint64) (*models.SMSCampaign, error) {
	campaign := &models.SMSCampaign{
		BaseModel:    models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		CampaignName: req.CampaignName,
		CampaignCode: req.CampaignCode,
		Description:  req.Description,
		Status:       "draft", // Default status
	}

	if req.ScheduledAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_at format: %w", err)
		}
		campaign.ScheduledAt = &parsedTime
	}

	if err := db.DB.Create(campaign).Error; err != nil {
		return nil, err
	}
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
