package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type LoanOrganizationProfileService struct{}

func NewLoanOrganizationProfileService() *LoanOrganizationProfileService {
	return &LoanOrganizationProfileService{}
}

func (s *LoanOrganizationProfileService) CreateProfile(req dtos.CreateLoanOrganizationProfileRequest, userID uint64) (*models.LoanOrganizationProfile, error) {
	profile := &models.LoanOrganizationProfile{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		NextLevel:       req.NextLevel,
		AstraID:         req.AstraID,
		LinkStatus:      req.LinkStatus,
		UUID:            req.UUID,
		Version:         req.Version,
		ProductID:       req.ProductID,
		CompanyDetailID: req.CompanyDetailID,
		ManuallyRatify:  req.ManuallyRatify,
	}
	if err := db.DB.Create(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *LoanOrganizationProfileService) GetProfiles(page, limit int) ([]dtos.LoanOrganizationProfileResponse, int64, error) {
	var results []dtos.LoanOrganizationProfileResponse
	var total int64
	db.DB.Model(&models.LoanOrganizationProfile{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.LoanOrganizationProfile{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *LoanOrganizationProfileService) GetProfile(id string) (*dtos.LoanOrganizationProfileResponse, error) {
	var result dtos.LoanOrganizationProfileResponse
	err := db.DB.Model(&models.LoanOrganizationProfile{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *LoanOrganizationProfileService) UpdateProfile(id string, req dtos.UpdateLoanOrganizationProfileRequest, userID uint64) error {
	var profile models.LoanOrganizationProfile
	if err := db.DB.First(&profile, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"next_level":        req.NextLevel,
		"astra_id":          req.AstraID,
		"link_status":       req.LinkStatus,
		"uuid":              req.UUID,
		"version":           req.Version,
		"product_id":        req.ProductID,
		"company_detail_id": req.CompanyDetailID,
		"manually_ratify":   req.ManuallyRatify,
		"updated_by":        userID,
	}

	return db.DB.Model(&profile).Updates(updates).Error
}

func (s *LoanOrganizationProfileService) DeleteProfile(id string, userID uint64) error {
	var profile models.LoanOrganizationProfile
	if err := db.DB.First(&profile, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&profile).Update("updated_by", userID).Delete(&profile).Error
}
