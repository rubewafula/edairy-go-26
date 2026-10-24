package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

// Disbursement queue
//
// Mobile loans are not paid out inside the OTP request. VerifyOTP approves the loan, reserves it against
// the credit limit, adds the repayment deduction and records a PENDING disbursement marked "queued".
// This worker picks queued disbursements up and sends them to the provider (Safaricom M-Pesa B2C):
//   - provider accepted          -> PROCESSING; the Daraja result callback later confirms SUCCESS/FAILED
//   - provider confirmed at once -> SUCCESS (finalised via ConfirmDisbursement)
//   - could not be sent          -> stays queued and is retried (M-Pesa not configured, network error, ...)
//   - gave up after max attempts -> FAILED; the loan is released (credit limit restored, deduction cancelled)
//
// Settings (env): DISBURSEMENT_WORKER_ENABLED (default true), DISBURSEMENT_WORKER_INTERVAL_SECONDS (15),
// DISBURSEMENT_MAX_ATTEMPTS (20). In ENV=development only, DISBURSEMENT_SIMULATE_SUCCESS=true marks queued
// payouts as successful with a SIM- reference instead of calling M-Pesa.

const disbursementWorkerUser uint64 = 0

var disbursementWorkerOnce sync.Once

func envInt(key string, def int) int {
	if v, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key))); err == nil && v > 0 {
		return v
	}
	return def
}

// StartDisbursementWorker runs the queue worker in the background for the life of the process.
func StartDisbursementWorker() {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("DISBURSEMENT_WORKER_ENABLED")), "false") {
		log.Printf("[DisbursementQueue] worker disabled (DISBURSEMENT_WORKER_ENABLED=false)")
		return
	}
	disbursementWorkerOnce.Do(func() {
		interval := time.Duration(envInt("DISBURSEMENT_WORKER_INTERVAL_SECONDS", 15)) * time.Second
		log.Printf("[DisbursementQueue] worker started (every %s)", interval)
		appendAppLog("disbursement.log", fmt.Sprintf("worker started (every %s)", interval))
		go func() {
			svc := NewLoanModuleService()
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			lastBeat := time.Time{}
			for {
				svc.ProcessDisbursementQueue()
				if time.Since(lastBeat) >= 15*time.Minute {
					lastBeat = time.Now()
					writeDisbursementHeartbeat()
				}
				<-ticker.C
			}
		}()
	})
}

// ProcessDisbursementQueue sends every due queued disbursement to its provider once.
func (s *LoanModuleService) ProcessDisbursementQueue() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[DisbursementQueue] panic: %v", r)
		}
	}()
	var queued []models.LoanDisbursement
	if err := db.DB.Where("status IN ? AND JSON_EXTRACT(channel_payload, '$.queued') = true",
		[]string{models.DisbursementStatusPending, models.DisbursementStatusMpesaError}).
		Order("id ASC").Limit(50).Find(&queued).Error; err != nil {
		log.Printf("[DisbursementQueue] list: %v", err)
		return
	}
	now := time.Now()
	for i := range queued {
		d := &queued[i]
		payload := decodeChannelPayload([]byte(d.ChannelPayload))
		if next := payloadStringFromJSON(payload, "next_attempt_at"); next != "" {
			if t, err := time.Parse(time.RFC3339, next); err == nil && now.Before(t) {
				continue
			}
		}
		s.processQueuedDisbursement(d.ID)
	}
}

