package services

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// handleDisbursementOutcome runs after a disbursement reaches SUCCESS or FAILED and its transaction has committed.
func (s *LoanModuleService) handleDisbursementOutcome(disbursementID, userID uint64) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[LoanDisbursementOutcome] panic for disbursement %d: %v", disbursementID, r)
			}
		}()
		var d models.LoanDisbursement
		if err := db.DB.First(&d, disbursementID).Error; err != nil {
			return
		}
		var contract models.LoanContract
		if err := db.DB.First(&contract, d.LoanContractID).Error; err != nil {
			return
		}
		switch d.Status {
		case models.DisbursementStatusSuccess:
			s.onDisbursementSucceeded(&contract, &d, userID)
		case models.DisbursementStatusFailed:
			// Mobile loans, and any loan cancelled after too many / too old payout failures, are released.
			if isMobileContract(&contract) || isAutoCancelCode(d.FailureCode) {
				reason := strings.TrimSpace(d.FailureMessage)
				if reason == "" {
					reason = d.FailureCode
				}
				if err := s.releaseMobileContract(contract.ID, reason, userID); err != nil {
					log.Printf("[LoanDisbursementOutcome] release contract %d: %v", contract.ID, err)
				}
			}
		}
	}()
}

func (s *LoanModuleService) onDisbursementSucceeded(contract *models.LoanContract, d *models.LoanDisbursement, userID uint64) {
	// The amount was already taken off the stored credit limit when the loan was approved,
	// so disbursement only refreshes can_borrow; it does not recompute the limit.
	credit := NewMemberCreditService()
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		return credit.RefreshBorrowStatus(tx, contract.MemberID)
	}); err != nil {
		log.Printf("[LoanDisbursementOutcome] borrow status refresh for member %d: %v", contract.MemberID, err)
	}
	txnRef := ""
	if d.ExternalReference != "" {
		txnRef = " (M-Pesa ref " + d.ExternalReference + ")"
	}
	message := fmt.Sprintf("Dear %s, your loan of KES %s%s for %s has been sent to you. Amount to be deducted from your milk: KES %s. Your milk balance after this deduction: KES %s.",
		memberFirstName(contract.MemberID), formatKES(d.Amount), txnRef, contract.ContractNo,
		formatKES(contractRemainingRepayable(db.DB, contract.ID)), formatKES(memberActualBalance(contract.MemberID)))
	if err := notifyMemberSMS(contract.MemberID, message, userID); err != nil {
		log.Printf("[LoanDisbursementOutcome] confirmation SMS for member %d: %v", contract.MemberID, err)
	}
}

// releaseMobileContract cancels a mobile loan whose payout definitively failed, so its reserved amount returns to the member.
func (s *LoanModuleService) releaseMobileContract(contractID uint64, reason string, userID uint64) error {
	var contract models.LoanContract
	released := false
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&contract, contractID).Error; err != nil {
			return err
		}
		if contract.Status != models.LoanContractPending || contract.DisbursedAmount > 0 {
			return nil
		}
		var live int64
		if err := tx.Model(&models.LoanDisbursement{}).
			Where("loan_contract_id = ? AND status IN ?", contract.ID,
				[]string{models.DisbursementStatusPending, models.DisbursementStatusProcessing, models.DisbursementStatusSuccess}).
			Count(&live).Error; err != nil {
			return err
		}
		if live > 0 {
			return nil
		}
		if err := tx.Model(&contract).Updates(map[string]interface{}{
			"status":                models.LoanContractClosed,
			"outstanding_principal": 0,
			"outstanding_interest":  0,
			"outstanding_fees":      0,
			"outstanding_penalties": 0,
			"updated_by":            userID,
		}).Error; err != nil {
			return err
		}
		if contract.LoanApplicationID != nil {
			if err := tx.Model(&models.LoanApplication{}).Where("id = ?", *contract.LoanApplicationID).Updates(map[string]interface{}{
				"status":          models.LoanAppCancelled,
				"rejected_reason": "Disbursement failed: " + reason,
				"updated_by":      userID,
			}).Error; err != nil {
				return err
			}
		}
		// The repayment deduction was added when the loan was queued; nothing was paid out, so cancel it.
		if err := tx.Where("reference = ? AND customer_id = ? AND customer_type = ? AND settled = 0 AND paid_amount = 0",
			"LOAN-"+contract.ContractNo, contract.MemberID, "member").
			Delete(&models.RecurrentDeduction{}).Error; err != nil {
			return err
		}
		released = true
		return nil
	})
	if err != nil || !released {
		return err
	}

	// The closed contract no longer counts, so Recalculate gives its amount back to the stored limit.
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		_, err := NewMemberCreditService().Recalculate(tx, contract.MemberID, CreditChange{
			Action:     models.CreditActionLoanReleased,
			SourceType: "loan_contract",
			SourceID:   contract.ID,
			Reference:  contract.ContractNo,
			Narration:  "Mobile loan released after failed disbursement",
			UserID:     userID,
		})
		return err
	}); err != nil {
		log.Printf("[LoanDisbursementOutcome] credit recalculation for member %d: %v", contract.MemberID, err)
	}
	restored := memberActualBalance(contract.MemberID)
	message := fmt.Sprintf("Dear %s, your loan request of KES %s (%s) could not be sent and has been cancelled. Nothing will be deducted from your milk. Your available balance is restored to KES %s. Please try again.",
		memberFirstName(contract.MemberID), formatKES(contract.ApprovedAmount), contract.ContractNo, formatKES(restored))
	if err := notifyMemberSMS(contract.MemberID, message, userID); err != nil {
		log.Printf("[LoanDisbursementOutcome] failure SMS for member %d: %v", contract.MemberID, err)
	}
	return nil
}

