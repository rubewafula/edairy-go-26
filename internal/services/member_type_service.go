package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// MemberTypeService provides methods for managing member types.
type MemberTypeService struct{}

// NewMemberTypeService creates a new instance of MemberTypeService.
func NewMemberTypeService() *MemberTypeService {
	return &MemberTypeService{}
}

// CreateMemberType creates a new member type in the database.
func (s *MemberTypeService) CreateMemberType(req dtos.CreateMemberTypeRequest) (*dtos.MemberTypeResponse, error) {
	memberType := &models.MemberType{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(memberType).Error; err != nil {
		return nil, err
	}

	response := &dtos.MemberTypeResponse{
		ID:          memberType.ID,
		Name:        memberType.Name,
		Description: memberType.Description,
		CreatedAt:   memberType.CreatedAt,
		UpdatedAt:   memberType.UpdatedAt,
	}

	return response, nil
}

// GetMemberTypes retrieves a paginated list of member types.
func (s *MemberTypeService) GetMemberTypes(page, limit int) ([]dtos.MemberTypeResponse, int64, error) {
	var results []dtos.MemberTypeResponse
	var total int64

	db.DB.Model(&models.MemberType{}).Where("deleted_at IS NULL").Count(&total)

	offset := (page - 1) * limit

	err := db.DB.Model(&models.MemberType{}).
		Limit(limit).Offset(offset).Order("id DESC").
		Where("deleted_at IS NULL").
		Scan(&results).Error

	return results, total, err
}

// GetMemberType retrieves a single member type by its ID.
func (s *MemberTypeService) GetMemberType(id string) (*dtos.MemberTypeResponse, error) {
	var result dtos.MemberTypeResponse
	err := db.DB.Model(&models.MemberType{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &result, nil
}

// UpdateMemberType updates an existing member type.
func (s *MemberTypeService) UpdateMemberType(id string, req dtos.UpdateMemberTypeRequest) error {
	var memberType models.MemberType
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", id).First(&memberType).Error; err != nil {
		return err
	}

	memberType.Name = req.Name
	memberType.Description = req.Description

	return db.DB.Save(&memberType).Error
}

// DeleteMemberType soft deletes a member type by its ID.
func (s *MemberTypeService) DeleteMemberType(id string) error {
	var memberType models.MemberType
	if err := db.DB.Where("id = ? AND deleted_at IS NULL", id).First(&memberType).Error; err != nil {
		return err
	}
	return db.DB.Delete(&memberType).Error
}
