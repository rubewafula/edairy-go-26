package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type ShareDividendService struct{}

func NewShareDividendService() *ShareDividendService {
	return &ShareDividendService{}
}

func (s *ShareDividendService) CreateDividend(req dtos.CreateShareDividendRequest) (*models.ShareDividend, error) {
	dividend := &models.ShareDividend{
		DeclarationID: req.DeclarationID,
		MemberID:      req.MemberID,
		FiscalYear:    req.FiscalYear,
		Period:        req.Period,
		ShareUnits:    req.ShareUnits,
		Status:        req.Status,
		RatePerShare:  req.RatePerShare,
		TaxAmount:     req.TaxAmount,
		NetAmount:     req.NetAmount,
		TransactionID: req.TransactionID,
	}

	if err := db.DB.Create(dividend).Error; err != nil {
		return nil, err
	}
	return dividend, nil
}

func (s *ShareDividendService) GetDividends() ([]dtos.ShareDividendResponse, int64, error) {
	var results []dtos.ShareDividendResponse
	var total int64
	db.DB.Model(&models.ShareDividend{}).Count(&total)

	query := `
		SELECT 
			sd.id, sd.declaration_id, sd.member_id, m.member_no, m.first_name, m.last_name,
			sd.fiscal_year, sd.period, sd.share_units, sd.status,
			sd.rate_per_share, sd.tax_amount, sd.net_amount, sd.transaction_id,
			sd.created_at, sd.updated_at
		FROM share_dividends sd
		LEFT JOIN member_registrations m ON sd.member_id = m.id
		WHERE sd.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ShareDividendService) GetDividend(id string) (*dtos.ShareDividendResponse, error) {
	var result dtos.ShareDividendResponse
	query := `
		SELECT 
			sd.id, sd.declaration_id, sd.member_id, m.member_no, m.first_name, m.last_name,
			sd.fiscal_year, sd.period, sd.share_units, sd.status,
			sd.rate_per_share, sd.tax_amount, sd.net_amount, sd.transaction_id,
			sd.created_at, sd.updated_at
		FROM share_dividends sd
		LEFT JOIN member_registrations m ON sd.member_id = m.id
		WHERE sd.id = ? AND sd.deleted_at IS NULL
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

func (s *ShareDividendService) UpdateDividend(id string, req dtos.UpdateShareDividendRequest) error {
	var dividend models.ShareDividend
	if err := db.DB.First(&dividend, id).Error; err != nil {
		return err
	}

	dividend.DeclarationID = req.DeclarationID
	dividend.MemberID = req.MemberID
	dividend.FiscalYear = req.FiscalYear
	dividend.Period = req.Period
	dividend.ShareUnits = req.ShareUnits
	dividend.Status = req.Status
	dividend.RatePerShare = req.RatePerShare
	dividend.TaxAmount = req.TaxAmount
	dividend.NetAmount = req.NetAmount
	dividend.TransactionID = req.TransactionID

	return db.DB.Save(&dividend).Error
}

func (s *ShareDividendService) DeleteDividend(id string) error {
	return db.DB.Delete(&models.ShareDividend{}, id).Error
}
