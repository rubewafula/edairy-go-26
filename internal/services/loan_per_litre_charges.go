package services

import (
	"fmt"
	"log"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// ChargeTypePerLitre is the loan_charges.charge_type for the member-loan per-litre charge.
const ChargeTypePerLitre = "PER_LITRE"

// ApplyPerLitreCharges charges members who took a cashout (member loan) in the pay period:
// charge = product.ChargePerLitre × confirmed litres the member delivered in the period.
//
// Charged once per member per period (re-running payroll does not charge again), on the member's latest
// cashout contract of the period. The charge is recorded in loan_charges, added to that contract's
// schedule (fee due) and outstanding fees, and to its LOAN-<contract_no> milk deduction so payroll
// recovers it; the member's stored credit limit is then reduced by the charge.
func (s *LoanModuleService) ApplyPerLitreCharges(pdr models.MemberPayDateRange, userID uint64) (int, error) {
	product, err := resolveMemberLoanProduct()
	if err != nil || product.ChargePerLitre <= 0 {
		return 0, nil
	}
	start := pdr.StartDate
	endExclusive := pdr.EndDate.AddDate(0, 0, 1)
	chargedDate := time.Date(pdr.EndDate.Year(), pdr.EndDate.Month(), pdr.EndDate.Day(), 0, 0, 0, 0, pdr.EndDate.Location())

	// Latest live member-loan contract per member taken in the period.
	var contracts []models.LoanContract
	if err := db.DB.Where(`loan_product_id = ? AND status IN ? AND created_at >= ? AND created_at < ? AND deleted_at IS NULL
		AND id IN (SELECT MAX(c2.id) FROM loan_contracts c2 WHERE c2.loan_product_id = ? AND c2.status IN ?
		           AND c2.created_at >= ? AND c2.created_at < ? AND c2.deleted_at IS NULL GROUP BY c2.member_id)`,
		product.ID, perLitreContractStatuses, start, endExclusive,
		product.ID, perLitreContractStatuses, start, endExclusive).
		Find(&contracts).Error; err != nil {
		return 0, err
	}

	applied := 0
	credit := NewMemberCreditService()
	for i := range contracts {
		c := contracts[i]
		ok, err := s.applyPerLitreCharge(&c, product.ChargePerLitre, start, pdr.EndDate, chargedDate, userID)
		if err != nil {
			log.Printf("[PerLitreCharge] member %d contract %s: %v", c.MemberID, c.ContractNo, err)
			continue
		}
		if ok {
			applied++
			if _, err := credit.RecalculateNow(c.MemberID, CreditChange{
				Action:     models.CreditActionLoanInterest,
				SourceType: "loan_charge",
				SourceID:   c.ID,
				Reference:  c.ContractNo,
				Narration:  fmt.Sprintf("Per-litre cashout charge for %s", pdr.Name),
				UserID:     userID,
			}); err != nil {
				log.Printf("[PerLitreCharge] credit update for member %d: %v", c.MemberID, err)
			}
		}
	}
	if applied > 0 {
		log.Printf("[PerLitreCharge] %s: charged %d member(s) at KES %.2f/litre", pdr.Name, applied, product.ChargePerLitre)
	}
	return applied, nil
}

var perLitreContractStatuses = []string{models.LoanContractPending, models.LoanContractActive, models.LoanContractOverdue}

func (s *LoanModuleService) applyPerLitreCharge(c *models.LoanContract, rate float64, start, end, chargedDate time.Time, userID uint64) (bool, error) {
	applied := false
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Once per member per period, whichever of their contracts carried it.
		var existing int64
		if err := tx.Table("loan_charges ch").
			Joins("JOIN loan_contracts lc ON lc.id = ch.loan_contract_id").
			Where("lc.member_id = ? AND ch.charge_type = ? AND ch.charged_date = ? AND ch.deleted_at IS NULL",
				c.MemberID, ChargeTypePerLitre, chargedDate.Format("2006-01-02")).
			Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return nil
		}

		var litres float64
		if err := tx.Raw(`SELECT COALESCE(SUM(e.quantity), 0) FROM milk_journal_entries e
			JOIN milk_journals j ON j.id = e.milk_journal_id
			WHERE e.member_id = ? AND e.deleted_at IS NULL AND j.deleted_at IS NULL AND j.confirmed = 1
			AND j.journal_date BETWEEN ? AND ?`, c.MemberID, start.Format("2006-01-02"), end.Format("2006-01-02")).
			Scan(&litres).Error; err != nil {
			return err
		}
		amount := roundMoney(litres * rate)
		if amount <= 0 {
			return nil
		}

		if err := tx.Create(&models.LoanChargeRecord{
			BaseModel:      models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			LoanContractID: c.ID,
			ChargeType:     ChargeTypePerLitre,
			Amount:         amount,
			ChargedDate:    chargedDate,
		}).Error; err != nil {
			return err
		}

		// Owed on the contract: last instalment's fee, contract outstanding fees.
		var inst models.LoanScheduleInstallment
		if err := tx.Where("loan_contract_id = ? AND deleted_at IS NULL", c.ID).
			Order("installment_no DESC").First(&inst).Error; err == nil {
			if err := tx.Model(&inst).Updates(map[string]interface{}{
				"fee_due":   gorm.Expr("fee_due + ?", amount),
				"total_due": gorm.Expr("total_due + ?", amount),
			}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&models.LoanContract{}).Where("id = ?", c.ID).
			Update("outstanding_fees", gorm.Expr("outstanding_fees + ?", amount)).Error; err != nil {
			return err
		}

		// Recovered through the loan's milk deduction.
		if err := tx.Model(&models.RecurrentDeduction{}).
			Where("reference = ? AND customer_id = ? AND customer_type = 'member' AND settled = 0", "LOAN-"+c.ContractNo, c.MemberID).
			Updates(map[string]interface{}{
				"total_amount":     gorm.Expr("total_amount + ?", amount),
				"recurrent_amount": gorm.Expr("recurrent_amount + ?", amount),
			}).Error; err != nil {
			return err
		}
		log.Printf("[PerLitreCharge] member %d: %.2f L × KES %.2f = KES %.2f on %s", c.MemberID, litres, rate, amount, c.ContractNo)
		applied = true
		return nil
	})
	return applied, err
}
