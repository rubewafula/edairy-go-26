package services

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

// unpaidOutstandingSQL is what is still to be repaid on contract c: from its schedule when it has
// installments, else from the contract's outstanding balances, else principal less what was paid.
// Loans without a schedule (e.g. advances) would otherwise count as 0 unpaid.
const unpaidOutstandingSQL = `(CASE
	WHEN EXISTS (SELECT 1 FROM loan_schedule_installments i0 WHERE i0.loan_contract_id = c.id AND i0.deleted_at IS NULL)
		THEN ` + contractOutstandingSQL + `
	WHEN (COALESCE(c.outstanding_principal, 0) + COALESCE(c.outstanding_interest, 0) + COALESCE(c.outstanding_fees, 0) + COALESCE(c.outstanding_penalties, 0)) > 0
		THEN COALESCE(c.outstanding_principal, 0) + COALESCE(c.outstanding_interest, 0) + COALESCE(c.outstanding_fees, 0) + COALESCE(c.outstanding_penalties, 0)
	ELSE GREATEST(COALESCE(c.principal_amount, 0) - COALESCE(c.total_paid, 0), 0)
END)`

// enrichApplicationsCumulativeUnpaid fills UnpaidAmount and CumulativeUnpaid on a page of applications.
// Unpaid = DISBURSED loans whose contract is not settled/closed/written off, valued at what is still
// to be repaid (principal + interest + charges - paid). The running total covers all of the member's
// unpaid loans up to each row (date applied, then id), not only the rows on this page.
func enrichApplicationsCumulativeUnpaid(items []models.LoanApplication) {
	if len(items) == 0 {
		return
	}
	memberSet := map[uint64]bool{}
	var memberIDs []uint64
	for _, it := range items {
		if !memberSet[it.MemberID] {
			memberSet[it.MemberID] = true
			memberIDs = append(memberIDs, it.MemberID)
		}
	}
	var rows []struct {
		ID          uint64
		MemberID    uint64
		DateApplied time.Time
		Outstanding float64
	}
	if err := db.DB.Raw(`SELECT a.id, a.member_id, a.date_applied,
			`+unpaidOutstandingSQL+` AS outstanding
		FROM loan_applications a
		JOIN loan_contracts c ON c.id = a.loan_contract_id AND c.deleted_at IS NULL
		WHERE a.deleted_at IS NULL AND a.status = 'DISBURSED' AND a.member_id IN ?
		AND c.status NOT IN ('SETTLED', 'CLOSED', 'WRITTEN_OFF')`, memberIDs).Scan(&rows).Error; err != nil {
		return
	}
	before := func(aDate time.Time, aID uint64, bDate time.Time, bID uint64) bool {
		if !aDate.Equal(bDate) {
			return aDate.Before(bDate)
		}
		return aID <= bID
	}
	for i := range items {
		it := &items[i]
		var cum, own float64
		for _, r := range rows {
			if r.MemberID != it.MemberID {
				continue
			}
			if r.ID == it.ID {
				own = r.Outstanding
			}
			if before(r.DateApplied, r.ID, it.DateApplied, it.ID) {
				cum += r.Outstanding
			}
		}
		it.UnpaidAmount = roundMoney(own)
		it.CumulativeUnpaid = roundMoney(cum)
	}
}
