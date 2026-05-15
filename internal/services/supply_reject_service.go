package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type SupplyRejectService struct{}

func NewSupplyRejectService() *SupplyRejectService {
	return &SupplyRejectService{}
}

func (s *SupplyRejectService) CreateReject(req dtos.CreateSupplyRejectRequest, userID uint64) (*models.SupplyReject, error) {
	reject := &models.SupplyReject{
		BaseModel: models.BaseModel{CreatedBy: userID},
		ItemID:    req.ItemID,
		SupplyID:  req.SupplyID,
		Quantity:  req.Quantity,
		Reason:    req.Reason,
	}
	err := db.DB.Create(reject).Error
	return reject, err
}

func (s *SupplyRejectService) GetRejects(page, limit int) ([]dtos.SupplyRejectResponse, int64, error) {
	var results []dtos.SupplyRejectResponse
	var total int64
	db.DB.Model(&models.SupplyReject{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT sr.*, si.item_name, 
		       CASE WHEN v.company_name != '' THEN v.company_name ELSE CONCAT(v.first_name, ' ', v.last_name) END as vendor_name
		FROM supply_rejects sr
		LEFT JOIN store_items si ON sr.item_id = si.id
		LEFT JOIN supplies s ON sr.supply_id = s.id
		LEFT JOIN suppliers v ON s.vendor_id = v.id
		WHERE sr.deleted_at IS NULL
		ORDER BY sr.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *SupplyRejectService) GetRejectsBySupply(supplyID string) ([]dtos.SupplyRejectResponse, error) {
	var results []dtos.SupplyRejectResponse
	query := `
		SELECT sr.*, si.item_name, 
		       CASE WHEN v.company_name != '' THEN v.company_name ELSE CONCAT(v.first_name, ' ', v.last_name) END as vendor_name
		FROM supply_rejects sr
		LEFT JOIN store_items si ON sr.item_id = si.id
		LEFT JOIN supplies s ON sr.supply_id = s.id
		LEFT JOIN suppliers v ON s.vendor_id = v.id
		WHERE sr.supply_id = ? AND sr.deleted_at IS NULL
	`
	err := db.DB.Raw(query, supplyID).Scan(&results).Error
	return results, err
}

func (s *SupplyRejectService) GetReject(id string) (*dtos.SupplyRejectResponse, error) {
	var result dtos.SupplyRejectResponse
	query := `
		SELECT sr.*, si.item_name, 
		       CASE WHEN v.company_name != '' THEN v.company_name ELSE CONCAT(v.first_name, ' ', v.last_name) END as vendor_name
		FROM supply_rejects sr
		LEFT JOIN store_items si ON sr.item_id = si.id
		LEFT JOIN supplies s ON sr.supply_id = s.id
		LEFT JOIN suppliers v ON s.vendor_id = v.id
		WHERE sr.id = ? AND sr.deleted_at IS NULL
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

func (s *SupplyRejectService) UpdateReject(id string, req dtos.UpdateSupplyRejectRequest, userID uint64) error {
	var reject models.SupplyReject
	if err := db.DB.First(&reject, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"quantity":   req.Quantity,
		"reason":     req.Reason,
		"updated_by": userID,
	}

	return db.DB.Model(&reject).Updates(updates).Error
}

func (s *SupplyRejectService) DeleteReject(id string, userID uint64) error {
	var reject models.SupplyReject
	if err := db.DB.First(&reject, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&reject).Update("updated_by", userID).Delete(&reject).Error
}
