package services

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// Farmer login PINs (4 digits), stored hashed in member_login_pins (one row per member).
// - Set when the member is registered (mandatory on POST /members).
// - Members without a PIN, or who forgot it / got locked out, set a new one with an SMS code
//   sent to their registered phone (phone + ID number + code).
// - 5 wrong PINs lock the PIN until it is reset by SMS code.

const (
	memberPinMaxFailures   = 5
	memberPinOTPTTL        = 5 * time.Minute
	memberPinOTPCooldown   = 60 * time.Second
	memberPinOTPMaxPerHour = 5
	memberPinOTPMaxTries   = 5
)

func ensureMemberPinTable() error {
	return db.DB.Exec(`CREATE TABLE IF NOT EXISTS member_login_pins (
		member_id BIGINT UNSIGNED NOT NULL PRIMARY KEY,
		pin_hash VARCHAR(100) NULL,
		pin_set_at DATETIME NULL,
		failed_attempts INT NOT NULL DEFAULT 0,
		locked TINYINT(1) NOT NULL DEFAULT 0,
		otp_hash VARCHAR(100) NULL,
		otp_expires_at DATETIME NULL,
		otp_attempts INT NOT NULL DEFAULT 0,
		otp_last_sent_at DATETIME NULL,
		otp_window_start DATETIME NULL,
		otp_window_count INT NOT NULL DEFAULT 0,
		updated_at DATETIME NULL)`).Error
}

type memberPinRow struct {
	MemberID       uint64
	PinHash        *string
	FailedAttempts int
	Locked         bool
	OtpHash        *string
	OtpExpiresAt   *time.Time
	OtpAttempts    int
	OtpLastSentAt  *time.Time
	OtpWindowStart *time.Time
	OtpWindowCount int
}

func loadMemberPin(memberID uint64) memberPinRow {
	var r memberPinRow
	db.DB.Raw(`SELECT member_id, pin_hash, failed_attempts, locked, otp_hash, otp_expires_at, otp_attempts,
		otp_last_sent_at, otp_window_start, otp_window_count FROM member_login_pins WHERE member_id = ?`, memberID).Scan(&r)
	r.MemberID = memberID
	return r
}

func pinDigits(s string) string { return strings.TrimSpace(s) }

// ValidateNewMemberPin checks a new PIN: exactly 4 digits and not easy to guess (repeated or
// straight-sequence digits, or the last 4 digits of the member's phone or ID number).
func ValidateNewMemberPin(pin, phone, idNo string) error {
	pin = pinDigits(pin)
	if len(pin) != 4 || strings.Trim(pin, "0123456789") != "" {
		return serviceErr(http.StatusUnprocessableEntity, "The PIN must be exactly 4 digits.")
	}
	if strings.Count(pin, pin[:1]) == 4 {
		return serviceErr(http.StatusUnprocessableEntity, "Choose a PIN that is harder to guess (not the same digit 4 times).")
	}
	if strings.Contains("0123456789", pin) || strings.Contains("9876543210", pin) {
		return serviceErr(http.StatusUnprocessableEntity, "Choose a PIN that is harder to guess (not 1234, 4321 and similar).")
	}
	last4 := func(v string) string {
		d := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, v)
		if len(d) < 4 {
			return ""
		}
		return d[len(d)-4:]
	}
	if l := last4(phone); l != "" && l == pin {
		return serviceErr(http.StatusUnprocessableEntity, "The PIN cannot be the last 4 digits of the phone number.")
	}
	if l := last4(idNo); l != "" && l == pin {
		return serviceErr(http.StatusUnprocessableEntity, "The PIN cannot be the last 4 digits of the ID number.")
	}
	return nil
}

// SetMemberPin stores a new PIN (hashed), clears failures, the lock and any pending SMS code.
func SetMemberPin(memberID uint64, pin string) error {
	if err := ensureMemberPinTable(); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pinDigits(pin)), 10)
	if err != nil {
		return err
	}
	return db.DB.Exec(`INSERT INTO member_login_pins (member_id, pin_hash, pin_set_at, failed_attempts, locked, otp_hash, otp_expires_at, otp_attempts, updated_at)
		VALUES (?, ?, NOW(), 0, 0, NULL, NULL, 0, NOW())
		ON DUPLICATE KEY UPDATE pin_hash = VALUES(pin_hash), pin_set_at = NOW(), failed_attempts = 0, locked = 0,
			otp_hash = NULL, otp_expires_at = NULL, otp_attempts = 0, updated_at = NOW()`, memberID, string(hash)).Error
}

// ResetMemberPinByOffice clears a member's PIN (admin "Reset PIN"); the member sets a new one by SMS code.
func ResetMemberPinByOffice(memberID uint64) error {
	if err := ensureMemberPinTable(); err != nil {
		return err
	}
	return db.DB.Exec(`UPDATE member_login_pins SET pin_hash = NULL, failed_attempts = 0, locked = 0, updated_at = NOW() WHERE member_id = ?`, memberID).Error
}

