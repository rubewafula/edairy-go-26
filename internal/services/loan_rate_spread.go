package services

import (
	"math"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// Member loans (advances) are priced on the milk-rate spread instead of a product interest rate:
// the member borrows against milk valued at the CREDIT rate (default_milk_rates.advance_rate) and
// repays the same milk at the DEFAULT rate (default_milk_rates.rate, the payout rate). So
//   repayable = borrowed x default / credit, interest = repayable - borrowed.
// In credit terms one shilling owed from milk uses credit/default of a shilling of credit.

// MemberRates are the two rates that apply to a member today.
type MemberRates struct {
	CreditRate  float64 `json:"credit_rate"`
	DefaultRate float64 `json:"default_rate"`
}

// Ratio is credit/default (0 when either rate is missing).
func (r MemberRates) Ratio() float64 {
	if r.CreditRate <= 0 || r.DefaultRate <= 0 {
		return 0
	}
	return math.Min(1, r.CreditRate/r.DefaultRate)
}

func currentPayDateRangeID(tx *gorm.DB) uint64 {
	var id uint64
	tx.Raw(`SELECT id FROM member_pay_date_ranges WHERE ? BETWEEN start_date AND end_date AND deleted_at IS NULL
		ORDER BY id DESC LIMIT 1`, time.Now().Format("2006-01-02")).Scan(&id)
	return id
}

var memberRatesCache sync.Map // memberID -> cachedRates

type cachedRates struct {
	rates MemberRates
	at    time.Time
}

// GetMemberRates resolves the member's credit and default milk rates (cached for a few seconds:
// quotes and the max-request search ask for them many times in one request).
func GetMemberRates(tx *gorm.DB, memberID uint64) MemberRates {
	if v, ok := memberRatesCache.Load(memberID); ok {
		if c := v.(cachedRates); time.Since(c.at) < 5*time.Second {
			return c.rates
		}
	}
	r := loadMemberRates(tx, memberID)
	memberRatesCache.Store(memberID, cachedRates{rates: r, at: time.Now()})
	return r
}

// ForgetMemberRates drops cached rates (after a rate change).
func ForgetMemberRates() {
	memberRatesCache.Range(func(k, _ interface{}) bool { memberRatesCache.Delete(k); return true })
}

func loadMemberRates(tx *gorm.DB, memberID uint64) MemberRates {
	var routeID uint64
	tx.Raw(`SELECT COALESCE(route_id, 0) FROM member_registrations WHERE id = ?`, memberID).Scan(&routeID)
	out := MemberRates{CreditRate: MemberAdvanceRate(tx, memberID)}
	pdr := currentPayDateRangeID(tx)
	var pdrs []uint64
	if pdr > 0 {
		pdrs = []uint64{pdr}
	}
	if rates, err := loadMilkRateResolver(tx, pdrs); err == nil {
		out.DefaultRate = rates.Rate(pdr, memberID, routeID)
	}
	return out
}

// spreadInterest is the interest on a member-loan advance; ok=false for other products.
func spreadInterest(tx *gorm.DB, product *models.LoanProduct, memberID uint64, amount float64) (float64, bool) {
	if product == nil || !product.IsDefaultMemberLoan {
		return 0, false
	}
	ratio := GetMemberRates(tx, memberID).Ratio()
	if ratio <= 0 || ratio >= 1 {
		return 0, true
	}
	return roundMoney(amount/ratio - amount), true
}
