package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Approving a loan application from the desktop needs a one-time code sent to the MEMBER's registered
// phone: the officer can only approve once the member reads the code back, so nobody can apply for and
// approve a loan in a member's name without the member. (Mobile cashouts are confirmed by the member's
// own OTP and never reach this step.)

const (
	approvalOTPTTL         = 10 * time.Minute
	approvalOTPResendAfter = 60 * time.Second
	approvalOTPMaxAttempts = 5
)

// LoanApprovalOTP is the code for approving one application (table loan_approval_otps).
type LoanApprovalOTP struct {
	models.BaseModel
	LoanApplicationID uint64     `gorm:"column:loan_application_id;index;not null" json:"loan_application_id"`
	MemberID          uint64     `gorm:"column:member_id;index;not null" json:"member_id"`
	Phone             string     `gorm:"column:phone;size:20;not null" json:"-"`
	OTPHash           string     `gorm:"column:otp_hash;size:128;not null" json:"-"`
	ExpiresAt         time.Time  `gorm:"column:expires_at;not null" json:"expires_at"`
	Attempts          int        `gorm:"column:attempts;default:0" json:"attempts"`
	LastSentAt        time.Time  `gorm:"column:last_sent_at;not null" json:"last_sent_at"`
	VerifiedAt        *time.Time `gorm:"column:verified_at" json:"verified_at"`
	Status            string     `gorm:"column:status;size:16;not null;default:PENDING" json:"status"`
}

func (LoanApprovalOTP) TableName() string { return "loan_approval_otps" }

var approvalOTPSchemaOnce sync.Once

func ensureApprovalOTPSchema() {
	approvalOTPSchemaOnce.Do(func() {
		if err := db.DB.AutoMigrate(&LoanApprovalOTP{}); err != nil {
			log.Printf("[LoanApproval] migrate loan_approval_otps: %v", err)
		}
	})
}

func hashApprovalOTP(applicationID uint64, otp string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not configured")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "loan-approval-otp:%d:%s", applicationID, strings.TrimSpace(otp))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// ApprovalOTPSent is returned when the code goes out.
type ApprovalOTPSent struct {
	ApplicationID     uint64    `json:"application_id"`
	MaskedPhone       string    `json:"masked_phone"`
	ExpiresAt         time.Time `json:"expires_at"`
	ResendAvailableAt time.Time `json:"resend_available_at"`
}

// SendApprovalOTP sends (or re-sends) the approval code to the member's registered phone.
func (s *LoanModuleService) SendApprovalOTP(applicationID, userID uint64) (*ApprovalOTPSent, error) {
	ensureApprovalOTPSchema()
	var app models.LoanApplication
	if err := db.DB.First(&app, applicationID).Error; err != nil {
		return nil, serviceErr(http.StatusNotFound, "loan application not found")
	}
	if app.Status != models.LoanAppSubmitted && app.Status != models.LoanAppUnderReview {
		return nil, serviceErr(http.StatusConflict, "application cannot be approved in status %s", app.Status)
	}
	var member models.Member
	if err := db.DB.Select("id, primary_phone").First(&member, app.MemberID).Error; err != nil {
		return nil, serviceErr(http.StatusNotFound, "member not found")
	}
	phone := utils.NormalizePhone(member.PrimaryPhone)
	if phone == "" {
		return nil, serviceErr(http.StatusUnprocessableEntity, "the member has no valid registered phone number to send the approval code to")
	}

	now := time.Now()
	var otp string
	var record LoanApprovalOTP
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("loan_application_id = ? AND status = ?", app.ID, models.LoanOTPStatusPending).
			Order("id DESC").First(&record).Error
		if err == nil {
			if wait := record.LastSentAt.Add(approvalOTPResendAfter).Sub(now); wait > 0 {
				return serviceErr(http.StatusTooManyRequests, "wait %d seconds before sending another code", int(wait.Seconds())+1)
			}
		} else {
			record = LoanApprovalOTP{LoanApplicationID: app.ID, MemberID: app.MemberID}
			record.CreatedBy = userID
		}
		if otp, err = generateOTP(); err != nil {
			return err
		}
		hash, err := hashApprovalOTP(app.ID, otp)
		if err != nil {
			return err
		}
		record.Phone = phone
		record.OTPHash = hash
		record.ExpiresAt = now.Add(approvalOTPTTL)
		record.Attempts = 0
		record.LastSentAt = now
		record.Status = models.LoanOTPStatusPending
		record.UpdatedBy = userID
		return tx.Save(&record).Error
	})
	if err != nil {
		return nil, err
	}

	amount := app.RequestedAmount
	if os.Getenv("ENV") == "development" {
		log.Printf("[LoanApproval][development] approval OTP for loan application %d: %s", app.ID, otp)
		writeDevOTPLog(app.MemberID, otp+" (approval)", amount, app.ID)
	}
	message := fmt.Sprintf("eDairy: your loan application %s for KES %s is being approved. Approval code %s, valid %d minutes. Give it to the officer ONLY if you applied for this loan.",
		app.ApplicationNo, formatKES(amount), otp, int(approvalOTPTTL.Minutes()))
	if err := notifyMemberSMS(app.MemberID, message, userID); err != nil {
		return nil, serviceErr(http.StatusBadGateway, "could not send the approval code to the member's phone: %s", err.Error())
	}
	return &ApprovalOTPSent{ApplicationID: app.ID, MaskedPhone: maskPhone(phone), ExpiresAt: record.ExpiresAt,
		ResendAvailableAt: now.Add(approvalOTPResendAfter)}, nil
}

