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

type ShareTransferService struct{}

func NewShareTransferService() *ShareTransferService {
	return &ShareTransferService{}
}

func (s *ShareTransferService) CreateShareTransfer(req dtos.CreateShareTransferRequest) (*models.ShareTransfer, error) {
	// 1. Get posting rules for share transfers
	var rule models.TransactionPostingRule
	if err := db.DB.Where("transaction_type = ?", "SHARE_TRANSFERS").First(&rule).Error; err != nil {
		return nil, fmt.Errorf("posting rule for SHARE_TRANSFERS not found: %w", err)
	}

	// 2. Get share type to determine unit price and calculate transfer amount
	var shareType models.ShareType
	if err := db.DB.First(&shareType, req.ShareTypeID).Error; err != nil {
		return nil, fmt.Errorf("share type not found: %w", err)
	}

	// 3. Find share accounts for both members
	var fromAccount, toAccount models.ShareAccount
	if err := db.DB.Where("member_id = ? AND share_type_id = ?", req.FromMemberID, req.ShareTypeID).First(&fromAccount).Error; err != nil {
		return nil, fmt.Errorf("sender share account not found for this share type")
	}
	if err := db.DB.Where("member_id = ? AND share_type_id = ?", req.ToMemberID, req.ShareTypeID).First(&toAccount).Error; err != nil {
		return nil, fmt.Errorf("receiver share account not found for this share type")
	}

	transferAmount := req.ShareUnits * shareType.ShareValue
	transactionDate := utils.ParseDate(req.TransactionDate)
	status := req.Status
	if status == "" {
		status = "PENDING"
	}

	transfer := &models.ShareTransfer{
		FromMemberID:    req.FromMemberID,
		ToMemberID:      req.ToMemberID,
		ShareUnits:      req.ShareUnits,
		TransferAmount:  transferAmount,
		ShareTypeID:     req.ShareTypeID,
		Status:          status,
		TransactionDate: transactionDate,
		ApprovedBy:      req.ApprovedBy,
	}

	if req.DateApproved != "" {
		t := utils.ParseDate(req.DateApproved)
		transfer.DateApproved = &t
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 4. Create Main Transaction Record
		transaction := &models.Transaction{
			Reference:       fmt.Sprintf("SHR-TRF-%s-%03d", transactionDate.Format("20060102"), req.FromMemberID),
			TransactionName: "Share Transfer",
			TransactionType: "SHARE",
			TransactionDate: transactionDate,
			Description:     fmt.Sprintf("Transfer of %.2f shares from member %d to member %d", req.ShareUnits, req.FromMemberID, req.ToMemberID),
			Status:          "POSTED",
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 5. Link Transaction ID to Share Transfer and Save
		transfer.TransactionID = transaction.ID
		if err := tx.Create(transfer).Error; err != nil {
			return err
		}

		// 6. Record Member 101 Share Movement (TRANSFER_OUT)
		var fromPrevBalance float64
		tx.Model(&models.ShareTransaction{}).Where("share_account_id = ?", fromAccount.ID).Select("COALESCE(SUM(debit - credit), 0)").Scan(&fromPrevBalance)

		senderMovement := &models.ShareTransaction{
			TransactionID:   transaction.ID,
			ShareAccountID:  fromAccount.ID,
			MemberID:        req.FromMemberID,
			TransactionType: "TRANSFER_OUT",
			ShareUnits:      req.ShareUnits,
			UnitPrice:       shareType.ShareValue,
			Debit:           0.00,
			Credit:          transferAmount,
			BalanceAfter:    fromPrevBalance - transferAmount,
			TransactionDate: transactionDate,
		}
		if err := tx.Create(senderMovement).Error; err != nil {
			return err
		}

		// 7. Record Member 102 Share Movement (TRANSFER_IN)
		var toPrevBalance float64
		tx.Model(&models.ShareTransaction{}).Where("share_account_id = ?", toAccount.ID).Select("COALESCE(SUM(debit - credit), 0)").Scan(&toPrevBalance)

		receiverMovement := &models.ShareTransaction{
			TransactionID:   transaction.ID,
			ShareAccountID:  toAccount.ID,
			MemberID:        req.ToMemberID,
			TransactionType: "TRANSFER_IN",
			ShareUnits:      req.ShareUnits,
			UnitPrice:       shareType.ShareValue,
			Debit:           transferAmount,
			Credit:          0.00,
			BalanceAfter:    toPrevBalance + transferAmount,
			TransactionDate: transactionDate,
		}
		if err := tx.Create(receiverMovement).Error; err != nil {
			return err
		}

		// 8. General Ledger Entries
		// Debit receiving member share ledger
		debitGL := &models.GeneralLedgerEntry{
			TransactionID:   transaction.ID,
			AccountID:       rule.DebitAccountID,
			SubAccountID:    rule.DebitSubAccountID,
			Debit:           transferAmount,
			Credit:          0.00,
			TransactionDate: time.Now(),
			Description:     fmt.Sprintf("Share transfer received by member %d", req.ToMemberID),
		}
		if err := tx.Create(debitGL).Error; err != nil {
			return err
		}

		// Credit sending member share ledger
		creditGL := &models.GeneralLedgerEntry{
			TransactionID:   transaction.ID,
			AccountID:       rule.CreditAccountID,
			SubAccountID:    rule.CreditSubAccountID,
			Debit:           0.00,
			Credit:          transferAmount,
			TransactionDate: time.Now(),
			Description:     fmt.Sprintf("Share transfer issued by member %d", req.FromMemberID),
		}
		if err := tx.Create(creditGL).Error; err != nil {
			return err
		}

		// 9. Update Member Share Account Balances
		// Sender
		if err := tx.Model(&models.ShareAccount{}).Where("id = ?", fromAccount.ID).Updates(map[string]interface{}{
			"share_units":  gorm.Expr("share_units - ?", req.ShareUnits),
			"share_amount": gorm.Expr("share_amount - ?", transferAmount),
		}).Error; err != nil {
			return err
		}
		// Receiver
		if err := tx.Model(&models.ShareAccount{}).Where("id = ?", toAccount.ID).Updates(map[string]interface{}{
			"share_units":  gorm.Expr("share_units + ?", req.ShareUnits),
			"share_amount": gorm.Expr("share_amount + ?", transferAmount),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *ShareTransferService) GetShareTransfers() ([]dtos.ShareTransferResponse, int64, error) {
	var results []dtos.ShareTransferResponse
	var total int64
	db.DB.Model(&models.ShareTransfer{}).Count(&total)

	query := `
		SELECT 
			st.id, st.transaction_id,
			st.from_member_id, fm.member_no AS from_member_no, fm.first_name AS from_member_first_name, fm.last_name AS from_member_last_name, 
			st.to_member_id, tm.member_no AS to_member_no, tm.first_name AS to_member_first_name, tm.last_name AS to_member_last_name, 
			st.share_units, st.transfer_amount, t.reference AS reference_no, st.status, st.transaction_date,
			st.approved_by, u.name AS approved_by_user_name, st.date_approved,
			st.created_at, st.updated_at
		FROM share_transfers st
		LEFT JOIN member_registrations fm ON st.from_member_id = fm.id
		LEFT JOIN member_registrations tm ON st.to_member_id = tm.id
		LEFT JOIN users u ON st.approved_by = u.id
		LEFT JOIN transactions t ON st.transaction_id = t.id
		WHERE st.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ShareTransferService) GetShareTransfer(id string) (*dtos.ShareTransferResponse, error) {
	var result dtos.ShareTransferResponse
	query := `
		SELECT 
			st.id, st.transaction_id, 
			st.from_member_id, fm.member_no AS from_member_no, fm.first_name AS from_member_first_name, fm.last_name AS from_member_last_name, 
			st.to_member_id, tm.member_no AS to_member_no, tm.first_name AS to_member_first_name, tm.last_name AS to_member_last_name, 
			st.share_units, st.transfer_amount, t.reference AS reference_no, st.status, st.transaction_date,
			st.approved_by, u.name AS approved_by_user_name, st.date_approved,
			st.created_at, st.updated_at
		FROM share_transfers st
		LEFT JOIN member_registrations fm ON st.from_member_id = fm.id
		LEFT JOIN member_registrations tm ON st.to_member_id = tm.id
		LEFT JOIN users u ON st.approved_by = u.id
		LEFT JOIN transactions t ON st.transaction_id = t.id
		WHERE st.id = ? AND st.deleted_at IS NULL
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

func (s *ShareTransferService) UpdateShareTransfer(id string, req dtos.UpdateShareTransferRequest) error {
	var transfer models.ShareTransfer
	if err := db.DB.First(&transfer, id).Error; err != nil {
		return err
	}

	// 1. Get posting rules for share transfers
	var rule models.TransactionPostingRule
	if err := db.DB.Where("transaction_type = ?", "SHARE_TRANSFERS").First(&rule).Error; err != nil {
		return fmt.Errorf("posting rule for SHARE_TRANSFERS not found: %w", err)
	}

	// 2. Get share type to determine current unit price
	var shareType models.ShareType
	if err := db.DB.First(&shareType, req.ShareTypeID).Error; err != nil {
		return fmt.Errorf("share type not found: %w", err)
	}

	// 3. Find share accounts for the existing members in this transfer
	var fromAccount, toAccount models.ShareAccount
	if err := db.DB.Where("member_id = ? AND share_type_id = ?", transfer.FromMemberID, req.ShareTypeID).First(&fromAccount).Error; err != nil {
		return fmt.Errorf("sender share account not found for this share type")
	}
	if err := db.DB.Where("member_id = ? AND share_type_id = ?", transfer.ToMemberID, req.ShareTypeID).First(&toAccount).Error; err != nil {
		return fmt.Errorf("receiver share account not found for this share type")
	}

	oldUnits := transfer.ShareUnits
	oldAmount := transfer.TransferAmount
	oldShareTypeID := transfer.ShareTypeID

	transferAmount := req.ShareUnits * shareType.ShareValue
	transactionDate := utils.ParseDate(req.TransactionDate)

	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 4. Update Share Transfer Header
		transfer.ShareUnits = req.ShareUnits
		transfer.TransferAmount = transferAmount
		transfer.ShareTypeID = req.ShareTypeID
		transfer.Status = req.Status
		transfer.TransactionDate = transactionDate
		transfer.ApprovedBy = req.ApprovedBy

		if req.DateApproved != "" {
			t := utils.ParseDate(req.DateApproved)
			transfer.DateApproved = &t
		} else {
			transfer.DateApproved = nil
		}

		if err := tx.Save(&transfer).Error; err != nil {
			return err
		}

		// 5. Update linked Main Transaction Record
		if err := tx.Model(&models.Transaction{}).Where("id = ?", transfer.TransactionID).Updates(map[string]interface{}{
			"transaction_date": transactionDate,
			"description":      fmt.Sprintf("Transfer of %.2f shares from member %d to member %d (updated)", req.ShareUnits, transfer.FromMemberID, transfer.ToMemberID),
			"status":           "POSTED",
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}

		// 6. Update Member Share Movements
		// Sender Movement (TRANSFER_OUT)
		var fromPrevBalance float64
		tx.Model(&models.ShareTransaction{}).
			Where("share_account_id = ? AND transaction_id != ?", fromAccount.ID, transfer.TransactionID).
			Select("COALESCE(SUM(debit - credit), 0)").Scan(&fromPrevBalance)

		if err := tx.Model(&models.ShareTransaction{}).
			Where("transaction_id = ? AND transaction_type = 'TRANSFER_OUT'", transfer.TransactionID).
			Updates(map[string]interface{}{
				"share_units":      req.ShareUnits,
				"unit_price":       shareType.ShareValue,
				"credit":           transferAmount,
				"balance_after":    fromPrevBalance - transferAmount,
				"transaction_date": transactionDate,
				"updated_at":       time.Now(),
			}).Error; err != nil {
			return err
		}

		// Receiver Movement (TRANSFER_IN)
		var toPrevBalance float64
		tx.Model(&models.ShareTransaction{}).
			Where("share_account_id = ? AND transaction_id != ?", toAccount.ID, transfer.TransactionID).
			Select("COALESCE(SUM(debit - credit), 0)").Scan(&toPrevBalance)

		if err := tx.Model(&models.ShareTransaction{}).
			Where("transaction_id = ? AND transaction_type = 'TRANSFER_IN'", transfer.TransactionID).
			Updates(map[string]interface{}{
				"share_units":      req.ShareUnits,
				"unit_price":       shareType.ShareValue,
				"debit":            transferAmount,
				"balance_after":    toPrevBalance + transferAmount,
				"transaction_date": transactionDate,
				"updated_at":       time.Now(),
			}).Error; err != nil {
			return err
		}

		// 7. Update General Ledger Entries
		// Update Debit GL entry (Receiver)
		tx.Model(&models.GeneralLedgerEntry{}).Where("transaction_id = ? AND debit > 0", transfer.TransactionID).Updates(map[string]interface{}{
			"debit": transferAmount, "description": fmt.Sprintf("Share transfer received by member %d (updated)", transfer.ToMemberID), "transaction_date": transactionDate, "updated_at": time.Now(),
		})

		// Update Credit GL entry (Sender)
		tx.Model(&models.GeneralLedgerEntry{}).Where("transaction_id = ? AND credit > 0", transfer.TransactionID).Updates(map[string]interface{}{
			"credit": transferAmount, "description": fmt.Sprintf("Share transfer issued by member %d (updated)", transfer.FromMemberID), "transaction_date": transactionDate, "updated_at": time.Now(),
		})

		// 8. Update Member Share Account Balances
		if oldShareTypeID == req.ShareTypeID {
			// Same accounts, update differences
			diffUnits := req.ShareUnits - oldUnits
			diffAmount := transferAmount - oldAmount

			// Sender: loses more if diff is positive
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", fromAccount.ID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units - ?", diffUnits),
				"share_amount": gorm.Expr("share_amount - ?", diffAmount),
			}).Error; err != nil {
				return err
			}
			// Receiver: gains more if diff is positive
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", toAccount.ID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units + ?", diffUnits),
				"share_amount": gorm.Expr("share_amount + ?", diffAmount),
			}).Error; err != nil {
				return err
			}
		} else {
			// Share type changed, accounts changed
			// Revert old accounts
			var oldFromAcc, oldToAcc models.ShareAccount
			if err := tx.Where("member_id = ? AND share_type_id = ?", transfer.FromMemberID, oldShareTypeID).First(&oldFromAcc).Error; err != nil {
				return err
			}
			if err := tx.Where("member_id = ? AND share_type_id = ?", transfer.ToMemberID, oldShareTypeID).First(&oldToAcc).Error; err != nil {
				return err
			}

			// Sender revert
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", oldFromAcc.ID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units + ?", oldUnits),
				"share_amount": gorm.Expr("share_amount + ?", oldAmount),
			}).Error; err != nil {
				return err
			}
			// Receiver revert
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", oldToAcc.ID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units - ?", oldUnits),
				"share_amount": gorm.Expr("share_amount - ?", oldAmount),
			}).Error; err != nil {
				return err
			}

			// Apply new accounts (fromAccount and toAccount were fetched based on req.ShareTypeID)
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", fromAccount.ID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units - ?", req.ShareUnits),
				"share_amount": gorm.Expr("share_amount - ?", transferAmount),
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", toAccount.ID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units + ?", req.ShareUnits),
				"share_amount": gorm.Expr("share_amount + ?", transferAmount),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ShareTransferService) DeleteShareTransfer(id string) error {
	return db.DB.Delete(&models.ShareTransfer{}, id).Error
}
