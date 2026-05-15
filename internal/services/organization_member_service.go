package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationMemberService struct{}

func NewOrganizationMemberService() *OrganizationMemberService {
	return &OrganizationMemberService{}
}

func (s *OrganizationMemberService) CreateMember(req dtos.CreateOrganizationMemberRequest, userID uint64) (*models.OrganizationMember, error) {
	member := &models.OrganizationMember{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		CustomerID:     req.CustomerID,
		CustomerType:   req.CustomerType,
		ManuallyRatify: req.ManuallyRatify,
		NextLevel:      req.NextLevel,
		Status:         req.Status,
		AstraID:        req.AstraID,
		CreditLimit:    req.CreditLimit,
		LinkStatus:     req.LinkStatus,
		LivenessPassed: req.LivenessPassed,
		AstraRemarks:   req.AstraRemarks,
		UUID:           req.UUID,
		AuthCreated:    req.AuthCreated,
		Locale:         req.Locale,
	}
	if err := db.DB.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

func (s *OrganizationMemberService) GetMembers(page, limit int) ([]dtos.OrganizationMemberResponse, int64, error) {
	var results []dtos.OrganizationMemberResponse
	var total int64
	db.DB.Model(&models.OrganizationMember{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			om.*, c.full_names as customer_name
		FROM organization_members om
		LEFT JOIN customers c ON om.customer_id = c.id
		WHERE om.deleted_at IS NULL
		ORDER BY om.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *OrganizationMemberService) GetMember(id string) (*dtos.OrganizationMemberResponse, error) {
	var result dtos.OrganizationMemberResponse
	query := `
		SELECT 
			om.*, c.full_names as customer_name
		FROM organization_members om
		LEFT JOIN customers c ON om.customer_id = c.id
		WHERE om.id = ? AND om.deleted_at IS NULL
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

func (s *OrganizationMemberService) UpdateMember(id string, req dtos.UpdateOrganizationMemberRequest, userID uint64) error {
	var member models.OrganizationMember
	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}

	// Update fields from request
	member.ManuallyRatify = req.ManuallyRatify
	member.NextLevel = req.NextLevel
	member.Status = req.Status
	member.AstraID = req.AstraID
	member.CreditLimit = req.CreditLimit
	member.LinkStatus = req.LinkStatus
	member.LivenessPassed = req.LivenessPassed
	member.AstraRemarks = req.AstraRemarks
	member.AuthCreated = req.AuthCreated
	member.Locale = req.Locale
	member.UpdatedBy = userID

	return db.DB.Save(&member).Error
}

func (s *OrganizationMemberService) DeleteMember(id string, userID uint64) error {
	var member models.OrganizationMember
	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&member).Update("updated_by", userID).Delete(&member).Error
}
