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

type SharePaymentService struct{}

func NewSharePaymentService() *SharePaymentService {
	return &SharePaymentService{}
}

func (s *SharePaymentService) CreateSharePayment(req dtos.CreateSharePaymentRequest) (*models.SharePayment, error) {
	// Get posting rules for shares contribution
	var rule models.TransactionPostingRule
	if err := db.DB.Where("transaction_type = ?", "SHARES_CONTRIBUTION").First(&rule).Error; err != nil {
		return nil, fmt.Errorf("posting rule for SHARES_CONTRIBUTION not found: %w", err)
	}

	// Get share account to determine share type and unit price
	var shareAccount models.ShareAccount
	if err := db.DB.First(&shareAccount, req.ShareAccountID).Error; err != nil {
		return nil, fmt.Errorf("share account not found: %w", err)
	}

	// Get share type to determine unit price
	var shareType models.ShareType
	if err := db.DB.First(&shareType, shareAccount.ShareTypeID).Error; err != nil {
		return nil, fmt.Errorf("share type not found for share account: %w", err)
	}

	transactionDate := utils.ParseDate(req.TransactionDate)
	status := req.Status
	if status == "" {
		status = "POSTED"
	}

	payment := &models.SharePayment{
		MemberID:        req.MemberID,
		ShareAccountID:  req.ShareAccountID,
		AmountPaid:      req.AmountPaid,
		ShareUnits:      req.ShareUnits,
		PaymentModeID:   req.PaymentModeID,
		Description:     req.Description,
		Status:          status,
		TransactionDate: transactionDate,
		ApprovedBy:      req.ApprovedBy,
	}

	if req.DateApproved != "" {
		t := utils.ParseDate(req.DateApproved)
		payment.DateApproved = &t
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create Main Transaction Record
		transaction := &models.Transaction{
			Reference:       fmt.Sprintf("SHR-%s-%04d", transactionDate.Format("200601"), req.MemberID),
			TransactionName: "Share Contribution",
			TransactionType: "SHARE",
			TransactionDate: transactionDate,
			Description:     req.Description,
			Status:          status,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 2. Link Transaction ID to Share Payment and Save
		payment.TransactionID = transaction.ID
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// 3. Record Member Share Movement (PURCHASE)
		var prevBalance float64
		tx.Model(&models.ShareTransaction{}).
			Where("share_account_id = ?", payment.ShareAccountID).
			Select("COALESCE(SUM(debit - credit), 0)").Scan(&prevBalance)

		shareMovement := &models.ShareTransaction{
			TransactionID:   transaction.ID,
			ShareAccountID:  payment.ShareAccountID,
			MemberID:        payment.MemberID,
			TransactionType: "PURCHASE",
			ShareUnits:      payment.ShareUnits,
			UnitPrice:       shareType.ShareValue,
			Debit:           payment.AmountPaid,
			Credit:          0.00,
			BalanceAfter:    prevBalance + payment.AmountPaid,
			TransactionDate: transactionDate,
		}
		if err := tx.Create(shareMovement).Error; err != nil {
			return err
		}

		// 4. Create General Ledger Debit Entry (typically Bank or Cash)
		debitGL := &models.GeneralLedgerEntry{
			TransactionID:   transaction.ID,
			AccountID:       rule.DebitAccountID,
			SubAccountID:    rule.DebitSubAccountID,
			Debit:           req.AmountPaid,
			Credit:          0.00,
			TransactionDate: transactionDate,
			Description:     fmt.Sprintf("Share contribution - %s", req.Description),
		}
		if err := tx.Create(debitGL).Error; err != nil {
			return err
		}

		// 5. Create General Ledger Credit Entry (typically Share Capital)
		creditGL := &models.GeneralLedgerEntry{
			TransactionID:   transaction.ID,
			AccountID:       rule.CreditAccountID,
			SubAccountID:    rule.CreditSubAccountID,
			Debit:           0.00,
			Credit:          req.AmountPaid,
			TransactionDate: transactionDate,
			Description:     fmt.Sprintf("Share contribution by member %d - %s", req.MemberID, req.Description),
		}
		if err := tx.Create(creditGL).Error; err != nil {
			return err
		}

		// 6. Update Share Account Balance
		if err := tx.Model(&models.ShareAccount{}).Where("id = ?", payment.ShareAccountID).Updates(map[string]interface{}{
			"share_units":  gorm.Expr("share_units + ?", payment.ShareUnits),
			"share_amount": gorm.Expr("share_amount + ?", payment.AmountPaid),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *SharePaymentService) GetSharePayments() ([]dtos.SharePaymentResponse, int64, error) {
	var results []dtos.SharePaymentResponse
	var total int64
	db.DB.Model(&models.SharePayment{}).Count(&total)

	query := `
		SELECT 
			sp.id, sp.transaction_id, sp.member_id, m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			sp.share_account_id, sp.amount_paid, sp.share_units, sp.payment_mode_id, pm.name AS payment_mode_name,
			t.reference AS reference_no, sp.description, sp.status, sp.transaction_date, sp.approved_by, u.name AS approved_by_user_name,
			sp.date_approved, sp.created_at, sp.updated_at
		FROM share_payments sp
		LEFT JOIN member_registrations m ON sp.member_id = m.id
		LEFT JOIN payment_modes pm ON sp.payment_mode_id = pm.id
		LEFT JOIN users u ON sp.approved_by = u.id
		LEFT JOIN transactions t ON sp.transaction_id = t.id
		WHERE sp.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *SharePaymentService) GetSharePayment(id string) (*dtos.SharePaymentResponse, error) {
	var result dtos.SharePaymentResponse
	query := `
		SELECT 
			sp.id, sp.transaction_id, sp.member_id, m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			sp.share_account_id, sp.amount_paid, sp.share_units, sp.payment_mode_id, pm.name AS payment_mode_name,
			t.reference AS reference_no, sp.description, sp.status, sp.transaction_date, sp.approved_by, u.name AS approved_by_user_name,
			sp.date_approved, sp.created_at, sp.updated_at
		FROM share_payments sp
		LEFT JOIN member_registrations m ON sp.member_id = m.id
		LEFT JOIN payment_modes pm ON sp.payment_mode_id = pm.id
		LEFT JOIN users u ON sp.approved_by = u.id
		LEFT JOIN transactions t ON sp.transaction_id = t.id
		WHERE sp.id = ? AND sp.deleted_at IS NULL
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

func (s *SharePaymentService) UpdateSharePayment(id string, req dtos.UpdateSharePaymentRequest) error {
	var payment models.SharePayment
	if err := db.DB.First(&payment, id).Error; err != nil {
		return err
	}

	// Get posting rules for shares contribution
	var rule models.TransactionPostingRule
	if err := db.DB.Where("transaction_type = ?", "SHARES_CONTRIBUTION").First(&rule).Error; err != nil {
		return fmt.Errorf("posting rule for SHARES_CONTRIBUTION not found: %w", err)
	}

	// Get share account to determine share type and unit price
	var shareAccount models.ShareAccount
	if err := db.DB.First(&shareAccount, payment.ShareAccountID).Error; err != nil {
		return fmt.Errorf("share account not found: %w", err)
	}

	// Get share type to determine unit price
	var shareType models.ShareType
	if err := db.DB.First(&shareType, shareAccount.ShareTypeID).Error; err != nil {
		return fmt.Errorf("share type not found for share account: %w", err)
	}

	oldUnits := payment.ShareUnits
	oldAmount := payment.AmountPaid
	oldAccountID := payment.ShareAccountID

	transactionDate := utils.ParseDate(req.TransactionDate)

	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update Share Payment record
		payment.ShareAccountID = req.ShareAccountID
		payment.AmountPaid = req.AmountPaid
		payment.ShareUnits = req.ShareUnits
		payment.PaymentModeID = req.PaymentModeID
		payment.Description = req.Description
		payment.Status = req.Status
		payment.TransactionDate = transactionDate
		payment.ApprovedBy = req.ApprovedBy

		if req.DateApproved != "" {
			t := utils.ParseDate(req.DateApproved)
			payment.DateApproved = &t
		} else {
			payment.DateApproved = nil
		}

		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		// 2. Update linked Transaction record
		if err := tx.Model(&models.Transaction{}).Where("id = ?", payment.TransactionID).Updates(map[string]interface{}{
			"transaction_date": transactionDate,
			"description":      payment.Description,
			"status":           payment.Status,
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}

		// 3. Update Member Share Movement (PURCHASE)
		var prevBalance float64
		// Calculate previous balance for the share account, excluding the current transaction's effect
		tx.Model(&models.ShareTransaction{}).
			Where("share_account_id = ? AND transaction_id != ?", payment.ShareAccountID, payment.TransactionID).
			Select("COALESCE(SUM(debit - credit), 0)").Scan(&prevBalance)

		if err := tx.Model(&models.ShareTransaction{}).
			Where("transaction_id = ? AND transaction_type = 'PURCHASE'", payment.TransactionID).
			Updates(map[string]interface{}{
				"share_account_id": payment.ShareAccountID,
				"share_units":      req.ShareUnits,
				"unit_price":       shareType.ShareValue,
				"debit":            req.AmountPaid,
				"credit":           0.00,
				"balance_after":    prevBalance + req.AmountPaid,
				"transaction_date": transactionDate,
				"updated_at":       time.Now(),
			}).Error; err != nil {
			return err
		}

		// 4. Update General Ledger Debit Entry
		if err := tx.Model(&models.GeneralLedgerEntry{}).
			Where("transaction_id = ? AND debit > 0", payment.TransactionID).
			Updates(map[string]interface{}{
				"account_id":       rule.DebitAccountID,
				"sub_account_id":   rule.DebitSubAccountID,
				"debit":            payment.AmountPaid,
				"transaction_date": transactionDate,
				"description":      fmt.Sprintf("Share contribution by member %d (updated) - %s", payment.MemberID, req.Description),
				"updated_at":       time.Now(),
			}).Error; err != nil {
			return err
		}

		// 5. Update General Ledger Credit Entry
		if err := tx.Model(&models.GeneralLedgerEntry{}).
			Where("transaction_id = ? AND credit > 0", payment.TransactionID).
			Updates(map[string]interface{}{
				"account_id":       rule.CreditAccountID,
				"sub_account_id":   rule.CreditSubAccountID,
				"credit":           payment.AmountPaid,
				"transaction_date": transactionDate,
				"description":      fmt.Sprintf("Share contribution by member %d (updated) - %s", payment.MemberID, req.Description),
				"updated_at":       time.Now(),
			}).Error; err != nil {
			return err
		}

		// 6. Update Share Account Balance
		if oldAccountID == req.ShareAccountID {
			// Same account, just update difference
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", oldAccountID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units + ?", req.ShareUnits-oldUnits),
				"share_amount": gorm.Expr("share_amount + ?", req.AmountPaid-oldAmount),
			}).Error; err != nil {
				return err
			}
		} else {
			// Account changed
			// Subtract from old
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", oldAccountID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units - ?", oldUnits),
				"share_amount": gorm.Expr("share_amount - ?", oldAmount),
			}).Error; err != nil {
				return err
			}
			// Add to new
			if err := tx.Model(&models.ShareAccount{}).Where("id = ?", req.ShareAccountID).Updates(map[string]interface{}{
				"share_units":  gorm.Expr("share_units + ?", req.ShareUnits),
				"share_amount": gorm.Expr("share_amount + ?", req.AmountPaid),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *SharePaymentService) DeleteSharePayment(id string) error {
	return db.DB.Delete(&models.SharePayment{}, id).Error
}
