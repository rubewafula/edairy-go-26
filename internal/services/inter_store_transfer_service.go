package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type InterStoreTransferService struct{}

func NewInterStoreTransferService() *InterStoreTransferService {
	return &InterStoreTransferService{}
}

func (s *InterStoreTransferService) CreateTransfer(req dtos.CreateInterStoreTransferRequest, userID uint64) (*models.InterStoreTransfer, error) {
	transfer := &models.InterStoreTransfer{
		BaseModel: models.BaseModel{
			CreatedBy: userID,
		},
		FromStoreID:  req.FromStoreID,
		ToStoreID:    req.ToStoreID,
		Reference:    req.Reference,
		TransferDate: utils.ParseDate(req.TransferDate),
		Status:       req.Status,
	}

	if err := db.DB.Create(transfer).Error; err != nil {
		return nil, err
	}
	return transfer, nil
}

func (s *InterStoreTransferService) GetTransfers(page, limit int) ([]dtos.InterStoreTransferResponse, int64, error) {
	var results []dtos.InterStoreTransferResponse
	var total int64
	db.DB.Model(&models.InterStoreTransfer{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			ist.*, fs.store AS from_store_name, ts.store AS to_store_name, ist.status
		FROM inter_store_transfers ist
		LEFT JOIN stores fs ON ist.from_store_id = fs.id
		LEFT JOIN stores ts ON ist.to_store_id = ts.id
		WHERE ist.deleted_at IS NULL
		ORDER BY ist.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *InterStoreTransferService) GetTransfer(id string) (*dtos.InterStoreTransferResponse, error) {
	var result dtos.InterStoreTransferResponse
	query := `
		SELECT 
			ist.*, fs.store AS from_store_name, ts.store AS to_store_name, ist.status
		FROM inter_store_transfers ist
		LEFT JOIN stores fs ON ist.from_store_id = fs.id
		LEFT JOIN stores ts ON ist.to_store_id = ts.id
		WHERE ist.id = ? AND ist.deleted_at IS NULL
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

func (s *InterStoreTransferService) UpdateTransfer(id string, req dtos.UpdateInterStoreTransferRequest, userID uint64) error {
	var transfer models.InterStoreTransfer
	if err := db.DB.First(&transfer, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&transfer).Updates(map[string]interface{}{
		"from_store_id": req.FromStoreID,
		"to_store_id":   req.ToStoreID,
		"reference":     req.Reference,
		"transfer_date": utils.ParseDate(req.TransferDate),
		"status":        req.Status,
		"updated_by":    userID,
	}).Error
}

func (s *InterStoreTransferService) DeleteTransfer(id string) error {
	var transfer models.InterStoreTransfer
	if err := db.DB.First(&transfer, id).Error; err != nil {
		return err
	}
	return db.DB.Delete(&transfer).Error
}