type memberLoginCandidate struct {
	ID           uint64
	MemberNo     string
	FirstName    string
	LastName     string
	PrimaryPhone string
	IDNo         string `gorm:"column:id_no"`
	Email        string
	Status       string
}

// membersByPhone finds members whose primary phone normalizes to phone (2547XXXXXXXX).
func membersByPhone(phone string) []memberLoginCandidate {
	var rows []memberLoginCandidate
	if len(phone) < 9 {
		return rows
	}
	db.DB.Raw(`SELECT id, member_no, first_name, last_name, primary_phone, id_no, email, status
		FROM member_registrations
		WHERE deleted_at IS NULL AND REPLACE(REPLACE(primary_phone, ' ', ''), '-', '') LIKE ?
		ORDER BY id`, "%"+phone[len(phone)-9:]).Scan(&rows)
	out := rows[:0]
	for _, m := range rows {
		if utils.NormalizePhone(m.PrimaryPhone) == phone {
			out = append(out, m)
		}
	}
	return out
}

func blockedMemberStatus(status string) bool {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "EXITED", "DECEASED", "DELETED", "CLOSED", "REJECTED":
		return true
	}
	return false
}

func normalizeIDNo(v string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(v), " ", ""))
}

// findMemberByPhoneAndID is used for PIN setup/reset: phone and ID number must belong to one member.
func findMemberByPhoneAndID(phone, idNo string) *memberLoginCandidate {
	for _, m := range membersByPhone(phone) {
		if normalizeIDNo(m.IDNo) == idNo && idNo != "" {
			c := m
			return &c
		}
	}
	return nil
}

// MemberPinOTPRequest is POST /api/member-pin/otp.
type MemberPinOTPRequest struct {
	Phone    string `json:"phone"`
	IDNumber string `json:"id_number"`
}

type MemberPinOTPResponse struct {
	MaskedPhone       string `json:"masked_phone"`
	ExpiresInSeconds  int    `json:"expires_in_seconds"`
	ResendAfterSecond int    `json:"resend_after_seconds"`
}

// RequestMemberPinOTP texts a 6-digit code to the member's registered phone for setting / resetting the PIN.
func RequestMemberPinOTP(req MemberPinOTPRequest, clientIP string) (*MemberPinOTPResponse, error) {
	phone := utils.NormalizePhone(req.Phone)
	idNo := normalizeIDNo(req.IDNumber)
	if phone == "" {
		return nil, serviceErr(http.StatusBadRequest, "Enter a valid phone number, e.g. 0712345678, 712345678 or 254712345678.")
	}
	if idNo == "" {
		return nil, serviceErr(http.StatusBadRequest, "Enter your ID number.")
	}
	keys := []string{"pinotp-phone:" + phone, "pinotp-ip:" + clientIP}
	if memberLoginBlocked(keys...) {
		return nil, serviceErr(http.StatusTooManyRequests, "Too many failed attempts. Please wait 15 minutes and try again.")
	}
	if err := ensureMemberPinTable(); err != nil {
		return nil, err
	}
	m := findMemberByPhoneAndID(phone, idNo)
	if m == nil {
		memberLoginFailed(keys...)
		return nil, serviceErr(http.StatusUnauthorized, "Phone number or ID number is incorrect. Use the phone number registered with the dairy.")
	}
	if blockedMemberStatus(m.Status) {
		return nil, serviceErr(http.StatusForbidden, "This member account is %s. Please contact the dairy office.", strings.ToLower(m.Status))
	}
	r := loadMemberPin(m.ID)
	now := time.Now()
	if r.OtpLastSentAt != nil && now.Sub(*r.OtpLastSentAt) < memberPinOTPCooldown {
		wait := int((memberPinOTPCooldown - now.Sub(*r.OtpLastSentAt)).Seconds()) + 1
		return nil, serviceErr(http.StatusTooManyRequests, "Please wait %d seconds before requesting another code.", wait)
	}
	windowCount := r.OtpWindowCount
	windowStart := now
	if r.OtpWindowStart != nil && now.Sub(*r.OtpWindowStart) < time.Hour {
		windowStart = *r.OtpWindowStart
	} else {
		windowCount = 0
	}
	if windowCount >= memberPinOTPMaxPerHour {
		return nil, serviceErr(http.StatusTooManyRequests, "Too many codes requested. Please try again in an hour or contact the dairy office.")
	}
	otp, err := generateOTP()
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(otp), 10)
	if err != nil {
		return nil, err
	}
	if err := db.DB.Exec(`INSERT INTO member_login_pins (member_id, otp_hash, otp_expires_at, otp_attempts, otp_last_sent_at, otp_window_start, otp_window_count, updated_at)
		VALUES (?, ?, ?, 0, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE otp_hash = VALUES(otp_hash), otp_expires_at = VALUES(otp_expires_at), otp_attempts = 0,
			otp_last_sent_at = VALUES(otp_last_sent_at), otp_window_start = VALUES(otp_window_start),
			otp_window_count = VALUES(otp_window_count), updated_at = NOW()`,
		m.ID, string(hash), now.Add(memberPinOTPTTL), now, windowStart, windowCount+1).Error; err != nil {
		return nil, err
	}
	if os.Getenv("ENV") == "development" {
		writeDevPinOTPLog(m.ID, m.MemberNo, phone, otp)
	}
	msg := fmt.Sprintf("Your eDairy PIN code is %s. Valid %d minutes. Do not share this code with anyone.", otp, int(memberPinOTPTTL.Minutes()))
	if err := notifyMemberSMS(m.ID, msg, 0); err != nil {
		log.Printf("[MemberPin] SMS to member %d failed: %v", m.ID, err)
		return nil, serviceErr(http.StatusBadGateway, "Could not send the SMS code. Please try again.")
	}
	memberLoginSucceeded(keys...)
	return &MemberPinOTPResponse{MaskedPhone: maskPhone(phone), ExpiresInSeconds: int(memberPinOTPTTL.Seconds()),
		ResendAfterSecond: int(memberPinOTPCooldown.Seconds())}, nil
}

