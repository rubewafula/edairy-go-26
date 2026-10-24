package services

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/rubewafula/edairy-go-26/internal/db"
)

// DefaultCashoutSupportPhone is shown on Processing cashouts until an admin sets another number.
const DefaultCashoutSupportPhone = "0705182016"

// CashoutSettings live in loan_cashout_settings (single row id = 1, created at API start).
type CashoutSettings struct {
	SupportPhone string `json:"support_phone"`
}

var cashoutSettingsOnce sync.Once

func ensureCashoutSettingsTable() error {
	var err error
	cashoutSettingsOnce.Do(func() {
		if err = db.DB.Exec(`CREATE TABLE IF NOT EXISTS loan_cashout_settings (
			id INT UNSIGNED NOT NULL PRIMARY KEY,
			support_phone VARCHAR(32) NOT NULL DEFAULT '` + DefaultCashoutSupportPhone + `',
			updated_by BIGINT UNSIGNED NULL,
			updated_at DATETIME NULL)`).Error; err != nil {
			return
		}
		err = db.DB.Exec(`INSERT IGNORE INTO loan_cashout_settings (id, support_phone, updated_at) VALUES (1, ?, NOW())`,
			DefaultCashoutSupportPhone).Error
	})
	return err
}

// GetCashoutSettings returns the stored settings (defaults if the row is missing or empty).
func GetCashoutSettings() CashoutSettings {
	s := CashoutSettings{SupportPhone: DefaultCashoutSupportPhone}
	if ensureCashoutSettingsTable() != nil {
		return s
	}
	var phone string
	db.DB.Raw(`SELECT support_phone FROM loan_cashout_settings WHERE id = 1`).Scan(&phone)
	if p := strings.TrimSpace(phone); p != "" {
		s.SupportPhone = p
	}
	return s
}

// UpdateCashoutSettings stores a new support phone (digits, spaces, + and - only).
func UpdateCashoutSettings(in CashoutSettings, userID uint64) (CashoutSettings, error) {
	phone := strings.TrimSpace(in.SupportPhone)
	digits := 0
	for _, r := range phone {
		switch {
		case r >= '0' && r <= '9':
			digits++
		case r == '+' || r == ' ' || r == '-':
		default:
			return CashoutSettings{}, serviceErr(http.StatusBadRequest, "support phone may only contain digits, spaces, + and -")
		}
	}
	if digits < 9 || digits > 15 {
		return CashoutSettings{}, serviceErr(http.StatusBadRequest, "enter a valid support phone number")
	}
	if err := ensureCashoutSettingsTable(); err != nil {
		return CashoutSettings{}, fmt.Errorf("cashout settings: %w", err)
	}
	if err := db.DB.Exec(`UPDATE loan_cashout_settings SET support_phone = ?, updated_by = ?, updated_at = NOW() WHERE id = 1`,
		phone, userID).Error; err != nil {
		return CashoutSettings{}, err
	}
	return GetCashoutSettings(), nil
}
