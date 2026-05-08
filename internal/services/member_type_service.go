package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type MemberTypeService struct{}

func NewMemberTypeService() *MemberTypeService {
	return &MemberTypeService{}
}

func (s *MemberTypeService) CreateMemberType(req dtos.CreateMemberTypeRequest) (*models.MemberType, error) {
	memberType := &models.MemberType{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(memberType).Error; err != nil {
		return nil, err
	}
	return memberType, nil
}

func (s *MemberTypeService) GetMemberTypes() ([]models.MemberType, int64, error) {
	var memberTypes []models.MemberType
	var total int64
	db.DB.Model(&models.MemberType{}).Count(&total)
	err := db.DB.Find(&memberTypes).Error
	return memberTypes, total, err
}

func (s *MemberTypeService) GetMemberType(id string) (*models.MemberType, error) {
	var memberType models.MemberType
	if err := db.DB.First(&memberType, id).Error; err != nil {
		return nil, err
	}
	return &memberType, nil
}

func (s *MemberTypeService) UpdateMemberType(id string, req dtos.UpdateMemberTypeRequest) error {
	var memberType models.MemberType
	if err := db.DB.First(&memberType, id).Error; err != nil {
		return err
	}

	memberType.Name = req.Name
	memberType.Description = req.Description

	return db.DB.Save(&memberType).Error
}

func (s *MemberTypeService) DeleteMemberType(id string) error {
	return db.DB.Delete(&models.MemberType{}, id).Error
}
