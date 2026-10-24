package services

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const creditCustomerTypeMember = "member"

// Contract statuses whose remaining schedule still has to be recovered from milk.
var creditObligationContractStatuses = []string{
	models.LoanContractPending,
	models.LoanContractActive,
	models.LoanContractOverdue,
	models.LoanContractRestructured,
}

type MemberCreditService struct {
	// excludeJournalIDs leaves these journals' milk out of Compute (used by the one-off migration only).
	excludeJournalIDs []uint64
}

func NewMemberCreditService() *MemberCreditService {
	return &MemberCreditService{}
}

// MemberCreditBreakdown explains how a member's credit limit was derived.
type MemberCreditBreakdown struct {
	MemberID        uint64     `json:"member_id"`
	MilkValueUnpaid float64    `json:"milk_value_unpaid"`
	RejectValue     float64    `json:"reject_value"`
	LoanObligations float64    `json:"loan_obligations"`
	OtherDeductions float64    `json:"other_deductions"`
	CreditLimit     float64    `json:"credit_limit"`
	UpdatedAt       *time.Time `json:"updated_at"`
	// MilkKilos is the confirmed, unpaid milk counted; ZeroRateKilos is the part with no milk rate (worth 0).
	MilkKilos     float64 `json:"milk_kilos"`
	ZeroRateKilos float64 `json:"zero_rate_kilos,omitempty"`
	// Two-rate view: CreditLimit is in credit-rate terms (what can be borrowed); DefaultBasisLimit is
	// the same balance at the default (payout) rate. MilkValueDefault is the milk at the default rate.
	MilkValueDefault  float64 `json:"milk_value_default"`
	DefaultBasisLimit float64 `json:"default_basis_limit"`
	CreditRate        float64 `json:"credit_rate"`
	DefaultRate       float64 `json:"default_rate"`
}

// CreditChange describes the event that triggered a recalculation.
type CreditChange struct {
	Action     string
	SourceType string
	SourceID   uint64
	Reference  string
	Narration  string
	UserID     uint64
}

func ensureMemberCreditSchema() error {
	m := db.DB.Migrator()
	for _, field := range []string{"CreditLimit", "CreditLimitUpdatedAt", "CreditLimitComputed", "CanBorrow", "CannotBorrowCode", "CannotBorrowReason", "MemberBalance"} {
		if !m.HasColumn(&models.Member{}, field) {
			if err := m.AddColumn(&models.Member{}, field); err != nil {
				return fmt.Errorf("add member_registrations.%s: %w", field, err)
			}
		}
	}

	logModel := &models.CreditLimitChangeLog{}
	if !m.HasTable(logModel) {
		return m.CreateTable(logModel)
	}
	for _, field := range []string{"Action", "Amount", "Reference", "SourceType", "SourceID", "Narration", "CreatedBy", "UpdatedBy", "CreatedAt", "UpdatedAt", "DeletedAt"} {
		if !m.HasColumn(logModel, field) {
			if err := m.AddColumn(logModel, field); err != nil {
				return fmt.Errorf("add credit_limit_change_log.%s: %w", field, err)
			}
		}
	}
	columns, err := m.ColumnTypes(logModel)
	if err != nil {
		return err
	}
	for _, col := range columns {
		name := col.Name()
		if name != "old_limit" && name != "credit_limit" {
			continue
		}
		if strings.EqualFold(col.DatabaseTypeName(), "DECIMAL") {
			continue
		}
		field := "OldLimit"
		if name == "credit_limit" {
			field = "CreditLimit"
		}
		if err := m.AlterColumn(logModel, field); err != nil {
			return fmt.Errorf("alter credit_limit_change_log.%s: %w", name, err)
		}
	}
	return nil
}

// Permission names follow middleware.RequireAutoPermission, which derives them from the route path and method.
var memberCreditPermissions = []string{
	"member-credit-limit.view",
	"members.credit-limit.view",
	"members.credit-limit-logs.view",
	"members.credit-limit.recalculate.create",
	"member-credit-limits.recalculate-all.create",
	"member-credit-limits.reconcile.create",
	"loan-module.mobile.loan-requests.create",
	"loan-module.mobile.loan-requests.view",
	"loan-module.mobile.loan-requests.resend-otp.create",
	"loan-module.mobile.loan-requests.verify-otp.create",
	"loan-module.mobile.transactions.view",
	"loan-module.mobile.quote.view",
	"loan-module.cashout-settings.view",
	"loan-module.cashout-settings.update",
	"members.reset-pin.create",
	"members.pin-status.view",
}