func (s *LoanModuleService) processQueuedDisbursement(id uint64) {
	// Claim the row so two workers (or two ticks) never send the same payout.
	claim := db.DB.Model(&models.LoanDisbursement{}).
		Where("id = ? AND status IN ?", id, []string{models.DisbursementStatusPending, models.DisbursementStatusMpesaError}).
		Update("status", models.DisbursementStatusProcessing)
	if claim.Error != nil || claim.RowsAffected != 1 {
		return
	}
	var d models.LoanDisbursement
	if err := db.DB.First(&d, id).Error; err != nil {
		return
	}
	var contract models.LoanContract
	if err := db.DB.First(&contract, d.LoanContractID).Error; err != nil {
		s.requeueDisbursement(&d, "contract not found: "+err.Error())
		return
	}
	payload := decodeChannelPayload([]byte(d.ChannelPayload))

	// Tried too often, or in error for over an hour: cancel the loan instead of trying again.
	if code, reason := expiredDisbursementReason(&d, payload); code != "" {
		_, _ = s.cancelFailedDisbursement(&d, code, reason)
		return
	}

	// Re-validate right before money leaves: if the member no longer qualifies, cancel the application
	// (the failure handling releases the loan: credit restored, deduction cancelled, member notified).
	if reason := revalidateBeforeDisbursement(&contract); reason != "" {
		log.Printf("[DisbursementQueue] disbursement %d (%s) cancelled on re-validation: %s", d.ID, contract.ContractNo, reason)
		db.DB.Model(&models.LoanDisbursement{}).Where("id = ?", d.ID).Update("channel_payload",
			mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{"queued": false, "last_error": reason}))
		if _, err := s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
			Success: false, FailureCode: "REVALIDATION_FAILED", FailureMessage: reason,
		}, disbursementWorkerUser); err != nil {
			log.Printf("[DisbursementQueue] cancel disbursement %d: %v", d.ID, err)
		}
		return
	}

	if os.Getenv("ENV") == "development" && strings.EqualFold(os.Getenv("DISBURSEMENT_SIMULATE_SUCCESS"), "true") {
		ref := fmt.Sprintf("SIM-%d", d.ID)
		log.Printf("[DisbursementQueue] SIMULATED payout of KES %s for %s (ref %s)", formatKES(d.Amount), contract.ContractNo, ref)
		if _, err := s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
			Success: true, ExternalReference: ref, PostToGL: true, CreateMilkDeduction: true,
		}, disbursementWorkerUser); err != nil {
			s.requeueDisbursement(&d, "simulated confirm failed: "+err.Error())
		}
		return
	}

	channelCode := d.ChannelCode
	channel, err := s.channels.GetByCode(channelCode)
	if err != nil {
		s.requeueDisbursement(&d, "channel "+channelCode+": "+err.Error())
		return
	}

	result := s.channels.Execute(DisbursementExecutionInput{
		Channel:    channel,
		NetAmount:  d.Amount,
		Payload:    payload,
		ContractNo: contract.ContractNo,
		MemberID:   contract.MemberID,
		Reference:  d.Reference,
	})

	switch result.Status {
	case models.DisbursementStatusProcessing, models.DisbursementStatusPending:
		// Sent; the provider's result callback will confirm success or failure.
		applyExecutionResult(&d, result)
		d.Status = models.DisbursementStatusProcessing
		d.ChannelPayload = mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{
			"queued": false, "sent_at": time.Now().Format(time.RFC3339), "last_error": "",
		})
		if err := db.DB.Save(&d).Error; err != nil {
			log.Printf("[DisbursementQueue] save sent disbursement %d: %v", d.ID, err)
		}
		log.Printf("[DisbursementQueue] disbursement %d (%s, KES %s) sent to %s, awaiting provider confirmation",
			d.ID, contract.ContractNo, formatKES(d.Amount), channelCode)
		appendAppLog("disbursement.log", fmt.Sprintf("disbursement %d (%s, KES %s) SENT to %s, awaiting M-Pesa result",
			d.ID, contract.ContractNo, formatKES(d.Amount), channelCode))
	case models.DisbursementStatusSuccess:
		applyExecutionResult(&d, result)
		d.Status = models.DisbursementStatusProcessing
		db.DB.Save(&d)
		if _, err := s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
			Success: true, ExternalReference: d.ExternalReference, PostToGL: true, CreateMilkDeduction: true,
		}, disbursementWorkerUser); err != nil {
			log.Printf("[DisbursementQueue] confirm disbursement %d: %v", d.ID, err)
		}
	default:
		s.requeueDisbursement(&d, strings.TrimSpace(result.FailureCode+": "+result.FailureMessage))
	}
}

