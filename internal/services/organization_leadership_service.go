package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationLeadershipService struct{}

func NewOrganizationLeadershipService() *OrganizationLeadershipService {
	return &OrganizationLeadershipService{}
}

func (s *OrganizationLeadershipService) CreateLeadership(req dtos.CreateOrganizationLeadershipRequest, userID uint64) (*models.OrganizationLeadership, error) {
	leadership := &models.OrganizationLeadership{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IDNo:         req.IDNo,
		Position:     req.Position,
		PrimaryPhone: req.Phone,
		Email:        req.Email,
		Status:       "active",
		LinkStatus:   "pending",
		Submitted:    false,
	}
	if err := db.DB.Create(leadership).Error; err != nil {
		return nil, err
	}
	return leadership, nil
}

func (s *OrganizationLeadershipService) GetLeaderships(page, limit int) ([]dtos.OrganizationLeadershipResponse, int64, error) {
	var results []dtos.OrganizationLeadershipResponse
	var total int64
	db.DB.Model(&models.OrganizationLeadership{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.OrganizationLeadership{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *OrganizationLeadershipService) GetLeadership(id string) (*dtos.OrganizationLeadershipResponse, error) {
	var result dtos.OrganizationLeadershipResponse
	err := db.DB.Model(&models.OrganizationLeadership{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *OrganizationLeadershipService) UpdateLeadership(id string, req dtos.UpdateOrganizationLeadershipRequest, userID uint64) error {
	var leadership models.OrganizationLeadership
	if err := db.DB.First(&leadership, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
		"id_no":         req.IDNo,
		"position":      req.Position,
		"primary_phone": req.Phone,
		"email":         req.Email,
		"updated_by":    userID,
	}

	return db.DB.Model(&leadership).Updates(updates).Error
}

func (s *OrganizationLeadershipService) DeleteLeadership(id string, userID uint64) error {
	var leadership models.OrganizationLeadership
	if err := db.DB.First(&leadership, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&leadership).Update("updated_by", userID).Delete(&leadership).Error
}

func (s *OrganizationLeadershipService) GetLeadershipByNationalID(nationalID string) (*dtos.OrganizationLeadershipResponse, error) {
	var result dtos.OrganizationLeadershipResponse
	err := db.DB.Model(&models.OrganizationLeadership{}).Where("id_no = ?", nationalID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
