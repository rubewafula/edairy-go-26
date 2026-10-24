package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

// A payout is given up and the loan cancelled (deduction reversed, credit restored, member told by SMS
// and asked to try again) when:
//   - it has been tried DISBURSEMENT_MAX_ATTEMPTS times (default 4), or
//   - it has been in error for more than DISBURSEMENT_ERROR_MAX_AGE_MINUTES (default 60); this is
//     checked when the worker next picks it up.
const (
	DisbursementCancelRetriesExhausted = "RETRIES_EXHAUSTED"
	DisbursementCancelErrorExpired     = "DISBURSEMENT_EXPIRED"
)

func maxDisbursementAttempts() int {
	if n := envInt("DISBURSEMENT_MAX_ATTEMPTS", 4); n > 0 {
		return n
	}
	return 4
}

func disbursementErrorMaxAge() time.Duration {
	if m := envInt("DISBURSEMENT_ERROR_MAX_AGE_MINUTES", 60); m > 0 {
		return time.Duration(m) * time.Minute
	}
	return time.Hour
}

func isAutoCancelCode(code string) bool {
	return code == DisbursementCancelRetriesExhausted || code == DisbursementCancelErrorExpired
}

func payloadAttempts(payload map[string]interface{}) int {
	switch v := payload["attempts"].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

func payloadLastError(payload map[string]interface{}) string {
	if v, ok := payload["last_error"].(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

// disbursementFirstErrorAt is when the payout first went into error (older rows: when it was queued).
func disbursementFirstErrorAt(d *models.LoanDisbursement, payload map[string]interface{}) time.Time {
	if v, ok := payload["first_error_at"].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
	}
	return d.CreatedAt
}

// firstErrorAtPatch keeps the time of the first error (set once).
func firstErrorAtPatch(payload map[string]interface{}) string {
	if v, ok := payload["first_error_at"].(string); ok && v != "" {
		return v
	}
	return time.Now().Format(time.RFC3339)
}

// expiredDisbursementReason says why a payout about to be retried must be cancelled instead ("" = retry).
func expiredDisbursementReason(d *models.LoanDisbursement, payload map[string]interface{}) (string, string) {
	max := maxDisbursementAttempts()
	if payloadAttempts(payload) >= max {
		return DisbursementCancelRetriesExhausted, fmt.Sprintf("payout failed after %d attempts", max)
	}
	if payloadLastError(payload) != "" {
		age := disbursementErrorMaxAge()
		if since := time.Since(disbursementFirstErrorAt(d, payload)); since > age {
			return DisbursementCancelErrorExpired, fmt.Sprintf("payout still failing after %s", age.Round(time.Minute))
		}
	}
	return "", ""
}

// cancelFailedDisbursement gives up on a payout: the disbursement is FAILED and the loan released
// (application cancelled, milk deduction reversed, credit restored, member SMS with the restored
// balance asking them to try again).
func (s *LoanModuleService) cancelFailedDisbursement(d *models.LoanDisbursement, code, reason string) (*models.LoanDisbursement, error) {
	db.DB.Model(&models.LoanDisbursement{}).Where("id = ?", d.ID).Update("channel_payload",
		mergeChannelPayload([]byte(d.ChannelPayload), map[string]interface{}{
			"queued": false, "cancelled_at": time.Now().Format(time.RFC3339), "cancel_reason": reason,
		}))
	out, err := s.ConfirmDisbursement(d.ID, dtos.ConfirmLoanDisbursementRequest{
		Success: false, Release: true, FailureCode: code, FailureMessage: reason,
	}, disbursementWorkerUser)
	if err != nil {
		log.Printf("[DisbursementQueue] cancel disbursement %d (%s): %v", d.ID, code, err)
		appendAppLog("disbursement.log", fmt.Sprintf("disbursement %d cancel FAILED (%s): %v", d.ID, reason, err))
		return nil, err
	}
	log.Printf("[DisbursementQueue] disbursement %d CANCELLED (%s): %s; loan released", d.ID, code, reason)
	appendAppLog("disbursement.log", fmt.Sprintf("disbursement %d CANCELLED (%s): %s - loan released, deduction reversed, member notified", d.ID, code, reason))
	return out, nil
}