func writeDevPinOTPLog(memberID uint64, memberNo, phone, otp string) {
	path := strings.TrimSpace(os.Getenv("DEV_OTP_LOG_FILE"))
	if path == "" {
		path = filepath.Join("logs", "dev-otp.log")
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s | PIN setup | member %d (%s) | phone %s | OTP %s\n", time.Now().Format("2006-01-02 15:04:05"), memberID, memberNo, phone, otp)
}

// MemberPinResetRequest is POST /api/member-pin/reset.
type MemberPinResetRequest struct {
	Phone    string `json:"phone"`
	IDNumber string `json:"id_number"`
	OTP      string `json:"otp"`
	Pin      string `json:"pin"`
}

// ResetMemberPinWithOTP checks the SMS code, sets the new PIN and signs the member in.
func ResetMemberPinWithOTP(req MemberPinResetRequest, clientIP string) (*MemberLoginResponse, error) {
	phone := utils.NormalizePhone(req.Phone)
	idNo := normalizeIDNo(req.IDNumber)
	keys := []string{"pinotp-phone:" + phone, "pinotp-ip:" + clientIP}
	if memberLoginBlocked(keys...) {
		return nil, serviceErr(http.StatusTooManyRequests, "Too many failed attempts. Please wait 15 minutes and try again.")
	}
	m := findMemberByPhoneAndID(phone, idNo)
	if m == nil {
		memberLoginFailed(keys...)
		return nil, serviceErr(http.StatusUnauthorized, "Phone number or ID number is incorrect.")
	}
	r := loadMemberPin(m.ID)
	if r.OtpHash == nil || r.OtpExpiresAt == nil || time.Now().After(*r.OtpExpiresAt) {
		return nil, serviceErr(http.StatusUnprocessableEntity, "The code has expired. Request a new code.")
	}
	if r.OtpAttempts >= memberPinOTPMaxTries {
		return nil, serviceErr(http.StatusUnprocessableEntity, "Too many wrong codes. Request a new code.")
	}
	if bcrypt.CompareHashAndPassword([]byte(*r.OtpHash), []byte(strings.TrimSpace(req.OTP))) != nil {
		db.DB.Exec(`UPDATE member_login_pins SET otp_attempts = otp_attempts + 1, updated_at = NOW() WHERE member_id = ?`, m.ID)
		return nil, serviceErr(http.StatusUnprocessableEntity, "The code is incorrect.")
	}
	if err := ValidateNewMemberPin(req.Pin, m.PrimaryPhone, m.IDNo); err != nil {
		return nil, err
	}
	if err := SetMemberPin(m.ID, req.Pin); err != nil {
		return nil, err
	}
	memberLoginSucceeded(keys...)
	log.Printf("[MemberPin] member %d (%s) set a new PIN by SMS code", m.ID, m.MemberNo)
	return issueMemberLogin(m, phone)
}

// MemberPinChangeRequest is POST /api/member-pin/change (signed-in member).
type MemberPinChangeRequest struct {
	CurrentPin string `json:"current_pin"`
	NewPin     string `json:"new_pin"`
}

func ChangeMemberPin(memberID uint64, req MemberPinChangeRequest) error {
	r := loadMemberPin(memberID)
	if r.PinHash == nil || bcrypt.CompareHashAndPassword([]byte(*r.PinHash), []byte(pinDigits(req.CurrentPin))) != nil {
		return serviceErr(http.StatusUnprocessableEntity, "The current PIN is incorrect.")
	}
	var m struct {
		PrimaryPhone string
		IDNo         string `gorm:"column:id_no"`
	}
	db.DB.Raw(`SELECT primary_phone, id_no FROM member_registrations WHERE id = ?`, memberID).Scan(&m)
	if err := ValidateNewMemberPin(req.NewPin, m.PrimaryPhone, m.IDNo); err != nil {
		return err
	}
	return SetMemberPin(memberID, req.NewPin)
}

// MemberHasPin reports whether a member has set a login PIN (for admin lists).
func MemberHasPin(memberID uint64) bool {
	var n int64
	db.DB.Raw(`SELECT COUNT(*) FROM member_login_pins WHERE member_id = ? AND pin_hash IS NOT NULL`, memberID).Scan(&n)
	return n > 0
}
