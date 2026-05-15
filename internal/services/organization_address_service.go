package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationAddressService struct{}

func NewOrganizationAddressService() *OrganizationAddressService {
	return &OrganizationAddressService{}
}

func (s *OrganizationAddressService) CreateAddress(req dtos.CreateOrganizationAddressRequest, userID uint64) (*models.OrganizationAddress, error) {
	address := &models.OrganizationAddress{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		AddressType: req.AddressType,
		City:        req.City,
		Code:        req.Code,
		Country:     req.Country,
		Line1:       req.Line1,
		Line2:       req.Line2,
		Line3:       req.Line3,
		State:       req.State,
	}
	if err := db.DB.Create(address).Error; err != nil {
		return nil, err
	}
	return address, nil
}

func (s *OrganizationAddressService) GetAddresses(page, limit int) ([]dtos.OrganizationAddressResponse, int64, error) {
	var results []dtos.OrganizationAddressResponse
	var total int64
	db.DB.Model(&models.OrganizationAddress{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.OrganizationAddress{}).
		Limit(limit).Offset(offset).Order("id DESC").Scan(&results).Error
	return results, total, err
}

func (s *OrganizationAddressService) GetAddress(id string) (*dtos.OrganizationAddressResponse, error) {
	var result dtos.OrganizationAddressResponse
	err := db.DB.Model(&models.OrganizationAddress{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *OrganizationAddressService) UpdateAddress(id string, req dtos.UpdateOrganizationAddressRequest, userID uint64) error {
	var address models.OrganizationAddress
	if err := db.DB.First(&address, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"address_type": req.AddressType,
		"city":         req.City,
		"code":         req.Code,
		"country":      req.Country,
		"line1":        req.Line1,
		"line2":        req.Line2,
		"line3":        req.Line3,
		"state":        req.State,
		"updated_by":   userID,
	}

	return db.DB.Model(&address).Updates(updates).Error
}

func (s *OrganizationAddressService) DeleteAddress(id string, userID uint64) error {
	var address models.OrganizationAddress
	if err := db.DB.First(&address, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&address).Update("updated_by", userID).Delete(&address).Error
}
