package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type AssetDepreciationService struct{}

func NewAssetDepreciationService() *AssetDepreciationService {
	return &AssetDepreciationService{}
}

func (s *AssetDepreciationService) CreateEntry(req dtos.CreateAssetDepreciationRequest) (*models.AssetDepreciationEntry, error) {
	// Get posting rules for depreciation
	var rule models.TransactionPostingRule
	// Fetch the asset to get its details for depreciation calculation
	var asset models.Asset
	if err := db.DB.First(&asset, req.AssetID).Error; err != nil {
		return nil, fmt.Errorf("asset not found: %w", err)
	}

	if err := db.DB.Where("transaction_type = ?", "DEPRECIATION_ACCRUAL").First(&rule).Error; err != nil {
		return nil, err
	}

	depreciationDate := utils.ParseDate(req.DepreciationDate)

	// Calculate new accumulated depreciation and book value
	newAccumulatedDepreciation := asset.AccumulatedDepreciation + req.DepreciationAmount
	newBookValue := asset.PurchaseCost - newAccumulatedDepreciation

	// Ensure book value does not go below salvage value
	if newBookValue < asset.SalvageValue {
		// Adjust depreciation amount to not go below salvage value
		adjustedDepreciationAmount := asset.BookValue - asset.SalvageValue
		if adjustedDepreciationAmount < 0 { // Already below salvage value, no more depreciation
			adjustedDepreciationAmount = 0
		}
		req.DepreciationAmount = adjustedDepreciationAmount
		newAccumulatedDepreciation = asset.AccumulatedDepreciation + adjustedDepreciationAmount
		newBookValue = asset.SalvageValue
	}

	// Determine the depreciation rate for the asset update
	// This logic is duplicated, consider extracting to a helper function if used frequently
	// or if more complex depreciation methods are introduced.
	// For now, it's kept here for clarity of the change.
	var depreciationRateForAssetUpdate float64

	if asset.DepreciationMethod == "STRAIGHT_LINE" && asset.UsefulLifeYears > 0 {
		// For straight-line, the annual rate is typically 1 / UsefulLifeYears
		depreciationRateForAssetUpdate = 1.0 / float64(asset.UsefulLifeYears)
	} else {

		depreciationRateForAssetUpdate = newAccumulatedDepreciation / asset.PurchaseCost

	}

	entry := &models.AssetDepreciationEntry{
		AssetID:                 req.AssetID,
		DepreciationDate:        depreciationDate,
		DepreciationAmount:      req.DepreciationAmount,
		AccumulatedDepreciation: newAccumulatedDepreciation,
		BookValue:               newBookValue,
		Notes:                   req.Notes,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create Main Transaction Record
		transaction := &models.Transaction{
			Reference:       fmt.Sprintf("DEP-%s-%04d", depreciationDate.Format("200601"), req.AssetID),
			TransactionName: "Asset Depreciation",
			TransactionType: "ASSETS",
			TransactionDate: depreciationDate,
			Description:     "Asset depreciation posting",
			Status:          "POSTED",
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 2. Link Transaction ID to Depreciation Entry and Save
		entry.TransactionID = transaction.ID
		if err := tx.Create(entry).Error; err != nil {
			return err
		}

		// 3. Create General Ledger Debit Entry
		debitGL := &models.GeneralLedgerEntry{
			TransactionID:   transaction.ID,
			AccountID:       rule.DebitAccountID,
			SubAccountID:    rule.DebitSubAccountID,
			Debit:           req.DepreciationAmount,
			Credit:          0.00,
			TransactionDate: time.Now(),
			Description:     "Monthly depreciation expense",
		}
		if err := tx.Create(debitGL).Error; err != nil {
			return err
		}

		// 4. Create General Ledger Credit Entry
		creditGL := &models.GeneralLedgerEntry{
			TransactionID:   transaction.ID,
			AccountID:       rule.CreditAccountID,
			SubAccountID:    rule.CreditSubAccountID,
			Debit:           0.00,
			Credit:          req.DepreciationAmount,
			TransactionDate: time.Now(),
			Description:     "Accumulated depreciation posting",
		}
		if err := tx.Create(creditGL).Error; err != nil {
			return err
		}

		// 5. Update Asset Master Record
		if err := tx.Model(&models.Asset{}).Where("id = ?", req.AssetID).Updates(map[string]interface{}{

			"accumulated_depreciation": newAccumulatedDepreciation,
			"book_value":               newBookValue,
			"depreciation_rate":        depreciationRateForAssetUpdate,
			"updated_at":               time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *AssetDepreciationService) UpdateEntry(id string, req dtos.UpdateAssetDepreciationRequest) (*models.AssetDepreciationEntry, error) {
	var oldEntry models.AssetDepreciationEntry
	if err := db.DB.First(&oldEntry, id).Error; err != nil {
		return nil, fmt.Errorf("asset depreciation entry not found: %w", err)
	}

	var asset models.Asset
	if err := db.DB.First(&asset, oldEntry.AssetID).Error; err != nil {
		return nil, fmt.Errorf("asset not found: %w", err)
	}

	var rule models.TransactionPostingRule
	if err := db.DB.Where("transaction_type = ?", "DEPRECIATION_ACCRUAL").First(&rule).Error; err != nil {
		return nil, err
	}

	depreciationDate := utils.ParseDate(req.DepreciationDate)

	// Calculate the change in depreciation amount
	changeInDepreciation := req.DepreciationAmount - oldEntry.DepreciationAmount

	// Calculate new accumulated depreciation and book value for the asset
	newAssetAccumulatedDepreciation := asset.AccumulatedDepreciation + changeInDepreciation
	newAssetBookValue := asset.PurchaseCost - newAssetAccumulatedDepreciation

	// Ensure book value does not go below salvage value
	if newAssetBookValue < asset.SalvageValue {
		// Adjust depreciation amount to not go below salvage value
		adjustedDepreciationAmount := asset.BookValue - asset.SalvageValue
		if adjustedDepreciationAmount < 0 { // Already below salvage value, no more depreciation
			adjustedDepreciationAmount = 0
		}
		req.DepreciationAmount = adjustedDepreciationAmount
		newAssetAccumulatedDepreciation = asset.AccumulatedDepreciation + (adjustedDepreciationAmount - oldEntry.DepreciationAmount)
		newAssetBookValue = asset.SalvageValue
	}

	var depreciationRateForAssetUpdate float64
	if asset.DepreciationMethod == "STRAIGHT_LINE" && asset.UsefulLifeYears > 0 {
		depreciationRateForAssetUpdate = 1.0 / float64(asset.UsefulLifeYears)
	} else {
		depreciationRateForAssetUpdate = newAssetAccumulatedDepreciation / asset.PurchaseCost
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update AssetDepreciationEntry
		if err := tx.Model(&oldEntry).Updates(map[string]interface{}{
			"depreciation_date":        depreciationDate,
			"depreciation_amount":      req.DepreciationAmount,
			"accumulated_depreciation": newAssetAccumulatedDepreciation,
			"book_value":               newAssetBookValue,
			"notes":                    req.Notes,
			"updated_at":               time.Now(),
		}).Error; err != nil {
			return err
		}

		// 2. Update Main Transaction Record
		if err := tx.Model(&models.Transaction{}).Where("id = ?", oldEntry.TransactionID).Updates(map[string]interface{}{
			"transaction_date": depreciationDate,
			"description":      req.Notes, // Assuming notes can update transaction description
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}

		// 3. Update General Ledger Debit Entry
		if err := tx.Model(&models.GeneralLedgerEntry{}).Where("transaction_id = ? AND debit > 0", oldEntry.TransactionID).Updates(map[string]interface{}{
			"debit":            req.DepreciationAmount,
			"transaction_date": time.Now(),
			"description":      "Asset depreciation expense (updated)",
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}

		// 4. Update General Ledger Credit Entry
		if err := tx.Model(&models.GeneralLedgerEntry{}).Where("transaction_id = ? AND credit > 0", oldEntry.TransactionID).Updates(map[string]interface{}{
			"credit":           req.DepreciationAmount,
			"transaction_date": time.Now(),
			"description":      "Accumulated depreciation posting (updated)",
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}

		// 5. Update Asset Master Record
		if err := tx.Model(&models.Asset{}).Where("id = ?", oldEntry.AssetID).Updates(map[string]interface{}{

			"accumulated_depreciation": newAssetAccumulatedDepreciation,
			"book_value":               newAssetBookValue,
			"depreciation_rate":        depreciationRateForAssetUpdate,
			"updated_at":               time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Re-fetch the updated entry to return
	var updatedEntry models.AssetDepreciationEntry
	// Using First with the ID ensures we get the latest state after the transaction
	if err := db.DB.First(&updatedEntry, id).Error; err != nil {
		return nil, err
	}

	return &updatedEntry, nil

}

func (s *AssetDepreciationService) GetEntries() ([]dtos.AssetDepreciationResponse, int64, error) {
	var results []dtos.AssetDepreciationResponse

	var total int64
	db.DB.Model(&models.AssetDepreciationEntry{}).Count(&total)

	query := `
		SELECT 
			ade.id, fa.id AS asset_id, fa.asset_name, fa.asset_code,
			ade.depreciation_date, ade.depreciation_amount, ade.accumulated_depreciation,
			ade.book_value, ade.transaction_id, ade.created_at
		FROM asset_depreciation_entries ade
		LEFT JOIN fixed_assets fa ON ade.asset_id = fa.id
		WHERE ade.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *AssetDepreciationService) GetEntry(id string) (*dtos.AssetDepreciationResponse, error) {
	var result dtos.AssetDepreciationResponse
	query := `
		SELECT 
			ade.id, ade.asset_id, fa.asset_name, fa.asset_code,
			ade.depreciation_date, ade.depreciation_amount, ade.accumulated_depreciation,
			ade.book_value, ade.transaction_id, ade.created_at
		FROM asset_depreciation_entries ade
		LEFT JOIN fixed_assets fa ON ade.asset_id = fa.id
		WHERE ade.id = ? AND ade.deleted_at IS NULL
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

func (s *AssetDepreciationService) DeleteEntry(id string) error {
	var entryToDelete models.AssetDepreciationEntry

	if err := db.DB.First(&entryToDelete, id).Error; err != nil {
		return fmt.Errorf("asset depreciation entry not found: %w", err)
	}

	var asset models.Asset
	if err := db.DB.First(&asset, entryToDelete.AssetID).Error; err != nil {

		return fmt.Errorf("asset not found: %w", err)
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Reverse the depreciation on the Asset master record
		newAssetAccumulatedDepreciation := asset.AccumulatedDepreciation - entryToDelete.DepreciationAmount
		newAssetBookValue := asset.BookValue + entryToDelete.DepreciationAmount

		// Ensure accumulated depreciation doesn't go below zero
		if newAssetAccumulatedDepreciation < 0 {
			newAssetAccumulatedDepreciation = 0
		}
		// Ensure book value doesn't exceed purchase cost
		if newAssetBookValue > asset.PurchaseCost {
			newAssetBookValue = asset.PurchaseCost
		}

		var depreciationRateForAssetUpdate float64
		if asset.DepreciationMethod == "STRAIGHT_LINE" && asset.UsefulLifeYears > 0 {
			depreciationRateForAssetUpdate = 1.0 / float64(asset.UsefulLifeYears)
		} else if asset.PurchaseCost > 0 {
			// Recalculate rate based on new accumulated depreciation
			depreciationRateForAssetUpdate = newAssetAccumulatedDepreciation / asset.PurchaseCost
		} else {
			depreciationRateForAssetUpdate = asset.DepreciationRate // Keep existing rate if no other calculation is possible
		}

		if err := tx.Model(&models.Asset{}).Where("id = ?", entryToDelete.AssetID).Updates(map[string]interface{}{
			"accumulated_depreciation": newAssetAccumulatedDepreciation,
			"book_value":               newAssetBookValue,
			"depreciation_rate":        depreciationRateForAssetUpdate,
			"updated_at":               time.Now(),
		}).Error; err != nil {
			return err
		}

		// 2. Soft delete General Ledger entries
		if err := tx.Where("transaction_id = ?", entryToDelete.TransactionID).Delete(&models.GeneralLedgerEntry{}).Error; err != nil {
			return err
		}

		// 3. Soft delete AssetDepreciationEntry
		if err := tx.Delete(&entryToDelete).Error; err != nil {
			return err
		}

		// 4. Soft delete Transaction record
		if err := tx.Delete(&models.Transaction{}, entryToDelete.TransactionID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