// VerifyApprovalOTP checks the code the member read back; it is used once.
func (s *LoanModuleService) VerifyApprovalOTP(applicationID uint64, code string, userID uint64) error {
	ensureApprovalOTPSchema()
	code = strings.TrimSpace(code)
	if code == "" {
		return serviceErr(http.StatusPreconditionRequired, "approval code required: send the code to the member's phone and enter the code they read back")
	}
	var verifyErr error
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var record LoanApprovalOTP
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("loan_application_id = ? AND status = ?", applicationID, models.LoanOTPStatusPending).
			Order("id DESC").First(&record).Error; err != nil {
			verifyErr = serviceErr(http.StatusPreconditionRequired, "no approval code has been sent for this application; send one first")
			return nil
		}
		now := time.Now()
		if now.After(record.ExpiresAt) {
			record.Status = models.LoanOTPStatusExpired
			verifyErr = serviceErr(http.StatusGone, "the approval code has expired; send a new one")
			return tx.Save(&record).Error
		}
		want, err := hashApprovalOTP(applicationID, code)
		if err != nil {
			return err
		}
		if !hmac.Equal([]byte(want), []byte(record.OTPHash)) {
			record.Attempts++
			if record.Attempts >= approvalOTPMaxAttempts {
				record.Status = models.LoanOTPStatusLocked
				verifyErr = serviceErr(http.StatusForbidden, "too many wrong codes; send a new code to the member")
			} else {
				verifyErr = serviceErr(http.StatusUnauthorized, "wrong approval code; %d attempt(s) left", approvalOTPMaxAttempts-record.Attempts)
			}
			return tx.Save(&record).Error
		}
		record.Status = models.LoanOTPStatusVerified
		record.VerifiedAt = &now
		record.UpdatedBy = userID
		return tx.Save(&record).Error
	})
	if err != nil {
		return err
	}
	return verifyErr
}

// expireApprovalOTPs voids open codes when an application is rejected.
func expireApprovalOTPs(applicationID uint64) {
	ensureApprovalOTPSchema()
	db.DB.Model(&LoanApprovalOTP{}).Where("loan_application_id = ? AND status = ?", applicationID, models.LoanOTPStatusPending).
		Update("status", models.LoanOTPStatusExpired)
}

// ReopenApprovalOTP puts back a code that was accepted moments ago when the approval itself then
// failed (e.g. validation), so the officer can retry without texting the member again.
func (s *LoanModuleService) ReopenApprovalOTP(applicationID uint64) {
	db.DB.Model(&LoanApprovalOTP{}).
		Where("loan_application_id = ? AND status = ? AND verified_at >= ?", applicationID, models.LoanOTPStatusVerified, time.Now().Add(-2*time.Minute)).
		Updates(map[string]interface{}{"status": models.LoanOTPStatusPending, "verified_at": nil})
}

// EnsureApprovalOTPPermission lets every role that can approve applications also send the approval code.
func EnsureApprovalOTPPermission() {
	ensureApprovalOTPSchema()
	const name = "loan-module.applications.approval-otp.create"
	db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
		SELECT ?, 'web', NOW(), NOW() FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = ?)`, name, name)
	if err := db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT DISTINCT rp.role_id, np.id FROM role_permissions rp
		JOIN permissions src ON src.id = rp.permission_id AND src.deleted_at IS NULL
			AND src.name IN ('loan-module.applications.approve.create', 'loan-module.applications.create')
		JOIN permissions np ON np.name = ? AND np.deleted_at IS NULL
		WHERE NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = rp.role_id AND x.permission_id = np.id)`, name).Error; err != nil {
		log.Printf("[LoanApproval] grant %s: %v", name, err)
	}
	db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r JOIN permissions p ON p.name = ? AND p.deleted_at IS NULL
		WHERE r.name = 'Super Admin' AND NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = r.id AND x.permission_id = p.id)`, name)

	// Retrying M-Pesa-failed payouts: whoever can confirm disbursements.
	for _, perm := range []string{"loan-module.disbursements.retry.create", "loan-module.disbursements.retry-mpesa-failed.create"} {
		db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
			SELECT ?, 'web', NOW(), NOW() FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = ?)`, perm, perm)
		db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
			SELECT DISTINCT rp.role_id, np.id FROM role_permissions rp
			JOIN permissions src ON src.id = rp.permission_id AND src.deleted_at IS NULL
				AND src.name IN ('loan-module.disbursements.confirm.create', 'loan-module.contracts.disburse.create')
			JOIN permissions np ON np.name = ? AND np.deleted_at IS NULL
			WHERE NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = rp.role_id AND x.permission_id = np.id)`, perm)
		db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
			SELECT r.id, p.id FROM roles r JOIN permissions p ON p.name = ? AND p.deleted_at IS NULL
			WHERE r.name = 'Super Admin' AND NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = r.id AND x.permission_id = p.id)`, perm)
	}
}
