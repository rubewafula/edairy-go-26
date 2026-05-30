package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type StatutoryDeductionConfigurationService struct{}

func NewStatutoryDeductionConfigurationService() *StatutoryDeductionConfigurationService {
	return &StatutoryDeductionConfigurationService{}
}

func (s *StatutoryDeductionConfigurationService) CreateConfiguration(req dtos.CreateStatutoryDeductionConfigurationRequest, userID uint64) (*models.StatutoryDeductionConfiguration, error) {
	config := &models.StatutoryDeductionConfiguration{
		BaseModel:             models.BaseModel{CreatedBy: userID},
		DeductionID:           req.DeductionID,
		EmployeeDeductionRate: req.EmployeeDeductionRate,
		EmployerDeductionRate: req.EmployerDeductionRate,
		MinAmount:             req.MinAmount,
		FixedAmount:           req.FixedAmount,
		BandLowerLimitAmount:  req.BandLowerLimitAmount,
		BandUpperLimitAmount:  req.BandUpperLimitAmount,
		MinApplicableAmount:   req.MinApplicableAmount,
	}

	if err := db.DB.Create(config).Error; err != nil {
		return nil, err
	}
	return config, nil
}

func (s *StatutoryDeductionConfigurationService) GetConfigurations(page, limit int) ([]dtos.StatutoryDeductionConfigurationResponse, int64, error) {
	var results []dtos.StatutoryDeductionConfigurationResponse
	var total int64
	db.DB.Model(&models.StatutoryDeductionConfiguration{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			sdc.*, edt.name AS deduction_type_name
		FROM statutory_deductions_configurations sdc
		LEFT JOIN employee_deduction_types edt ON sdc.deduction_id = edt.id
		WHERE sdc.deleted_at IS NULL
		ORDER BY sdc.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *StatutoryDeductionConfigurationService) GetConfiguration(id string) (*dtos.StatutoryDeductionConfigurationResponse, error) {
	var result dtos.StatutoryDeductionConfigurationResponse
	query := `
		SELECT 
			sdc.*, edt.name AS deduction_type_name
		FROM statutory_deductions_configurations sdc
		LEFT JOIN employee_deduction_types edt ON sdc.deduction_id = edt.id
		WHERE sdc.id = ? AND sdc.deleted_at IS NULL
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

func (s *StatutoryDeductionConfigurationService) UpdateConfiguration(id string, req dtos.UpdateStatutoryDeductionConfigurationRequest, userID uint64) error {
	var config models.StatutoryDeductionConfiguration
	if err := db.DB.First(&config, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"deduction_id":            req.DeductionID,
		"employee_deduction_rate": req.EmployeeDeductionRate,
		"employer_deduction_rate": req.EmployerDeductionRate,
		"min_amount":              req.MinAmount,
		"fixed_amount":            req.FixedAmount,
		"band_lower_limit_amount": req.BandLowerLimitAmount,
		"band_upper_limit_amount": req.BandUpperLimitAmount,
		"min_applicable_amount":   req.MinApplicableAmount,
		"updated_by":              userID,
	}

	return db.DB.Model(&config).Updates(updates).Error
}

func (s *StatutoryDeductionConfigurationService) DeleteConfiguration(id string) error {
	return db.DB.Delete(&models.StatutoryDeductionConfiguration{}, id).Error
}
