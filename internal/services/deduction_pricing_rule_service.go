package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type DeductionPricingRuleService struct{}

func NewDeductionPricingRuleService() *DeductionPricingRuleService {
	return &DeductionPricingRuleService{}
}

func (s *DeductionPricingRuleService) CreateRule(req dtos.CreateDeductionPricingRuleRequest) (*models.DeductionPricingRule, error) {
	rule := &models.DeductionPricingRule{
		DeductionTypeID: req.DeductionTypeID,
		MinCreditLimit:  req.MinCreditLimit,
		MaxLimit:        req.MaxLimit,
		BoardingFee:     req.BoardingFee,
		ProcessingFee:   req.ProcessingFee,
		InsuranceFee:    req.InsuranceFee,
		LegalFee:        req.LegalFee,
		InterestRate:    req.InterestRate,
		Status:          req.Status,
	}

	if rule.Status == "" {
		rule.Status = "ACTIVE"
	}

	if err := db.DB.Create(rule).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *DeductionPricingRuleService) GetRules(page, limit int) ([]dtos.DeductionPricingRuleResponse, int64, error) {
	var results []dtos.DeductionPricingRuleResponse
	var total int64
	db.DB.Model(&models.DeductionPricingRule{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			dpr.*, dt.description AS deduction_type_name
		FROM deduction_pricing_rules dpr
		LEFT JOIN deduction_types dt ON dpr.deduction_type_id = dt.id
		WHERE dpr.deleted_at IS NULL
		ORDER BY dpr.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *DeductionPricingRuleService) GetRule(id string) (*dtos.DeductionPricingRuleResponse, error) {
	var result dtos.DeductionPricingRuleResponse
	query := `
		SELECT 
			dpr.*, dt.description AS deduction_type_name
		FROM deduction_pricing_rules dpr
		LEFT JOIN deduction_types dt ON dpr.deduction_type_id = dt.id
		WHERE dpr.id = ? AND dpr.deleted_at IS NULL
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

func (s *DeductionPricingRuleService) UpdateRule(id string, req dtos.UpdateDeductionPricingRuleRequest) error {
	var rule models.DeductionPricingRule
	if err := db.DB.First(&rule, id).Error; err != nil {
		return err
	}

	rule.DeductionTypeID = req.DeductionTypeID
	rule.MinCreditLimit = req.MinCreditLimit
	rule.MaxLimit = req.MaxLimit
	rule.BoardingFee = req.BoardingFee
	rule.ProcessingFee = req.ProcessingFee
	rule.InsuranceFee = req.InsuranceFee
	rule.LegalFee = req.LegalFee
	rule.InterestRate = req.InterestRate
	rule.Status = req.Status

	return db.DB.Save(&rule).Error
}

func (s *DeductionPricingRuleService) DeleteRule(id string) error {
	return db.DB.Delete(&models.DeductionPricingRule{}, id).Error
}