func ensureMemberCreditPermissions() error {
	for _, name := range memberCreditPermissions {
		if err := db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
			SELECT ?, 'web', NOW(), NOW() FROM DUAL
			WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = ?)`, name, name).Error; err != nil {
			return fmt.Errorf("seed permission %s: %w", name, err)
		}
	}
	if err := db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r
		JOIN permissions p ON p.name IN ? AND p.deleted_at IS NULL
		WHERE r.name = 'Super Admin'
		AND NOT EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_id = r.id AND rp.permission_id = p.id)`,
		memberCreditPermissions).Error; err != nil {
		return err
	}
	// Anyone who can see a member's loans/cashouts can also see their loan transactions.
	return db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT DISTINCT rp.role_id, tp.id FROM role_permissions rp
		JOIN permissions src ON src.id = rp.permission_id AND src.deleted_at IS NULL
			AND src.name IN ('loan-module.mobile.loan-requests.view', 'loan-module.applications.view', 'member-credit-limit.view')
		JOIN permissions tp ON tp.name IN ('loan-module.mobile.transactions.view', 'loan-module.mobile.quote.view') AND tp.deleted_at IS NULL
		WHERE rp.deleted_at IS NULL
		AND NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = rp.role_id AND x.permission_id = tp.id)`).Error
}

// milkRateResolver applies the same precedence as member payroll generation:
// member special rate for the period, member default, route special for the period, route default, global default.
type milkRateResolver struct {
	specialMember map[uint64]map[uint64]float64
	specialRoute  map[uint64]map[uint64]float64
	defaultMember map[uint64]float64
	defaultRoute  map[uint64]float64
	global        float64
}

func newMilkRateResolver() *milkRateResolver {
	return &milkRateResolver{
		specialMember: map[uint64]map[uint64]float64{},
		specialRoute:  map[uint64]map[uint64]float64{},
		defaultMember: map[uint64]float64{},
		defaultRoute:  map[uint64]float64{},
	}
}

func loadMilkRateResolver(tx *gorm.DB, payDateRangeIDs []uint64) (*milkRateResolver, error) {
	r := newMilkRateResolver()
	if len(payDateRangeIDs) > 0 {
		var special []models.MilkSpecialRate
		if err := tx.Where("pay_date_range_id IN ? AND deleted_at IS NULL", payDateRangeIDs).Find(&special).Error; err != nil {
			return nil, err
		}
		for _, sr := range special {
			r.addSpecial(sr.MonthlyPayDateRangeID, sr.MemberID, sr.RouteID, sr.Rate)
		}
	}
	var defaults []models.DefaultMilkRate
	if err := tx.Where("deleted_at IS NULL").Find(&defaults).Error; err != nil {
		return nil, err
	}
	for _, dr := range defaults {
		r.addDefault(dr.MemberID, dr.RouteID, dr.Rate)
	}
	return r, nil
}

func (r *milkRateResolver) addSpecial(payDateRangeID uint64, memberID, routeID *uint64, rate float64) {
	memberID, routeID = nonZeroID(memberID), nonZeroID(routeID)
	switch {
	case memberID != nil:
		if r.specialMember[payDateRangeID] == nil {
			r.specialMember[payDateRangeID] = map[uint64]float64{}
		}
		r.specialMember[payDateRangeID][*memberID] = rate
	case routeID != nil:
		if r.specialRoute[payDateRangeID] == nil {
			r.specialRoute[payDateRangeID] = map[uint64]float64{}
		}
		r.specialRoute[payDateRangeID][*routeID] = rate
	}
}

func (r *milkRateResolver) addDefault(memberID, routeID *uint64, rate float64) {
	memberID, routeID = nonZeroID(memberID), nonZeroID(routeID)
	switch {
	case memberID != nil:
		r.defaultMember[*memberID] = rate
	case routeID != nil:
		r.defaultRoute[*routeID] = rate
	default:
		r.global = rate
	}
}

func (r *milkRateResolver) Rate(payDateRangeID, memberID, routeID uint64) float64 {
	if rate, ok := r.specialMember[payDateRangeID][memberID]; ok {
		return rate
	}
	if rate, ok := r.defaultMember[memberID]; ok {
		return rate
	}
	if rate, ok := r.specialRoute[payDateRangeID][routeID]; ok {
		return rate
	}
	if rate, ok := r.defaultRoute[routeID]; ok {
		return rate
	}
	return r.global
}

// nonZeroID treats a stored 0 like NULL (rate for all members / all routes).
func nonZeroID(id *uint64) *uint64 {
	if id == nil || *id == 0 {
		return nil
	}
	return id
}

// RateForJournal is the member's rate, falling back to the journal's route when the member's own
// route has no rate (e.g. the member is registered on another route, or on none).
func (r *milkRateResolver) RateForJournal(payDateRangeID, memberID, memberRouteID, journalRouteID uint64) float64 {
	rate := r.Rate(payDateRangeID, memberID, memberRouteID)
	if rate <= 0 && journalRouteID > 0 && journalRouteID != memberRouteID {
		rate = r.Rate(payDateRangeID, memberID, journalRouteID)
	}
	return rate
}

// Describe lists the loaded rates, for error messages and logs.
func (r *milkRateResolver) Describe() string {
	parts := []string{fmt.Sprintf("all members %.2f", r.global)}
	for id, v := range r.defaultRoute {
		parts = append(parts, fmt.Sprintf("route %d %.2f", id, v))
	}
	if n := len(r.defaultMember); n > 0 {
		parts = append(parts, fmt.Sprintf("%d member default rate(s)", n))
	}
	for pdr, m := range r.specialRoute {
		for id, v := range m {
			parts = append(parts, fmt.Sprintf("special route %d period %d %.2f", id, pdr, v))
		}
	}
	return strings.Join(parts, ", ")
}

type periodKilos struct {
	PdrID   uint64
	RouteID uint64 // journal's route
	Kilos   float64
}

const payDateRangeForDate = `COALESCE((SELECT r.id FROM member_pay_date_ranges r
	WHERE %[1]s BETWEEN r.start_date AND r.end_date AND r.deleted_at IS NULL
	ORDER BY r.id DESC LIMIT 1), 0)`

const notInApprovedPayroll = `NOT EXISTS (SELECT 1 FROM member_payrolls mp
	JOIN member_pay_date_ranges pr ON pr.id = mp.pay_date_range_id
	WHERE mp.status = 'approved' AND mp.deleted_at IS NULL
	AND %[1]s BETWEEN pr.start_date AND pr.end_date)`

// unpaidMilkKilos sums the member's milk on CONFIRMED journals that is not yet in an approved payroll.
// Confirming a milk journal is what raises the credit limit (unconfirming takes it back).
func (s *MemberCreditService) unpaidMilkKilos(tx *gorm.DB, memberID uint64) ([]periodKilos, error) {
	exclude := ""
	args := []interface{}{memberID}
	if len(s.excludeJournalIDs) > 0 {
		exclude = " AND j.id NOT IN ?"
		args = append(args, s.excludeJournalIDs)
	}
	query := fmt.Sprintf(`SELECT t.pdr_id, t.route_id, SUM(t.quantity) AS kilos FROM (
		SELECT e.quantity, `+payDateRangeForDate+` AS pdr_id, COALESCE(j.route_id, 0) AS route_id
		FROM milk_journal_entries e
		JOIN milk_journals j ON j.id = e.milk_journal_id
		WHERE e.member_id = ? AND e.deleted_at IS NULL AND j.deleted_at IS NULL AND j.confirmed = 1`+exclude+`
		AND `+notInApprovedPayroll+`
	) t GROUP BY t.pdr_id, t.route_id`, "j.journal_date")
	var rows []periodKilos
	err := tx.Raw(query, args...).Scan(&rows).Error
	return rows, err
}

func (s *MemberCreditService) unpaidRejectKilos(tx *gorm.DB, memberID uint64) ([]periodKilos, error) {
	query := fmt.Sprintf(`SELECT t.pdr_id, SUM(t.quantity) AS kilos FROM (
		SELECT mr.quantity, `+payDateRangeForDate+` AS pdr_id
		FROM milk_rejects mr
		WHERE mr.member_id = ? AND mr.deleted_at IS NULL
		AND `+notInApprovedPayroll+`
	) t GROUP BY t.pdr_id`, "mr.transaction_date")
	var rows []periodKilos
	err := tx.Raw(query, memberID).Scan(&rows).Error
	return rows, err
}

func (s *MemberCreditService) loanObligations(tx *gorm.DB, memberID uint64) (float64, error) {
	var total float64
	err := tx.Raw(`SELECT COALESCE(SUM(GREATEST(
			(i.principal_due + i.interest_due + i.fee_due + i.insurance_due + i.penalty_due) -
			(i.principal_paid + i.interest_paid + i.fee_paid + i.insurance_paid + i.penalty_paid), 0)), 0)
		FROM loan_schedule_installments i
		JOIN loan_contracts c ON c.id = i.loan_contract_id
		WHERE c.member_id = ? AND c.status IN ? AND c.deleted_at IS NULL AND i.deleted_at IS NULL`,
		memberID, creditObligationContractStatuses).Scan(&total).Error
	return total, err
}

func (s *MemberCreditService) otherDeductions(tx *gorm.DB, memberID uint64) (float64, error) {
	var total float64
	err := tx.Raw(`SELECT COALESCE(SUM(GREATEST(rd.total_amount - rd.paid_amount, 0)), 0)
		FROM recurrent_deductions rd
		WHERE rd.customer_id = ? AND rd.customer_type = ? AND rd.settled = 0 AND rd.deleted_at IS NULL
		AND NOT EXISTS (SELECT 1 FROM loan_contracts c WHERE CONCAT('LOAN-', c.contract_no) = rd.reference)`,
		memberID, creditCustomerTypeMember).Scan(&total).Error
	return total, err
}

// Compute derives the credit limit from source data without writing anything.
func (s *MemberCreditService) Compute(tx *gorm.DB, memberID, routeID uint64) (*MemberCreditBreakdown, error) {
	milk, err := s.unpaidMilkKilos(tx, memberID)
	if err != nil {
		return nil, fmt.Errorf("unpaid milk: %w", err)
	}
	rejects, err := s.unpaidRejectKilos(tx, memberID)
	if err != nil {
		return nil, fmt.Errorf("unpaid rejects: %w", err)
	}

	pdrSet := map[uint64]bool{}
	for _, row := range append(append([]periodKilos{}, milk...), rejects...) {
		if row.PdrID > 0 {
			pdrSet[row.PdrID] = true
		}
	}
	pdrIDs := make([]uint64, 0, len(pdrSet))
	for id := range pdrSet {
		pdrIDs = append(pdrIDs, id)
	}
	// Advances are valued at the CREDIT rate (Default Milk Rates -> Advance Rate); the DEFAULT rate
	// (payout) gives the same balance in milk-payment terms.
	rates, err := loadAdvanceRateResolver(tx)
	if err != nil {
		return nil, fmt.Errorf("advance rates: %w", err)
	}
	defRates, err := loadMilkRateResolver(tx, pdrIDs)
	if err != nil {
		return nil, fmt.Errorf("milk rates: %w", err)
	}

	b := &MemberCreditBreakdown{MemberID: memberID}
	var rejectDefault float64
	for _, row := range milk {
		rate := rates.RateForJournal(row.PdrID, memberID, routeID, row.RouteID)
		b.MilkKilos += row.Kilos
		if rate <= 0 {
			b.ZeroRateKilos += row.Kilos
		}
		b.MilkValueUnpaid += row.Kilos * rate
		b.MilkValueDefault += row.Kilos * defRates.RateForJournal(row.PdrID, memberID, routeID, row.RouteID)
	}
	for _, row := range rejects {
		b.RejectValue += row.Kilos * rates.Rate(row.PdrID, memberID, routeID)
		rejectDefault += row.Kilos * defRates.Rate(row.PdrID, memberID, routeID)
	}
	mr := GetMemberRates(tx, memberID)
	b.CreditRate, b.DefaultRate = mr.CreditRate, mr.DefaultRate
	ratio := mr.Ratio()
	if b.LoanObligations, err = s.loanObligations(tx, memberID); err != nil {
		return nil, fmt.Errorf("loan obligations: %w", err)
	}
	if b.OtherDeductions, err = s.otherDeductions(tx, memberID); err != nil {
		return nil, fmt.Errorf("other deductions: %w", err)
	}

	b.MilkValueUnpaid = roundMoney(b.MilkValueUnpaid)
	b.RejectValue = roundMoney(b.RejectValue)
	b.LoanObligations = roundMoney(b.LoanObligations)
	b.OtherDeductions = roundMoney(b.OtherDeductions)
	b.MilkValueDefault = roundMoney(b.MilkValueDefault)
	// Member loans: the Total Deductible (principal + spread + fees, what milk payroll recovers) comes
	// off the credit limit shilling for shilling, the same figure as the member balance and deduction.
	// Other deductions are owed out of milk payments (default-rate shillings) and use credit/default.
	b.CreditLimit = roundMoney(b.MilkValueUnpaid - b.RejectValue - b.LoanObligations - b.OtherDeductions*ratio)
	b.DefaultBasisLimit = roundMoney(b.MilkValueDefault - rejectDefault - b.LoanObligations - b.OtherDeductions)
	return b, nil
}

// Recalculate applies an event (milk, reject, payroll, loan) to the member's STORED credit limit.
//
// The stored credit_limit is the source of truth and is never overwritten with the formula value.
// Instead the formula (unpaid milk − rejects − loan obligations − deductions) is computed now and
// compared with its value at the previous recalculation (credit_limit_computed); the difference is
// added to credit_limit. So recording milk raises the limit, an approved loan lowers it by the amount
// to repay, payroll lowers it by the milk paid out — and a limit set by hand stays as the base.
// The first recalculation for a member only records the baseline and leaves credit_limit unchanged.
// It must run inside a transaction.
func (s *MemberCreditService) Recalculate(tx *gorm.DB, memberID uint64, change CreditChange) (*MemberCreditBreakdown, error) {
	var member models.Member
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id, route_id, credit_limit, credit_limit_computed, status").First(&member, memberID).Error; err != nil {
		return nil, fmt.Errorf("member %d not found", memberID)
	}

	b, err := s.Compute(tx, member.ID, member.RouteID)
	if err != nil {
		return nil, err
	}
	computed := b.CreditLimit

	oldLimit := roundMoney(member.CreditLimit)
	delta := 0.0
	if member.CreditLimitComputed != nil {
		delta = roundMoney(computed - *member.CreditLimitComputed)
	}
	newLimit := roundMoney(oldLimit + delta)
	b.CreditLimit = newLimit
	if change.SourceType != "credit_baseline" {
		var prev *float64
		if member.CreditLimitComputed != nil {
			v := *member.CreditLimitComputed
			prev = &v
		}
		writeCreditDiagnosticLog(member.ID, change, b, computed, prev, delta, oldLimit, newLimit)
	}

	borrow := s.EvaluateBorrowStatus(tx, member.ID, member.Status, newLimit)

	now := time.Now()
	// Member balance = the balance at the default rate. Loans already come off the credit limit at full
	// value, so add them back before converting and take them off again at full value: a loan lowers
	// the balance and the credit limit by the same Total Deductible.
	memberBalance := b.DefaultBasisLimit
	if r := (MemberRates{CreditRate: b.CreditRate, DefaultRate: b.DefaultRate}).Ratio(); r > 0 {
		memberBalance = roundMoney((newLimit+b.LoanObligations)/r - b.LoanObligations)
	}
	b.DefaultBasisLimit = memberBalance
	updates := map[string]interface{}{
		"member_balance":        memberBalance,
		"credit_limit_computed": computed,
		"can_borrow":            borrow.CanBorrow,
		"cannot_borrow_code":    borrow.Code,
		"cannot_borrow_reason":  borrow.Reason,
	}
	if math.Abs(delta) >= 0.005 {
		updates["credit_limit"] = newLimit
		updates["credit_limit_updated_at"] = now
	}
	if err := tx.Model(&models.Member{}).Where("id = ?", member.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	b.UpdatedAt = &now

	if math.Abs(delta) < 0.005 {
		return b, nil
	}

	action := change.Action
	if action == "" {
		action = models.CreditActionRecalculated
	}
	entry := &models.CreditLimitChangeLog{
		CustomerID:   member.ID,
		CustomerType: creditCustomerTypeMember,
		OldLimit:     oldLimit,
		CreditLimit:  newLimit,
		Amount:       delta,
		Action:       action,
		Reference:    change.Reference,
		SourceType:   change.SourceType,
		SourceID:     change.SourceID,
		Narration:    change.Narration,
		CreatedBy:    change.UserID,
		UpdatedBy:    change.UserID,
		CreatedAt:    &now,
		UpdatedAt:    now.Format("2006-01-02 15:04:05"),
	}
	if err := tx.Create(entry).Error; err != nil {
		return nil, fmt.Errorf("write credit limit log: %w", err)
	}
	return b, nil
}

// EnsureCreditBaselinesAsync records credit_limit_computed for members that have none yet, without
// changing their stored credit_limit, so later events adjust from the current stored value.
func (s *MemberCreditService) EnsureCreditBaselinesAsync() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[MemberCredit] baseline panic: %v", r)
			}
		}()
		var ids []uint64
		if err := db.DB.Model(&models.Member{}).Where("credit_limit_computed IS NULL").Pluck("id", &ids).Error; err != nil {
			log.Printf("[MemberCredit] baseline: %v", err)
			return
		}
		done := 0
		for _, id := range ids {
			if err := db.DB.Transaction(func(tx *gorm.DB) error {
				_, err := s.Recalculate(tx, id, CreditChange{Action: models.CreditActionRecalculated, SourceType: "credit_baseline"})
				return err
			}); err != nil {
				log.Printf("[MemberCredit] baseline member %d: %v", id, err)
				continue
			}
			done++
		}
		if len(ids) > 0 {
			log.Printf("[MemberCredit] recorded credit baselines for %d/%d members (stored limits unchanged)", done, len(ids))
		}
		s.migrateToConfirmedMilkCredit()
		s.runCreditReconcileOnce()
		s.runAdvanceBasisSwitchOnce()
		s.runLoanDeductibleSwitchOnce()
	}()
}

// RecalculateNow runs Recalculate in its own transaction.
func (s *MemberCreditService) RecalculateNow(memberID uint64, change CreditChange) (*MemberCreditBreakdown, error) {
	var result *MemberCreditBreakdown
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		b, err := s.Recalculate(tx, memberID, change)
		result = b
		return err
	})
	return result, err
}

// RecalculateAsync refreshes the given members in the background after the caller's transaction has committed.
// Failures are logged; the next event or read recomputes from source data anyway.
func (s *MemberCreditService) RecalculateAsync(memberIDs []uint64, change CreditChange) {
	ids := uniqueMemberIDs(memberIDs)
	if len(ids) == 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[MemberCredit] recalculation panic: %v", r)
			}
		}()
		for _, id := range ids {
			if _, err := s.RecalculateNow(id, change); err != nil {
				log.Printf("[MemberCredit] recalculate member %d (%s): %v", id, change.Action, err)
			}
		}
	}()
}

func uniqueMemberIDs(ids []uint64) []uint64 {
	seen := make(map[uint64]bool, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

// RecalculateAll refreshes every member in the background and returns how many were queued.
func (s *MemberCreditService) RecalculateAll(userID uint64) (int, error) {
	var ids []uint64
	if err := db.DB.Model(&models.Member{}).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	s.RecalculateAsync(ids, CreditChange{
		Action:     models.CreditActionRecalculated,
		SourceType: "backfill",
		Narration:  "Full credit limit recalculation",
		UserID:     userID,
	})
	return len(ids), nil
}

// Codes stored in member_registrations.cannot_borrow_code when can_borrow is false.
const (
	BorrowBlockMemberStatus       = "MEMBER_STATUS"
	BorrowBlockPendingApplication = "PENDING_APPLICATION"
	BorrowBlockNoCreditLimit      = "NO_CREDIT_LIMIT"
	BorrowBlockLimitInsufficient  = "LIMIT_INSUFFICIENT"
	BorrowBlockNoAdvanceRate      = "NO_ADVANCE_RATE"
	// BorrowBlockNoCashoutDefault: no default member cashout product is set in Loan Products.
	BorrowBlockNoCashoutDefault = "NO_CASHOUT_DEFAULT"
)

// pendingLoanStatuses are applications still being processed. A member's last loan counts as
// processed once it is APPROVED, DISBURSED, REJECTED or CANCELLED (which covers reversed
// disbursements and expired/abandoned OTP requests).
var pendingLoanStatuses = []string{models.LoanAppDraft, models.LoanAppSubmitted, models.LoanAppUnderReview}

// BorrowStatus is the result of EvaluateBorrowStatus. Code and Reason are nil when CanBorrow is true.
type BorrowStatus struct {
	CanBorrow          bool
	Code               *string
	Reason             *string
	MinBorrowAmount    float64
	PendingApplication *PendingLoanApplicationInfo
}

// PendingLoanApplicationInfo describes the member's last loan application while it is unprocessed.
type PendingLoanApplicationInfo struct {
	ID              uint64  `json:"id"`
	Status          string  `json:"status"`
	RequestedAmount float64 `json:"requested_amount"`
}

// minBorrowAmount is the Min Amount of the default member loan product (Loan Products settings).
func minBorrowAmount() float64 {
	if product, err := resolveMobileLoanProduct(0); err == nil {
		return product.MinAmount
	}
	return 0
}

// lastPendingLoanApplication returns the member's most recent application if it is still unprocessed.
// A mobile draft whose OTP has expired is treated as abandoned, not pending.
func lastPendingLoanApplication(tx *gorm.DB, memberID uint64) *models.LoanApplication {
	return lastPendingLoanApplicationExcluding(tx, memberID, 0)
}

// lastPendingLoanApplicationExcluding is lastPendingLoanApplication ignoring application excludeID
// (used when re-checking a request that is itself still awaiting OTP).
func lastPendingLoanApplicationExcluding(tx *gorm.DB, memberID, excludeID uint64) *models.LoanApplication {
	var last models.LoanApplication
	q := tx.Where("member_id = ?", memberID)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Order("id DESC").Limit(1).Find(&last).Error; err != nil || last.ID == 0 {
		return nil
	}
	pending := false
	for _, st := range pendingLoanStatuses {
		if last.Status == st {
			pending = true
			break
		}
	}
	if !pending {
		return nil
	}
	if last.Status == models.LoanAppDraft && last.Source == models.LoanAppSourceMobile {
		var live int64
		tx.Model(&models.LoanApplicationOTP{}).
			Where("loan_application_id = ? AND status = ? AND expires_at > ?", last.ID, models.LoanOTPStatusPending, time.Now()).
			Count(&live)
		if live == 0 {
			return nil
		}
	}
	return &last
}

func pendingApplicationReason(app *models.LoanApplication) string {
	switch app.Status {
	case models.LoanAppDraft:
		return fmt.Sprintf("Loan application #%d is awaiting OTP confirmation. Confirm it, or wait for it to expire, before requesting another cashout.", app.ID)
	default:
		return fmt.Sprintf("Loan application #%d is still being processed (%s). A new cashout can be requested once it is approved, disbursed or rejected.",
			app.ID, strings.ToLower(strings.ReplaceAll(app.Status, "_", " ")))
	}
}

func strPtr(v string) *string { return &v }

// EvaluateBorrowStatus decides whether a member can borrow given their current credit limit.
// Order of checks: member status, last loan still unprocessed, credit limit vs minimum borrowable amount.
// Milk deliveries raise the credit limit, so a member blocked for limit becomes eligible again once
// enough milk is recorded (each milk entry triggers a recalculation).
func (s *MemberCreditService) EvaluateBorrowStatus(tx *gorm.DB, memberID uint64, memberStatus string, creditLimit float64) BorrowStatus {
	return s.evaluateBorrowStatus(tx, memberID, memberStatus, creditLimit, 0)
}

// evaluateBorrowStatus is EvaluateBorrowStatus ignoring application excludeAppID when looking for
// another unprocessed loan.
func (s *MemberCreditService) evaluateBorrowStatus(tx *gorm.DB, memberID uint64, memberStatus string, creditLimit float64, excludeAppID uint64) BorrowStatus {
	minAmount := minBorrowAmount()
	result := BorrowStatus{MinBorrowAmount: minAmount}

	if !memberCanBorrow(memberStatus) {
		result.Code = strPtr(BorrowBlockMemberStatus)
		result.Reason = strPtr(fmt.Sprintf("Member account is %s and cannot borrow.", strings.ToLower(strings.TrimSpace(memberStatus))))
		return result
	}

	// Advances are priced on the advance rate: none set for this member (member / route / all) -> no advances.
	if MemberAdvanceRate(tx, memberID) <= 0 {
		result.Code = strPtr(BorrowBlockNoAdvanceRate)
		result.Reason = strPtr("Advances are not available yet: no advance rate is set for your route. Please contact the dairy office.")
		return result
	}

	if pending := lastPendingLoanApplicationExcluding(tx, memberID, excludeAppID); pending != nil {
		result.PendingApplication = &PendingLoanApplicationInfo{
			ID:              pending.ID,
			Status:          pending.Status,
			RequestedAmount: pending.RequestedAmount,
		}
		result.Code = strPtr(BorrowBlockPendingApplication)
		result.Reason = strPtr(pendingApplicationReason(pending))
		return result
	}

	available := math.Max(0, creditLimit)
	if available <= 0 {
		result.Code = strPtr(BorrowBlockNoCreditLimit)
		result.Reason = strPtr(fmt.Sprintf("No credit limit available. Deliver more milk to reach the minimum borrowable amount of KES %s.", formatKES(minAmount)))
		return result
	}
	if available < minAmount {
		result.Code = strPtr(BorrowBlockLimitInsufficient)
		result.Reason = strPtr(fmt.Sprintf("Credit limit not sufficient: KES %s is below the minimum borrowable amount of KES %s. Delivering more milk raises the limit.",
			formatKES(available), formatKES(minAmount)))
		return result
	}

	result.CanBorrow = true
	return result
}

// RefreshBorrowStatusAsync re-evaluates can_borrow (against the stored credit limit) after a loan
// status change has been committed. It never changes credit_limit.
func (s *MemberCreditService) RefreshBorrowStatusAsync(memberID uint64, sourceType string, sourceID uint64, userID uint64) {
	if memberID == 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[MemberCredit] borrow status refresh panic: %v", r)
			}
		}()
		if err := db.DB.Transaction(func(tx *gorm.DB) error {
			return s.RefreshBorrowStatus(tx, memberID)
		}); err != nil {
			log.Printf("[MemberCredit] refresh borrow status for member %d (%s %d): %v", memberID, sourceType, sourceID, err)
		}
	}()
}

// RefreshBorrowStatus re-evaluates can_borrow against the member's STORED credit limit.
// It never changes credit_limit.
func (s *MemberCreditService) RefreshBorrowStatus(tx *gorm.DB, memberID uint64) error {
	var member models.Member
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id, status, credit_limit").First(&member, memberID).Error; err != nil {
		return fmt.Errorf("member %d not found", memberID)
	}
	borrow := s.EvaluateBorrowStatus(tx, member.ID, member.Status, member.CreditLimit)
	return tx.Model(&models.Member{}).Where("id = ?", member.ID).Updates(map[string]interface{}{
		"can_borrow":           borrow.CanBorrow,
		"cannot_borrow_code":   borrow.Code,
		"cannot_borrow_reason": borrow.Reason,
	}).Error
}

// AdjustCreditLimit moves the member's stored credit limit by delta (negative when a loan is
// taken against it, positive when it is released), logs the change and refreshes can_borrow.
// It does not recompute the limit from milk or loans. Must run inside a transaction.
func (s *MemberCreditService) AdjustCreditLimit(tx *gorm.DB, memberID uint64, delta float64, change CreditChange) (float64, error) {
	var member models.Member
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id, status, credit_limit").First(&member, memberID).Error; err != nil {
		return 0, fmt.Errorf("member %d not found", memberID)
	}
	oldLimit := roundMoney(member.CreditLimit)
	newLimit := roundMoney(oldLimit + delta)
	borrow := s.EvaluateBorrowStatus(tx, member.ID, member.Status, newLimit)

	now := time.Now()
	if err := tx.Model(&models.Member{}).Where("id = ?", member.ID).Updates(map[string]interface{}{
		"credit_limit":            newLimit,
		"credit_limit_updated_at": now,
		"can_borrow":              borrow.CanBorrow,
		"cannot_borrow_code":      borrow.Code,
		"cannot_borrow_reason":    borrow.Reason,
	}).Error; err != nil {
		return 0, err
	}
	if math.Abs(delta) < 0.005 {
		return newLimit, nil
	}
	action := change.Action
	if action == "" {
		action = models.CreditActionRecalculated
	}
	entry := &models.CreditLimitChangeLog{
		CustomerID:   member.ID,
		CustomerType: creditCustomerTypeMember,
		OldLimit:     oldLimit,
		CreditLimit:  newLimit,
		Amount:       roundMoney(delta),
		Action:       action,
		Reference:    change.Reference,
		SourceType:   change.SourceType,
		SourceID:     change.SourceID,
		Narration:    change.Narration,
		CreatedBy:    change.UserID,
		UpdatedBy:    change.UserID,
		CreatedAt:    &now,
		UpdatedAt:    now.Format("2006-01-02 15:04:05"),
	}
	if err := tx.Create(entry).Error; err != nil {
		return 0, fmt.Errorf("write credit limit log: %w", err)
	}
	return newLimit, nil
}

// StoredCreditLimit returns member_registrations.credit_limit without recalculating it.
func (s *MemberCreditService) StoredCreditLimit(memberID uint64) float64 {
	var member models.Member
	if err := db.DB.Select("id, credit_limit").First(&member, memberID).Error; err != nil {
		return 0
	}
	return member.CreditLimit
}

// MemberCreditLimitResponse is returned by GET /api/member-credit-limit.
// It is read straight from member_registrations; nothing is calculated on read.
// credit_limit and can_borrow are maintained by Recalculate when milk, loans or payroll change.
type MemberCreditLimitResponse struct {
	MemberID             uint64     `json:"member_id"`
	MemberNo             string     `json:"member_no"`
	MemberName           string     `json:"member_name"`
	CreditLimit          float64    `json:"credit_limit"`
	CreditLimitUpdatedAt *time.Time `json:"credit_limit_updated_at"`
	AvailableCredit      float64    `json:"available_credit"`
	MinLoanAmount        float64    `json:"min_loan_amount"`
	MaxLoanAmount        float64    `json:"max_loan_amount"`
	// ChargePerLitre is the member loan product's per-litre cashout charge (KES), for the app's terms.
	ChargePerLitre float64 `json:"charge_per_litre"`
	// AdvanceRate is KES per litre of confirmed milk counted towards advances (0 = not set: no advances).
	AdvanceRate float64 `json:"advance_rate"`
	// Two rates: CreditRate (= AdvanceRate) and DefaultRate (payout). CreditLimitDefault is the same
	// balance at the default rate (credit_limit x default / credit).
	CreditRate         float64 `json:"credit_rate"`
	DefaultRate        float64 `json:"default_rate"`
	CreditLimitDefault float64 `json:"credit_limit_default"`
	// MemberBalance is member_registrations.member_balance (= CreditLimitDefault).
	MemberBalance float64 `json:"member_balance"`
	// MaxRequest: the largest advance the credit limit covers after charges; MaxReceive: what the
	// member gets for it (after fees taken from the payout); ProcessingFeeAtMax: the fee on MaxRequest.
	MaxRequest         float64 `json:"max_request"`
	// SupportPhone: who to call when a cashout stays in Processing (loan_cashout_settings).
	SupportPhone string `json:"support_phone"`
	MaxReceive         float64 `json:"max_receive"`
	ProcessingFeeAtMax float64 `json:"processing_fee_at_max"`
	RepayableAtMax     float64 `json:"repayable_at_max"`
	// CanBorrow mirrors member_registrations.can_borrow. EligibleForLoan is kept for older clients.
	CanBorrow       bool `json:"can_borrow"`
	EligibleForLoan bool `json:"eligible_for_loan"`
	// CannotBorrowCode/CannotBorrowReason are null when CanBorrow is true.
	CannotBorrowCode   *string `json:"cannot_borrow_code"`
	CannotBorrowReason *string `json:"cannot_borrow_reason"`
}

// GetCreditLimit returns the member's stored credit limit and borrow status. It does not recalculate.
func (s *MemberCreditService) GetCreditLimit(memberID, userID uint64) (*MemberCreditLimitResponse, error) {
	var member models.Member
	if err := db.DB.Select("id, member_no, first_name, last_name, other_names, credit_limit, credit_limit_updated_at, can_borrow, cannot_borrow_code, cannot_borrow_reason, member_balance").
		First(&member, memberID).Error; err != nil {
		return nil, fmt.Errorf("member not found")
	}

	resp := &MemberCreditLimitResponse{
		MemberID:             member.ID,
		MemberNo:             member.MemberNo,
		MemberName:           memberDisplayName(member),
		CreditLimit:          member.CreditLimit,
		CreditLimitUpdatedAt: member.CreditLimitUpdatedAt,
		AvailableCredit:      math.Max(0, member.CreditLimit),
		CanBorrow:            member.CanBorrow,
		EligibleForLoan:      member.CanBorrow,
		CannotBorrowCode:     member.CannotBorrowCode,
		CannotBorrowReason:   member.CannotBorrowReason,
	}
	resp.SupportPhone = GetCashoutSettings().SupportPhone
	// The member loan is borrowed strictly against the credit limit, so that is the maximum.
	resp.MaxLoanAmount = resp.AvailableCredit
	if product, err := resolveMobileLoanProduct(0); err == nil {
		resp.MinLoanAmount = product.MinAmount
		resp.ChargePerLitre = product.ChargePerLitre
	}
	mr := GetMemberRates(db.DB, member.ID)
	resp.AdvanceRate, resp.CreditRate, resp.DefaultRate = mr.CreditRate, mr.CreditRate, mr.DefaultRate
	resp.MemberBalance = member.MemberBalance
	if ratio := mr.Ratio(); ratio > 0 {
		resp.CreditLimitDefault = roundMoney(member.CreditLimit / ratio)
		if resp.MemberBalance == 0 && member.CreditLimit != 0 {
			resp.MemberBalance = resp.CreditLimitDefault
		}
	}
	if resp.MemberBalance != 0 {
		resp.CreditLimitDefault = resp.MemberBalance
	}
	if product, err := resolveMobileLoanProduct(0); err == nil && resp.AvailableCredit > 0 {
		ms := NewMobileLoanService()
		term := productTermMonths(product)
		resp.MaxRequest = ms.maxRequestFor(product, member.ID, resp.AvailableCredit, term)
		if resp.MaxRequest > 0 {
			fee, upfront := upfrontFees(product, resp.MaxRequest)
			resp.ProcessingFeeAtMax = fee
			resp.MaxReceive = roundMoney(resp.MaxRequest - upfront)
			resp.RepayableAtMax = ms.projectedRepayable(product, member.ID, resp.MaxRequest, term)
		}
		resp.MaxLoanAmount = resp.MaxRequest
	}
	if resp.AdvanceRate <= 0 && resp.CanBorrow {
		// Read-only: the stored flag catches up on the next credit event.
		code, reason := BorrowBlockNoAdvanceRate, "Advances are not available yet: no advance rate is set for your route. Please contact the dairy office."
		resp.CanBorrow, resp.EligibleForLoan = false, false
		resp.CannotBorrowCode, resp.CannotBorrowReason = &code, &reason
	}
	if _, err := resolveMobileLoanProduct(0); err != nil && resp.CanBorrow {
		// No default member cashout product: tell the member before they try (read-only, like above).
		code, reason := BorrowBlockNoCashoutDefault, "Member cashouts are not available right now: the default member cashout product has not been set. Please contact the admin."
		resp.CanBorrow, resp.EligibleForLoan = false, false
		resp.CannotBorrowCode, resp.CannotBorrowReason = &code, &reason
	}
	return resp, nil
}

func (s *MemberCreditService) ListLogs(memberID uint64, page, limit int) ([]models.CreditLimitChangeLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	var logs []models.CreditLimitChangeLog
	var total int64
	q := db.DB.Model(&models.CreditLimitChangeLog{}).
		Where("customer_id = ? AND customer_type = ?", memberID, creditCustomerTypeMember)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id DESC").Limit(limit).Offset((page - 1) * limit).Find(&logs).Error
	return logs, total, err
}

// memberCanBorrow blocks members whose registration has been suspended or ended.
func memberCanBorrow(status string) bool {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "SUSPENDED", "INACTIVE", "EXITED", "DECEASED", "DELETED", "CLOSED", "REJECTED":
		return false
	default:
		return true
	}
}
