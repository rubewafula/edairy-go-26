package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// Member self-service: farmers sign in to the mobile app with their phone number and national ID
// number. Each member gets a hidden users row (users.member_id = member_registrations.id) with the
// "Member" role, so the normal JWT/auth middleware works. middleware.MemberSelfServiceScope then
// limits that login to the member's own data and a short list of endpoints.

const MemberRoleName = "Member"

// MemberLoginRequest is POST /api/member-login.
type MemberLoginRequest struct {
	Phone    string `json:"phone"`
	IDNumber string `json:"id_number"`
}

type MemberLoginProfile struct {
	ID       uint64 `json:"id"`
	MemberNo string `json:"member_no"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	IDNumber string `json:"id_number"`
	Status   string `json:"status"`
}

type MemberLoginResponse struct {
	Token       string             `json:"token"`
	UserID      uint64             `json:"user_id"`
	Email       string             `json:"email"`
	Roles       []string           `json:"roles"`
	Permissions []string           `json:"permissions"`
	AccountType string             `json:"account_type"` // "member"
	Member      MemberLoginProfile `json:"member"`
}

var memberSelfServiceOnce sync.Once

// EnsureMemberSelfServiceSchema adds users.member_id and the Member role (idempotent).
func EnsureMemberSelfServiceSchema() error {
	var err error
	memberSelfServiceOnce.Do(func() {
		var n int64
		db.DB.Raw(`SELECT COUNT(*) FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'member_id'`).Scan(&n)
		if n == 0 {
			if err = db.DB.Exec(`ALTER TABLE users ADD COLUMN member_id BIGINT UNSIGNED NULL, ADD INDEX idx_users_member_id (member_id)`).Error; err != nil {
				return
			}
			log.Printf("[MemberLogin] added users.member_id")
		}
		err = db.DB.Exec(`INSERT INTO roles (name, guard_name, created_at, updated_at)
			SELECT ?, 'web', NOW(), NOW() FROM DUAL
			WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = ? AND deleted_at IS NULL)`, MemberRoleName, MemberRoleName).Error
	})
	return err
}

// ---- brute-force guard: 5 failed attempts per phone or per IP in 15 minutes ----

type loginAttempts struct {
	count int
	first time.Time
}

var (
	memberLoginMu       sync.Mutex
	memberLoginFailures = map[string]*loginAttempts{}
)

const (
	memberLoginMaxFailures = 5
	memberLoginWindow      = 15 * time.Minute
)

func memberLoginBlocked(keys ...string) bool {
	memberLoginMu.Lock()
	defer memberLoginMu.Unlock()
	for _, k := range keys {
		if a := memberLoginFailures[k]; a != nil {
			if time.Since(a.first) > memberLoginWindow {
				delete(memberLoginFailures, k)
				continue
			}
			if a.count >= memberLoginMaxFailures {
				return true
			}
		}
	}
	return false
}

func memberLoginFailed(keys ...string) {
	memberLoginMu.Lock()
	defer memberLoginMu.Unlock()
	for _, k := range keys {
		a := memberLoginFailures[k]
		if a == nil || time.Since(a.first) > memberLoginWindow {
			a = &loginAttempts{first: time.Now()}
			memberLoginFailures[k] = a
		}
		a.count++
	}
}

func memberLoginSucceeded(keys ...string) {
	memberLoginMu.Lock()
	defer memberLoginMu.Unlock()
	for _, k := range keys {
		delete(memberLoginFailures, k)
	}
}

// MemberLogin checks phone (2547…, 07…, 7…, +254…) + ID number against member_registrations.
func MemberLogin(req MemberLoginRequest, clientIP string) (*MemberLoginResponse, error) {
	phone := utils.NormalizePhone(req.Phone)
	idNo := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(req.IDNumber), " ", ""))
	if phone == "" {
		return nil, serviceErr(http.StatusBadRequest, "Enter a valid phone number, e.g. 0712345678, 712345678 or 254712345678.")
	}
	if idNo == "" {
		return nil, serviceErr(http.StatusBadRequest, "Enter your ID number.")
	}
	keys := []string{"phone:" + phone, "ip:" + clientIP}
	if memberLoginBlocked(keys...) {
		return nil, serviceErr(http.StatusTooManyRequests, "Too many failed attempts. Please wait 15 minutes and try again.")
	}
	if err := EnsureMemberSelfServiceSchema(); err != nil {
		return nil, fmt.Errorf("member login setup: %w", err)
	}

	var candidates []struct {
		ID           uint64
		MemberNo     string
		FirstName    string
		LastName     string
		PrimaryPhone string
		IDNo         string `gorm:"column:id_no"`
		Email        string
		Status       string
	}
	db.DB.Raw(`SELECT id, member_no, first_name, last_name, primary_phone, id_no, email, status
		FROM member_registrations
		WHERE deleted_at IS NULL AND UPPER(REPLACE(TRIM(id_no), ' ', '')) = ?`, idNo).Scan(&candidates)

	var match *MemberLoginProfile
	var memberEmail string
	for _, m := range candidates {
		if utils.NormalizePhone(m.PrimaryPhone) == phone {
			match = &MemberLoginProfile{
				ID: m.ID, MemberNo: m.MemberNo, Phone: phone, IDNumber: m.IDNo, Status: m.Status,
				Name: strings.TrimSpace(m.FirstName + " " + m.LastName),
			}
			memberEmail = strings.TrimSpace(m.Email)
			break
		}
	}
	if match == nil {
		memberLoginFailed(keys...)
		return nil, serviceErr(http.StatusUnauthorized, "Phone number or ID number is incorrect. Use the phone number registered with the dairy.")
	}
	switch strings.ToUpper(strings.TrimSpace(match.Status)) {
	case "EXITED", "DECEASED", "DELETED", "CLOSED", "REJECTED":
		return nil, serviceErr(http.StatusForbidden, "This member account is %s. Please contact the dairy office.", strings.ToLower(match.Status))
	}
	memberLoginSucceeded(keys...)

	userID, email, err := ensureMemberUser(match.ID, match.Name, memberEmail)
	if err != nil {
		return nil, fmt.Errorf("member login: %w", err)
	}
	token, err := memberToken(userID, email)
	if err != nil {
		return nil, err
	}
	forgetMemberUserCache(userID)
	log.Printf("[MemberLogin] member %d (%s) signed in as user %d", match.ID, match.MemberNo, userID)
	return &MemberLoginResponse{
		Token: token, UserID: userID, Email: email,
		Roles: []string{MemberRoleName}, Permissions: []string{},
		AccountType: "member", Member: *match,
	}, nil
}

// ensureMemberUser finds or creates the hidden login user for a member.
func ensureMemberUser(memberID uint64, name, memberEmail string) (uint64, string, error) {
	var existing struct {
		ID    uint64
		Email string
	}
	db.DB.Raw(`SELECT id, email FROM users WHERE member_id = ? AND deleted_at IS NULL ORDER BY id LIMIT 1`, memberID).Scan(&existing)
	if existing.ID == 0 {
		email := fmt.Sprintf("member-%d@members.edairy.local", memberID)
		// Use the member's own email only if no other login already has it (staff accounts keep theirs).
		if memberEmail != "" {
			var taken int64
			db.DB.Raw(`SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER(?)`, memberEmail).Scan(&taken)
			if taken == 0 {
				email = memberEmail
			}
		}
		buf := make([]byte, 24)
		if _, err := rand.Read(buf); err != nil {
			return 0, "", err
		}
		// Random password: this user can only sign in through /api/member-login.
		hash, err := bcrypt.GenerateFromPassword([]byte(hex.EncodeToString(buf)), 10)
		if err != nil {
			return 0, "", err
		}
		user := models.User{Name: name, Email: email, Password: string(hash), IsVerified: true}
		if err := db.DB.Create(&user).Error; err != nil {
			return 0, "", err
		}
		if err := db.DB.Exec(`UPDATE users SET member_id = ? WHERE id = ?`, memberID, user.ID).Error; err != nil {
			return 0, "", err
		}
		existing.ID, existing.Email = user.ID, user.Email
	}
	var role models.Role
	if err := db.DB.Where("name = ?", MemberRoleName).First(&role).Error; err == nil {
		var has int64
		db.DB.Raw(`SELECT COUNT(*) FROM user_roles WHERE user_id = ? AND role_id = ?`, existing.ID, role.ID).Scan(&has)
		if has == 0 {
			u := models.User{}
			u.ID = existing.ID
			if err := db.DB.Model(&u).Association("Roles").Append(&role); err != nil {
				log.Printf("[MemberLogin] could not give user %d the Member role: %v", existing.ID, err)
			}
		}
	}
	return existing.ID, existing.Email, nil
}

// memberToken: same claims as staff tokens; lasts MEMBER_JWT_EXPIRY_HOURS (default 30 days).
func memberToken(userID uint64, email string) (string, error) {
	hours, _ := strconv.Atoi(os.Getenv("MEMBER_JWT_EXPIRY_HOURS"))
	if hours <= 0 {
		hours = 24 * 30
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"roles":   []string{MemberRoleName},
		"member":  true,
		"exp":     utils.Now().Add(time.Duration(hours) * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ---- users.member_id lookup (cached) ----

type cachedMemberUser struct {
	memberID uint64
	at       time.Time
}

var memberUserCache sync.Map

// MemberIDForUser returns the member a login belongs to (0 for staff).
func MemberIDForUser(userID uint64) uint64 {
	if v, ok := memberUserCache.Load(userID); ok {
		if c := v.(cachedMemberUser); time.Since(c.at) < time.Minute {
			return c.memberID
		}
	}
	var memberID uint64
	if EnsureMemberSelfServiceSchema() == nil {
		db.DB.Raw(`SELECT COALESCE(member_id, 0) FROM users WHERE id = ? AND deleted_at IS NULL`, userID).Scan(&memberID)
	}
	memberUserCache.Store(userID, cachedMemberUser{memberID: memberID, at: time.Now()})
	return memberID
}

func forgetMemberUserCache(userID uint64) { memberUserCache.Delete(userID) }

// LoanApplicationMemberID is the member a loan application belongs to (0 if not found).
func LoanApplicationMemberID(applicationID uint64) uint64 {
	var memberID uint64
	db.DB.Raw(`SELECT member_id FROM loan_applications WHERE id = ? AND deleted_at IS NULL`, applicationID).Scan(&memberID)
	return memberID
}

// MemberSelfProfile is GET /api/member-me.
func MemberSelfProfile(memberID uint64) (*MemberLoginProfile, error) {
	var m struct {
		ID           uint64
		MemberNo     string
		FirstName    string
		LastName     string
		PrimaryPhone string
		IDNo         string `gorm:"column:id_no"`
		Status       string
	}
	db.DB.Raw(`SELECT id, member_no, first_name, last_name, primary_phone, id_no, status
		FROM member_registrations WHERE id = ? AND deleted_at IS NULL`, memberID).Scan(&m)
	if m.ID == 0 {
		return nil, serviceErr(http.StatusNotFound, "member not found")
	}
	return &MemberLoginProfile{ID: m.ID, MemberNo: m.MemberNo, Name: strings.TrimSpace(m.FirstName + " " + m.LastName),
		Phone: utils.NormalizePhone(m.PrimaryPhone), IDNumber: m.IDNo, Status: m.Status}, nil
}
