package services

import (
	"sort"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

// MemberLoanTransaction is one line of a member's loan statement (GET /api/loan-module/mobile/transactions).
type MemberLoanTransaction struct {
	ID         string    `json:"id"`
	Date       time.Time `json:"date"`
	Kind       string    `json:"kind"` // disbursed | cost | repaid
	Type       string    `json:"type"` // e.g. MOBILE_MPESA, PROCESSING, INSURANCE, PER_LITRE, MILK_DEDUCTION
	Label      string    `json:"label"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	ContractNo string    `json:"contract_no"`
	Product    string    `json:"product"`
	Reference  string    `json:"reference"`
}

// MemberLoanTransactionTotals sums a statement: money sent, what it cost, what was repaid, what is owed.
type MemberLoanTransactionTotals struct {
	Disbursed   float64 `json:"disbursed"`
	Cost        float64 `json:"cost"`
	Repaid      float64 `json:"repaid"`
	Outstanding float64 `json:"outstanding"`
}

var chargeLabels = map[string]string{
	"PROCESSING":       "Processing fee",
	"INSURANCE":        "Insurance fee",
	ChargeTypePerLitre: "Per-litre cashout charge",
}

func labelFor(m map[string]string, key, fallback string) string {
	if l, ok := m[strings.ToUpper(key)]; ok {
		return l
	}
	if key == "" {
		return fallback
	}
	return fallback + " (" + strings.ToLower(strings.ReplaceAll(key, "_", " ")) + ")"
}

func inRange(t time.Time, from, to string) bool {
	d := t.Format("2006-01-02")
	return (from == "" || d >= from) && (to == "" || d <= to)
}

// ListMemberLoanTransactions returns the member's loan statement lines (newest first) and their totals.
// from/to (YYYY-MM-DD, optional) filter by line date.
func (s *LoanModuleService) ListMemberLoanTransactions(memberID uint64, from, to string) ([]MemberLoanTransaction, MemberLoanTransactionTotals, error) {
	var totals MemberLoanTransactionTotals
	var contracts []models.LoanContract
	if err := db.DB.Preload("LoanProduct").Where("member_id = ?", memberID).Find(&contracts).Error; err != nil {
		return nil, totals, err
	}
	// Only money that actually moved: loans whose payout succeeded, their costs, and repayments.
	var paidOut []models.LoanDisbursement
	if len(contracts) > 0 {
		all := make([]uint64, 0, len(contracts))
		for _, c := range contracts {
			all = append(all, c.ID)
		}
		db.DB.Where("loan_contract_id IN ? AND status = ?", all, models.DisbursementStatusSuccess).Find(&paidOut)
	}
	disbursedContract := map[uint64]bool{}
	for _, d := range paidOut {
		disbursedContract[d.LoanContractID] = true
	}
	ids := make([]uint64, 0, len(contracts))
	byID := map[uint64]*models.LoanContract{}
	for i := range contracts {
		if disbursedContract[contracts[i].ID] {
			ids = append(ids, contracts[i].ID)
			byID[contracts[i].ID] = &contracts[i]
		}
	}
	if len(ids) == 0 {
		return []MemberLoanTransaction{}, totals, nil
	}
	product := func(c *models.LoanContract) string {
		if c != nil && c.LoanProduct != nil {
			return c.LoanProduct.ProductName
		}
		return ""
	}
	var lines []MemberLoanTransaction
	add := func(t MemberLoanTransaction) {
		if inRange(t.Date, from, to) {
			lines = append(lines, t)
		}
	}

	for _, d := range paidOut {
		c := byID[d.LoanContractID]
		date := d.DisbursementDate
		if d.CompletedAt != nil {
			date = *d.CompletedAt
		}
		label := "Loan disbursed"
		if d.ChannelCode == models.DisbursementCodeMobileMpesa {
			label += " via M-Pesa"
		}
		add(MemberLoanTransaction{ID: "D" + itoa(d.ID), Date: date, Kind: "disbursed", Type: d.ChannelCode,
			Label: label, Amount: d.Amount, Status: d.Status, ContractNo: c.ContractNo, Product: product(c),
			Reference: firstNonEmptyString(d.ExternalReference, d.Reference)})
	}

	var charges []models.LoanChargeRecord
	db.DB.Where("loan_contract_id IN ?", ids).Find(&charges)
	for _, ch := range charges {
		c := byID[ch.LoanContractID]
		add(MemberLoanTransaction{ID: "C" + itoa(ch.ID), Date: ch.ChargedDate, Kind: "cost", Type: ch.ChargeType,
			Label: labelFor(chargeLabels, ch.ChargeType, "Loan charge"), Amount: ch.Amount, Status: "CHARGED",
			ContractNo: c.ContractNo, Product: product(c)})
	}

	// Interest scheduled on each contract, shown once on the contract's date.
	type interestRow struct {
		LoanContractID uint64
		Interest       float64
	}
	var interest []interestRow
	db.DB.Raw(`SELECT loan_contract_id, SUM(interest_due) AS interest FROM loan_schedule_installments
		WHERE loan_contract_id IN ? AND deleted_at IS NULL GROUP BY loan_contract_id HAVING SUM(interest_due) > 0`, ids).Scan(&interest)
	for _, r := range interest {
		c := byID[r.LoanContractID]
		date := c.CreatedAt
		if c.DisbursementDate != nil {
			date = *c.DisbursementDate
		}
		add(MemberLoanTransaction{ID: "I" + itoa(c.ID), Date: date, Kind: "cost", Type: "INTEREST",
			Label: "Interest on loan", Amount: roundMoney(r.Interest), Status: c.Status, ContractNo: c.ContractNo, Product: product(c)})
	}

	var penalties []models.LoanPenaltyRecord
	db.DB.Where("loan_contract_id IN ? AND waived = 0", ids).Find(&penalties)
	for _, p := range penalties {
		c := byID[p.LoanContractID]
		add(MemberLoanTransaction{ID: "P" + itoa(p.ID), Date: p.PenaltyDate, Kind: "cost", Type: p.PenaltyType,
			Label: labelFor(map[string]string{}, p.PenaltyType, "Late penalty"), Amount: p.Amount, Status: "CHARGED",
			ContractNo: c.ContractNo, Product: product(c)})
	}

	var repayments []models.LoanRepaymentRecord
	db.DB.Where("loan_contract_id IN ?", ids).Find(&repayments)
	repayLabels := map[string]string{
		models.LoanChannelMilkDeduction: "Loan repaid from milk",
		models.LoanChannelMobileMoney:   "Loan repaid via mobile money",
		models.LoanChannelCash:          "Loan repaid in cash",
		models.LoanChannelBank:          "Loan repaid via bank",
	}
	for _, r := range repayments {
		c := byID[r.LoanContractID]
		add(MemberLoanTransaction{ID: "R" + itoa(r.ID), Date: r.PaymentDate, Kind: "repaid", Type: r.Channel,
			Label: labelFor(repayLabels, r.Channel, "Loan repaid"), Amount: r.Amount, Status: "PAID",
			ContractNo: c.ContractNo, Product: product(c), Reference: r.Reference})
	}

	sort.SliceStable(lines, func(i, j int) bool { return lines[i].Date.After(lines[j].Date) })
	for _, l := range lines {
		switch l.Kind {
		case "disbursed":
			totals.Disbursed += l.Amount
		case "cost":
			totals.Cost += l.Amount
		case "repaid":
			totals.Repaid += l.Amount
		}
	}
	totals.Disbursed = roundMoney(totals.Disbursed)
	totals.Cost = roundMoney(totals.Cost)
	totals.Repaid = roundMoney(totals.Repaid)
	// What the member still owes on all their loans (not limited by the date filter).
	var owed float64
	for _, id := range ids {
		owed += contractRemainingRepayable(db.DB, id)
	}
	totals.Outstanding = roundMoney(owed)
	if lines == nil {
		lines = []MemberLoanTransaction{}
	}
	return lines, totals, nil
}

func itoa(v uint64) string {
	const digits = "0123456789"
	if v == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = digits[v%10]
		v /= 10
	}
	return string(b[i:])
}

func firstNonEmptyString(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
