package services

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type DividendDeclarationService struct{}

func NewDividendDeclarationService() *DividendDeclarationService {
	return &DividendDeclarationService{}
}

func (s *DividendDeclarationService) CreateDeclaration(req dtos.CreateDividendDeclarationRequest) (*models.DividendDeclaration, error) {
	status := req.Status
	if status == "" {
		status = "DRAFT"
	}

	declaration := &models.DividendDeclaration{
		FiscalYear:      req.FiscalYear,
		Period:          req.Period,
		TotalPool:       req.TotalPool,
		RatePerShare:    req.RatePerShare,
		CalculationType: req.CalculationType,
		Status:          status,
		ApprovedBy:      req.ApprovedBy,
	}

	if req.ApprovedAt != "" {
		declaration.ApprovedAt = utils.ParseDate(req.ApprovedAt)
	}

	if err := db.DB.Create(declaration).Error; err != nil {
		return nil, err
	}
	return declaration, nil
}

func (s *DividendDeclarationService) GetDeclarations() ([]dtos.DividendDeclarationResponse, int64, error) {
	var results []dtos.DividendDeclarationResponse
	var total int64
	db.DB.Model(&models.DividendDeclaration{}).Count(&total)

	query := `
		SELECT 
			dd.id, dd.fiscal_year, dd.period, dd.total_pool, dd.rate_per_share,
			dd.calculation_type, dd.status, dd.approved_by, u.name AS approved_by_user_name,
			dd.approved_at, dd.created_at, dd.updated_at
		FROM dividend_declarations dd
		LEFT JOIN users u ON dd.approved_by = u.id
		WHERE dd.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *DividendDeclarationService) GetDeclaration(id string) (*dtos.DividendDeclarationResponse, error) {
	var result dtos.DividendDeclarationResponse
	query := `
		SELECT 
			dd.id, dd.fiscal_year, dd.period, dd.total_pool, dd.rate_per_share,
			dd.calculation_type, dd.status, dd.approved_by, u.name AS approved_by_user_name,
			dd.approved_at, dd.created_at, dd.updated_at
		FROM dividend_declarations dd
		LEFT JOIN users u ON dd.approved_by = u.id
		WHERE dd.id = ? AND dd.deleted_at IS NULL
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

func (s *DividendDeclarationService) UpdateDeclaration(id string, req dtos.UpdateDividendDeclarationRequest) error {
	var declaration models.DividendDeclaration
	if err := db.DB.First(&declaration, id).Error; err != nil {
		return err
	}

	declaration.FiscalYear = req.FiscalYear
	declaration.Period = req.Period
	declaration.TotalPool = req.TotalPool
	declaration.RatePerShare = req.RatePerShare
	declaration.CalculationType = req.CalculationType
	declaration.Status = req.Status
	declaration.ApprovedBy = req.ApprovedBy
	if req.ApprovedAt != "" {
		declaration.ApprovedAt = utils.ParseDate(req.ApprovedAt)
	} else {
		// Optionally set to zero time if approved_at is being cleared
		declaration.ApprovedAt = time.Time{}
	}

	return db.DB.Save(&declaration).Error
}

func (s *DividendDeclarationService) DeleteDeclaration(id string) error {
	return db.DB.Delete(&models.DividendDeclaration{}, id).Error
}
