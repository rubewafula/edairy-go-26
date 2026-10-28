package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

func isMpesaChannel(code string) bool {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case models.DisbursementCodeMobileMpesa, models.DisbursementCodeMobileDtbMpesa:
		return true
	}
	return false
}

// markMpesaFailed parks a payout that M-Pesa could not make. Nothing is released: the application
// stays APPROVED, the contract PENDING_DISBURSEMENT, the credit stays reserved and the member keeps
// seeing "Processing". An admin retries it (POST /disbursements/:id/retry) or releases it
// (POST /disbursements/:id/confirm with success=false, release=true).
func (s *LoanModuleService) markMpesaFailed(d *models.LoanDisbursement, code, message string) (*models.LoanDisbursement, error) {
	// M-Pesa reported a failure (or could not be used): keep the loan, flag MPESA_ERROR and re-queue so
	// the worker tries again with a backoff. An admin can retry at once or release the loan.
	code = strings.TrimSpace(code)
	if code == "" {
		code = "MPESA_ERROR"
	}
	payload := decodeChannelPayload([]byte(d.ChannelPayload))
	attempts := payloadAttempts(payload) + 1
	if max := maxDisbursementAttempts(); attempts >= max {
		d.ChannelPayload = mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{
			"attempts": attempts, "last_error": strings.TrimSpace(code + ": " + message),
			"first_error_at": firstErrorAtPatch(payload),
		})
		return s.cancelFailedDisbursement(d, DisbursementCancelRetriesExhausted,
			fmt.Sprintf("payout failed after %d attempts: %s", attempts, strings.TrimSpace(code+": "+message)))
	}
	backoff := time.Duration(30*(1<<minInt(attempts, 6))) * time.Second
	newPayload := mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{
		"queued": true, "attempts": attempts, "last_error": strings.TrimSpace(code + ": " + message),
		"first_error_at":  firstErrorAtPatch(payload),
		"next_attempt_at": time.Now().Add(backoff).Format(time.RFC3339),
	})
	if err := db.DB.Model(&models.LoanDisbursement{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
		"status":          models.DisbursementStatusMpesaError,
		"failure_code":    code,
		"failure_message": message,
		"channel_payload": newPayload,
	}).Error; err != nil {
		return nil, err
	}
	d.Status, d.FailureCode, d.FailureMessage, d.ChannelPayload = models.DisbursementStatusMpesaError, code, message, newPayload
	log.Printf("[DisbursementQueue] disbursement %d MPESA_ERROR (%s), re-queued in %s: %s", d.ID, code, backoff, message)
	appendAppLog("disbursement.log", fmt.Sprintf("disbursement %d MPESA_ERROR (%s), re-queued in %s: %s - member still sees Processing", d.ID, code, backoff, message))
	return d, nil
}

// RetryMpesaFailed puts M-Pesa-failed payouts back in the queue (id 0 = all of them). The worker
// re-validates each loan before sending again.
func (s *LoanModuleService) RetryMpesaFailed(disbursementID, userID uint64) (int, error) {
	var list []models.LoanDisbursement
	q := db.DB.Where("status IN ?", []string{models.DisbursementStatusMpesaFailed, models.DisbursementStatusMpesaError})
	if disbursementID > 0 {
		q = q.Where("id = ?", disbursementID)
	}
	if err := q.Find(&list).Error; err != nil {
		return 0, err
	}
	if disbursementID > 0 && len(list) == 0 {
		return 0, serviceErr(409, "disbursement %d is not waiting on M-Pesa (MPESA_ERROR)", disbursementID)
	}
	for _, d := range list {
		payload := mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{
			"queued": true, "attempts": 0, "next_attempt_at": "", "retried_at": time.Now().Format(time.RFC3339),
			"last_error": "", "first_error_at": "",
			"retried_by": userID,
		})
		db.DB.Model(&models.LoanDisbursement{}).Where("id = ? AND status IN ?", d.ID,
			[]string{models.DisbursementStatusMpesaFailed, models.DisbursementStatusMpesaError}).
			Updates(map[string]interface{}{"status": models.DisbursementStatusPending, "channel_payload": payload, "updated_by": userID})
		appendAppLog("disbursement.log", fmt.Sprintf("disbursement %d re-queued by user %d", d.ID, userID))
	}
	return len(list), nil
}

// writeDisbursementHeartbeat records that the worker is alive and what is waiting.
func writeDisbursementHeartbeat() {
	var rows []struct {
		Status string
		N      int64
	}
	db.DB.Raw(`SELECT status, COUNT(*) AS n FROM loan_disbursements WHERE deleted_at IS NULL
		AND status IN ('PENDING','PROCESSING','MPESA_ERROR','MPESA_FAILED') GROUP BY status`).Scan(&rows)
	parts := []string{}
	for _, r := range rows {
		parts = append(parts, fmt.Sprintf("%s %d", r.Status, r.N))
	}
	if len(parts) == 0 {
		parts = append(parts, "nothing waiting")
	}
	appendAppLog("disbursement.log", "worker alive: "+strings.Join(parts, ", "))
}

// RequeueParkedMpesaFailures turns payouts parked as MPESA_FAILED (older builds) into MPESA_ERROR so
// the worker keeps retrying them.
func RequeueParkedMpesaFailures() {
	res := db.DB.Exec(`UPDATE loan_disbursements SET status = ?,
		channel_payload = JSON_SET(COALESCE(channel_payload, '{}'), '$.queued', true, '$.next_attempt_at', '')
		WHERE status = ? AND deleted_at IS NULL`, models.DisbursementStatusMpesaError, models.DisbursementStatusMpesaFailed)
	if res.Error == nil && res.RowsAffected > 0 {
		appendAppLog("disbursement.log", fmt.Sprintf("%d MPESA_FAILED payout(s) re-queued as MPESA_ERROR", res.RowsAffected))
	}
}
