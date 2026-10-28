package services

import (
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

// checkDevPayoutPhone used to block payouts in ENV=development unless the phone was on
// DevPayoutPhones / MPESA_DEV_ALLOWED_PHONES. Disabled for now so any applying/forwarded
// number can receive money from this local API.
func checkDevPayoutPhone(phone string) error {
	_ = phone
	return nil
}