func isMobileContract(contract *models.LoanContract) bool {
	if contract.LoanApplicationID == nil {
		return false
	}
	var app models.LoanApplication
	if err := db.DB.Select("id, source").First(&app, *contract.LoanApplicationID).Error; err != nil {
		return false
	}
	return app.Source == models.LoanAppSourceMobile
}

// contractRemainingRepayable is what is still to be recovered from the member's milk for this contract.
func contractRemainingRepayable(tx *gorm.DB, contractID uint64) float64 {
	var total float64
	tx.Raw(`SELECT COALESCE(SUM(GREATEST(
			(principal_due + interest_due + fee_due + insurance_due + penalty_due) -
			(principal_paid + interest_paid + fee_paid + insurance_paid + penalty_paid), 0)), 0)
		FROM loan_schedule_installments WHERE loan_contract_id = ? AND deleted_at IS NULL`, contractID).Scan(&total)
	return roundMoney(total)
}

// memberActualBalance is the member's milk balance (valued at the payout rate) after all loan
// deductions, i.e. member_registrations.member_balance as kept by the credit recalculation.
func memberActualBalance(memberID uint64) float64 {
	var m models.Member
	if err := db.DB.Select("id, credit_limit, member_balance").First(&m, memberID).Error; err != nil {
		return 0
	}
	bal := m.MemberBalance
	if bal == 0 && m.CreditLimit != 0 {
		if r := GetMemberRates(db.DB, memberID).Ratio(); r > 0 {
			bal = m.CreditLimit / r
		}
	}
	return roundMoney(math.Max(0, bal))
}

// notifyMemberSMS sends an SMS to the member's registered primary phone.
func notifyMemberSMS(memberID uint64, message string, userID uint64) error {
	var member models.Member
	if err := db.DB.Select("id, primary_phone").First(&member, memberID).Error; err != nil {
		return fmt.Errorf("member %d not found", memberID)
	}
	phone := utils.NormalizePhone(member.PrimaryPhone)
	if phone == "" {
		return fmt.Errorf("member %d has no valid primary phone", memberID)
	}
	if os.Getenv("SMS_DRY_RUN") == "true" {
		log.Printf("[SMS dry run] to %s: %s", phone, message)
		return nil
	}
	_, err := NewSMSService().CreateSMSOutbox(dtos.CreateSMSOutboxRequest{PhoneNo: phone, Message: message}, userID)
	return err
}

func memberFirstName(memberID uint64) string {
	var member models.Member
	if err := db.DB.Select("id, first_name, last_name").First(&member, memberID).Error; err != nil {
		return "Member"
	}
	if name := strings.TrimSpace(member.FirstName); name != "" {
		return name
	}
	return "Member"
}

// formatKES renders an amount with thousands separators and two decimals, e.g. 12,500.00.
func formatKES(v float64) string {
	negative := v < 0
	s := strconv.FormatFloat(math.Abs(roundMoney(v)), 'f', 2, 64)
	whole, frac := s[:len(s)-3], s[len(s)-2:]
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	out := b.String() + "." + frac
	if negative {
		return "-" + out
	}
	return out
}
