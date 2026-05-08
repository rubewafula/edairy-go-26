package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type ShareTypeService struct{}

func NewShareTypeService() *ShareTypeService {
	return &ShareTypeService{}
}

func (s *ShareTypeService) CreateShareType(req dtos.CreateShareTypeRequest) (*models.ShareType, error) {
	shareType := &models.ShareType{
		ShareCode:         req.ShareCode,
		ShareType:         req.ShareType,
		Description:       req.Description,
		Rate:              req.Rate,
		Mandatory:         req.Mandatory,
		HasShareValue:     req.HasShareValue,
		RepayMethod:       req.RepayMethod,
		CalculatingMethod: req.CalculatingMethod,
		ShareValue:        req.ShareValue,
		DeductionTypeID:   req.DeductionTypeID,
		Priority:          req.Priority,
	}

	if err := db.DB.Create(shareType).Error; err != nil {
		return nil, err
	}
	return shareType, nil
}

func (s *ShareTypeService) GetShareTypes() ([]dtos.ShareTypeResponse, int64, error) {
	var results []dtos.ShareTypeResponse
	var total int64
	db.DB.Model(&models.ShareType{}).Count(&total)

	query := `
		SELECT 
			st.id, st.share_code, st.share_type, st.description, st.rate, 
			st.madatory AS mandatory, st.has_share_value, st.repay_method, 
			st.calculating_method, st.share_value, st.deduction_type_id, 
			dt.description AS deduction_type_name,
			st.priority, st.created_at, st.updated_at
		FROM share_types st
		LEFT JOIN deduction_types dt ON st.deduction_type_id = dt.id
		WHERE st.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ShareTypeService) GetShareType(id string) (*dtos.ShareTypeResponse, error) {
	var result dtos.ShareTypeResponse
	query := `
		SELECT 
			st.id, st.share_code, st.share_type, st.description, st.rate, 
			st.madatory AS mandatory, st.has_share_value, st.repay_method, 
			st.calculating_method, st.share_value, st.deduction_type_id, 
			dt.description AS deduction_type_name,
			st.priority, st.created_at, st.updated_at
		FROM share_types st
		LEFT JOIN deduction_types dt ON st.deduction_type_id = dt.id
		WHERE st.id = ? AND st.deleted_at IS NULL
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

func (s *ShareTypeService) UpdateShareType(id string, req dtos.UpdateShareTypeRequest) error {
	var shareType models.ShareType
	if err := db.DB.First(&shareType, id).Error; err != nil {
		return err
	}

	shareType.ShareCode = req.ShareCode
	shareType.ShareType = req.ShareType
	shareType.Description = req.Description
	shareType.Rate = req.Rate
	shareType.Mandatory = req.Mandatory
	shareType.HasShareValue = req.HasShareValue
	shareType.RepayMethod = req.RepayMethod
	shareType.CalculatingMethod = req.CalculatingMethod
	shareType.ShareValue = req.ShareValue
	shareType.DeductionTypeID = req.DeductionTypeID
	shareType.Priority = req.Priority

	return db.DB.Save(&shareType).Error
}

func (s *ShareTypeService) DeleteShareType(id string) error {
	return db.DB.Delete(&models.ShareType{}, id).Error
}
