package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type LivestockProductionService struct{}

func NewLivestockProductionService() *LivestockProductionService {
	return &LivestockProductionService{}
}

func (s *LivestockProductionService) CreateProduction(req dtos.CreateLivestockProductionRequest, userID uint64) (*models.LivestockProductionRecord, error) {
	// Validate livestock status
	var livestock models.Livestock
	if err := db.DB.First(&livestock, req.LivestockID).Error; err != nil {
		return nil, fmt.Errorf("livestock not found")
	}

	if livestock.Status != "ACTIVE" {
		return nil, fmt.Errorf("cannot record production for livestock with status: %s", livestock.Status)
	}

	prod := &models.LivestockProductionRecord{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		LivestockID:    req.LivestockID,
		ProductionType: req.ProductionType,
		ProductionDate: utils.ParseDate(req.ProductionDate),
		Quantity:       req.Quantity,
		Unit:           req.Unit,
		Remarks:        req.Remarks,
	}

	if err := db.DB.Create(prod).Error; err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *LivestockProductionService) GetProductionRecords(livestockID string, page, limit int) ([]dtos.LivestockProductionResponse, int64, error) {
	var results []dtos.LivestockProductionResponse
	var total int64

	queryBuilder := db.DB.Model(&models.LivestockProductionRecord{})
	if livestockID != "" {
		queryBuilder = queryBuilder.Where("livestock_id = ?", livestockID)
	}
	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT lpr.*, l.tag_no as livestock_tag_no
		FROM livestock_production_records lpr
		JOIN livestocks l ON lpr.livestock_id = l.id
		WHERE (? = '' OR lpr.livestock_id = ?) AND lpr.deleted_at IS NULL
		ORDER BY lpr.production_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, livestockID, livestockID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *LivestockProductionService) GetProductionRecord(id string) (*dtos.LivestockProductionResponse, error) {
	var result dtos.LivestockProductionResponse
	query := `
		SELECT lpr.*, l.tag_no as livestock_tag_no
		FROM livestock_production_records lpr
		JOIN livestocks l ON lpr.livestock_id = l.id
		WHERE lpr.id = ? AND lpr.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil || result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *LivestockProductionService) UpdateProductionRecord(id string, req dtos.UpdateLivestockProductionRequest, userID uint64) error {
	var prod models.LivestockProductionRecord
	if err := db.DB.First(&prod, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"production_type": req.ProductionType,
		"production_date": utils.ParseDate(req.ProductionDate),
		"quantity":        req.Quantity,
		"unit":            req.Unit,
		"remarks":         req.Remarks,
		"updated_by":      userID,
	}
	return db.DB.Model(&prod).Updates(updates).Error
}

func (s *LivestockProductionService) DeleteProductionRecord(id string, userID uint64) error {
	return db.DB.Model(&models.LivestockProductionRecord{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.LivestockProductionRecord{}).Error
}

func (s *LivestockProductionService) GetProductionStats(livestockID string) (map[string]interface{}, error) {
	var stats struct {
		TotalQuantity float64 `json:"total_quantity"`
		AvgQuantity   float64 `json:"avg_quantity"`
		RecordCount   int64   `json:"record_count"`
	}

	err := db.DB.Model(&models.LivestockProductionRecord{}).
		Select("SUM(quantity) as total_quantity, AVG(quantity) as avg_quantity, COUNT(*) as record_count").
		Where("livestock_id = ? AND deleted_at IS NULL", livestockID).
		Scan(&stats).Error

	return map[string]interface{}{"stats": stats}, err
}