// requeueDisbursement puts a payout that could not be sent back in the queue with a backoff, or gives
// up after DISBURSEMENT_MAX_ATTEMPTS (the loan is then released by the normal failure handling).
func (s *LoanModuleService) requeueDisbursement(d *models.LoanDisbursement, reason string) {
	payload := decodeChannelPayload([]byte(d.ChannelPayload))
	attempts := payloadAttempts(payload) + 1 // every try counts

	configProblem := strings.Contains(strings.ToLower(reason), "not configured") ||
		strings.Contains(strings.ToLower(reason), "payouts are disabled")

	if max := maxDisbursementAttempts(); attempts >= max {
		d.ChannelPayload = mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{
			"attempts": attempts, "last_error": reason, "first_error_at": firstErrorAtPatch(payload),
		})
		_, _ = s.cancelFailedDisbursement(d, DisbursementCancelRetriesExhausted,
			fmt.Sprintf("payout failed after %d attempts: %s", attempts, reason))
		return
	}

	backoff := time.Duration(30*(1<<minInt(attempts, 6))) * time.Second // 1m, 2m, 4m ... capped at 32m
	if configProblem {
		backoff = 5 * time.Minute
	}
	patch := map[string]interface{}{
		"queued":          true,
		"attempts":        attempts,
		"last_error":      reason,
		"first_error_at":  firstErrorAtPatch(payload),
		"next_attempt_at": time.Now().Add(backoff).Format(time.RFC3339),
	}
	code := "MPESA_SEND_FAILED"
	if configProblem {
		code = "MPESA_NOT_AVAILABLE"
	}
	if err := db.DB.Model(&models.LoanDisbursement{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
		"status":          models.DisbursementStatusMpesaError,
		"channel_payload": mergeChannelPayload([]byte(d.ChannelPayload), patch),
		"failure_code":    code,
		"failure_message": reason,
	}).Error; err != nil {
		log.Printf("[DisbursementQueue] requeue disbursement %d: %v", d.ID, err)
		return
	}
	log.Printf("[DisbursementQueue] disbursement %d MPESA_ERROR, re-queued (attempt %d, retry in %s): %s", d.ID, attempts, backoff, reason)
	appendAppLog("disbursement.log", fmt.Sprintf("disbursement %d MPESA_ERROR, re-queued (attempt %d, retry in %s): %s", d.ID, attempts, backoff, reason))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// revalidateBeforeDisbursement checks a queued loan is still valid just before it is sent.
// The loan has already been taken off the credit limit, so the limit must simply not be negative
// (still covered). Returns "" when it may be sent, else the reason it may not.
func revalidateBeforeDisbursement(contract *models.LoanContract) string {
	var member models.Member
	if err := db.DB.Select("id, status, credit_limit").First(&member, contract.MemberID).Error; err != nil {
		return "member not found"
	}
	if !memberCanBorrow(member.Status) {
		return fmt.Sprintf("Member account is %s and cannot borrow.", strings.ToLower(strings.TrimSpace(member.Status)))
	}
	if contract.LoanApplicationID != nil {
		var app models.LoanApplication
		if err := db.DB.Select("id, status").First(&app, *contract.LoanApplicationID).Error; err == nil &&
			app.Status != models.LoanAppApproved {
			return fmt.Sprintf("Loan application is %s, not approved.", strings.ToLower(app.Status))
		}
	}
	if contract.Status != models.LoanContractPending {
		return fmt.Sprintf("Loan contract is %s, not awaiting disbursement.", strings.ToLower(strings.ReplaceAll(contract.Status, "_", " ")))
	}
	if member.CreditLimit < -0.005 {
		return fmt.Sprintf("Credit limit is KES %s after this loan; it is no longer covered by the member's credit.", formatKES(member.CreditLimit))
	}
	return ""
}
