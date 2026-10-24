package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/utils"
)

// defaultDevPayoutPhones are the only numbers a development API may send real money to (M-Pesa, DTB,
// Equitel, Airtel). Override with MPESA_DEV_ALLOWED_PHONES (comma-separated, any Kenyan format).
var defaultDevPayoutPhones = []string{"254705182016", "254726986944", "254721695524"}

// DevPayoutPhones returns the normalized allow-list used when ENV=development.
func DevPayoutPhones() []string {
	raw := defaultDevPayoutPhones
	if v := strings.TrimSpace(os.Getenv("MPESA_DEV_ALLOWED_PHONES")); v != "" {
		raw = strings.Split(v, ",")
	}
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		if n, err := utils.NormalizeKenyaPhone(p, ""); err == nil {
			out = append(out, n)
		}
	}
	return out
}

// checkDevPayoutPhone blocks a payout in ENV=development unless the phone is on the allow-list.
// The error text contains "payouts are disabled" so the disbursement queue treats it as a
// configuration problem: the loan waits (retried every 5 min, no attempts used) instead of failing.
func checkDevPayoutPhone(phone string) error {
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("ENV")), "development") {
		return nil
	}
	n, err := utils.NormalizeKenyaPhone(phone, "")
	if err != nil {
		return fmt.Errorf("real payouts are disabled in development for invalid phone %q", phone)
	}
	allowed := DevPayoutPhones()
	for _, a := range allowed {
		if a == n {
			return nil
		}
	}
	return fmt.Errorf("real payouts are disabled in development for %s: only %s may receive money (MPESA_DEV_ALLOWED_PHONES)",
		n, strings.Join(allowed, ", "))
}
