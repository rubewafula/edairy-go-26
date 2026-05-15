package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type OrganizationWalletService struct{}

func NewOrganizationWalletService() *OrganizationWalletService {
	return &OrganizationWalletService{}
}

func (s *OrganizationWalletService) CreateWallet(req dtos.CreateOrganizationWalletRequest, userID uint64) (*models.OrganizationWallet, error) {
	wallet := &models.OrganizationWallet{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		WalletTypeID: req.WalletTypeID,
		WalletID:     req.WalletID,
		WalletName:   req.WalletName,
	}
	if err := db.DB.Create(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *OrganizationWalletService) GetWallets(page, limit int) ([]dtos.OrganizationWalletResponse, int64, error) {
	var results []dtos.OrganizationWalletResponse
	var total int64
	db.DB.Model(&models.OrganizationWallet{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT ow.*
		FROM organization_wallets ow
		WHERE ow.deleted_at IS NULL
		ORDER BY ow.id DESC LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *OrganizationWalletService) GetWallet(id string) (*dtos.OrganizationWalletResponse, error) {
	var result dtos.OrganizationWalletResponse
	query := `
		SELECT ow.*
		FROM organization_wallets ow
		WHERE ow.id = ? AND ow.deleted_at IS NULL LIMIT 1
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

func (s *OrganizationWalletService) UpdateWallet(id string, req dtos.UpdateOrganizationWalletRequest, userID uint64) error {
	var wallet models.OrganizationWallet
	if err := db.DB.First(&wallet, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"wallet_type_id": req.WalletTypeID,
		"wallet_id":      req.WalletID,
		"wallet_name":    req.WalletName,
		"updated_by":     userID,
	}

	return db.DB.Model(&wallet).Updates(updates).Error
}

func (s *OrganizationWalletService) DeleteWallet(id string, userID uint64) error {
	var wallet models.OrganizationWallet
	if err := db.DB.First(&wallet, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&wallet).Update("updated_by", userID).Delete(&wallet).Error
}
