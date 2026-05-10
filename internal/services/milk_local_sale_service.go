package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MilkLocalSaleService struct{}

func NewMilkLocalSaleService() *MilkLocalSaleService {
	return &MilkLocalSaleService{}
}

func (s *MilkLocalSaleService) CreateMilkLocalSale(req dtos.CreateMilkLocalSaleRequest) (*models.MilkLocalSale, error) {
	sale := &models.MilkLocalSale{
		Quantity:        req.Quantity,
		Rate:            req.Rate,
		GradeID:         req.GradeID,
		RefNumber:       req.RefNumber,
		TransactionDate: utils.ParseDate(req.TransactionDate),
		TransporterID:   req.TransporterID,
		Amount:          req.Amount,
	}

	if err := db.DB.Create(sale).Error; err != nil {
		return nil, err
	}
	return sale, nil
}

func (s *MilkLocalSaleService) GetMilkLocalSales(page, limit int) ([]dtos.MilkLocalSaleResponse, int64, error) {
	var results []dtos.MilkLocalSaleResponse
	var total int64
	db.DB.Model(&models.MilkLocalSale{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			mls.id, mls.quantity, mls.rate, mls.grade_id, pg.name AS grade_name,
			mls.ref_number, mls.transaction_date, mls.transporter_id, t.transporter_no AS transporter_name,
			mls.amount, mls.created_at, mls.updated_at
		FROM milk_local_sales mls
		LEFT JOIN product_grades pg ON mls.grade_id = pg.id
		LEFT JOIN transporters t ON mls.transporter_id = t.id
		WHERE mls.deleted_at IS NULL
		ORDER BY mls.transaction_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *MilkLocalSaleService) GetMilkLocalSale(id string) (*dtos.MilkLocalSaleResponse, error) {
	var result dtos.MilkLocalSaleResponse
	query := `
		SELECT 
			mls.id, mls.quantity, mls.rate, mls.grade_id, pg.name AS grade_name,
			mls.ref_number, mls.transaction_date, mls.transporter_id, t.transporter_no AS transporter_name,
			mls.amount, mls.created_at, mls.updated_at
		FROM milk_local_sales mls
		LEFT JOIN product_grades pg ON mls.grade_id = pg.id
		LEFT JOIN transporters t ON mls.transporter_id = t.id
		WHERE mls.id = ? AND mls.deleted_at IS NULL
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

func (s *MilkLocalSaleService) UpdateMilkLocalSale(id string, req dtos.UpdateMilkLocalSaleRequest) error {
	var sale models.MilkLocalSale
	if err := db.DB.First(&sale, id).Error; err != nil {
		return err
	}

	sale.Quantity = req.Quantity
	sale.Rate = req.Rate
	sale.GradeID = req.GradeID
	sale.RefNumber = req.RefNumber
	sale.TransactionDate = utils.ParseDate(req.TransactionDate)
	sale.TransporterID = req.TransporterID
	sale.Amount = req.Amount

	return db.DB.Save(&sale).Error
}

func (s *MilkLocalSaleService) DeleteMilkLocalSale(id string) error {
	return db.DB.Delete(&models.MilkLocalSale{}, id).Error
}
