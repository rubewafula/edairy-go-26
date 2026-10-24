package services

import (
	"log"
	"net/http"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// DefaultDBUMaxAmount is the DBU allocation maximum a loan product gets unless one is set (KES).
const DefaultDBUMaxAmount = 500000.0

// validateDBUQuota checks a product's DBU allocation quota (KES).
func validateDBUQuota(minAmount, maxAmount float64) error {
	if minAmount < 0 || maxAmount < 0 {
		return serviceErr(http.StatusBadRequest, "DBU minimum and maximum cannot be negative")
	}
	if maxAmount < minAmount {
		return serviceErr(http.StatusBadRequest, "DBU maximum cannot be less than the DBU minimum")
	}
	return nil
}

// EnsureLoanProductDBUColumns adds dbu_min_amount / dbu_max_amount to loan_products when missing.
// Existing products get minimum 0 and maximum 500,000 (column defaults).
func EnsureLoanProductDBUColumns() {
	m := db.DB.Migrator()
	for _, f := range []struct{ col, field string }{
		{"dbu_min_amount", "DBUMinAmount"},
		{"dbu_max_amount", "DBUMaxAmount"},
		{"processing_fee_amount", "ProcessingFeeAmount"},
	} {
		if !m.HasColumn(&models.LoanProduct{}, f.col) {
			if err := m.AddColumn(&models.LoanProduct{}, f.field); err != nil {
				log.Printf("[LoanProduct] add column %s: %v", f.col, err)
			}
		}
	}
}

// FixLoanRecurrentAmounts sets recurrent_amount = total_amount on open loan milk deductions that still
// carry the old system default (the contract's monthly installment, or total / term). Amounts that
// were changed by hand (anything else) are left alone.
func FixLoanRecurrentAmounts() {
	res := db.DB.Exec(`UPDATE recurrent_deductions rd
		JOIN loan_contracts c ON rd.reference = CONCAT('LOAN-', c.contract_no) AND rd.customer_id = c.member_id
		SET rd.recurrent_amount = rd.total_amount, rd.updated_at = NOW()
		WHERE rd.deleted_at IS NULL AND rd.settled = 0 AND rd.customer_type = 'member'
		AND rd.recurrent_amount < rd.total_amount - 0.005
		AND (
			(c.installment_amount IS NOT NULL AND ABS(rd.recurrent_amount - c.installment_amount) < 0.01)
			OR (c.term_months > 0 AND ABS(rd.recurrent_amount - ROUND(rd.total_amount / c.term_months, 2)) < 0.01)
		)`)
	if res.Error != nil {
		log.Printf("[LoanDeduction] fix recurrent amounts: %v", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("[LoanDeduction] %d loan deduction(s): recurrent amount set to the total amount", res.RowsAffected)
	}
}

// RefreshAllBorrowStatus re-evaluates can_borrow / cannot_borrow_reason for every member (e.g. after
// the default member loan product or its terms change). Credit limits are not changed.
func RefreshAllBorrowStatus(why string) {
	var ids []uint64
	if err := db.DB.Model(&models.Member{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("[MemberLoan] refresh borrow status (%s): %v", why, err)
		return
	}
	credit := NewMemberCreditService()
	failed := 0
	for _, id := range ids {
		if err := db.DB.Transaction(func(tx *gorm.DB) error { return credit.RefreshBorrowStatus(tx, id) }); err != nil {
			failed++
		}
	}
	if p, err := resolveMobileLoanProduct(0); err == nil {
		log.Printf("[MemberLoan] %s: mobile loans now use product %d %q (min KES %.2f, max KES %.2f); borrow status refreshed for %d member(s), %d failed",
			why, p.ID, p.ProductName, p.MinAmount, p.MaxAmount, len(ids)-failed, failed)
	}
}
