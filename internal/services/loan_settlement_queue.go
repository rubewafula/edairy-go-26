package services

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

// Settling DISBURSED loans selected on the desktop is queued: each loan is settled on its own by a
// worker and the officer gets a notification per loan (settled / could not be settled).

const (
	SettlementQueued     = "QUEUED"
	SettlementProcessing = "PROCESSING"
	SettlementDone       = "DONE"
	SettlementFailed     = "FAILED"
)

// LoanSettlementJob is one queued settlement (table loan_settlement_jobs).
type LoanSettlementJob struct {
	models.BaseModel
	LoanApplicationID uint64     `gorm:"column:loan_application_id;index;not null" json:"loan_application_id"`
	LoanContractID    uint64     `gorm:"column:loan_contract_id;index;not null" json:"loan_contract_id"`
	RequestedBy       uint64     `gorm:"column:requested_by;index" json:"requested_by"`
	Status            string     `gorm:"column:status;size:16;index;not null;default:QUEUED" json:"status"`
	SettledAmount     float64    `gorm:"column:settled_amount;type:decimal(15,2)" json:"settled_amount"`
	Error             string     `gorm:"column:error;type:text" json:"error"`
	ProcessedAt       *time.Time `gorm:"column:processed_at" json:"processed_at"`
}

func (LoanSettlementJob) TableName() string { return "loan_settlement_jobs" }

var (
	settlementSchemaOnce sync.Once
	settlementWorkerOnce sync.Once
	settlementKick       = make(chan struct{}, 1)
)

func ensureSettlementSchema() {
	settlementSchemaOnce.Do(func() {
		if err := db.DB.AutoMigrate(&LoanSettlementJob{}); err != nil {
			log.Printf("[LoanSettlement] migrate: %v", err)
		}
	})
}

// SettleQueueResult says what was queued and what was skipped (and why).
type SettleQueueResult struct {
	Queued  int               `json:"queued"`
	Skipped map[uint64]string `json:"skipped"`
}

// QueueSettlements queues the given applications; only DISBURSED ones with an open contract are taken.
func (s *LoanModuleService) QueueSettlements(appIDs []uint64, userID uint64) (*SettleQueueResult, error) {
	ensureSettlementSchema()
	res := &SettleQueueResult{Skipped: map[uint64]string{}}
	var apps []models.LoanApplication
	if err := db.DB.Where("id IN ?", appIDs).Find(&apps).Error; err != nil {
		return nil, err
	}
	found := map[uint64]bool{}
	for _, a := range apps {
		found[a.ID] = true
		if a.Status != models.LoanAppDisbursed || a.LoanContractID == nil {
			res.Skipped[a.ID] = fmt.Sprintf("only disbursed loans can be settled (status %s)", a.Status)
			continue
		}
		var c models.LoanContract
		if err := db.DB.Select("id, status").First(&c, *a.LoanContractID).Error; err != nil {
			res.Skipped[a.ID] = "loan contract not found"
			continue
		}
		if c.Status == models.LoanContractSettled || c.Status == models.LoanContractClosed || c.Status == models.LoanContractWrittenOff {
			res.Skipped[a.ID] = "loan is already " + strings.ToLower(c.Status)
			continue
		}
		if contractHasPayoutInFlight(c.ID) {
			res.Skipped[a.ID] = "payout not completed yet (only disbursed loans are settled)"
			continue
		}
		var open int64
		db.DB.Model(&LoanSettlementJob{}).Where("loan_contract_id = ? AND status IN ?", c.ID,
			[]string{SettlementQueued, SettlementProcessing}).Count(&open)
		if open > 0 {
			res.Skipped[a.ID] = "already queued for settlement"
			continue
		}
		job := LoanSettlementJob{LoanApplicationID: a.ID, LoanContractID: c.ID, RequestedBy: userID, Status: SettlementQueued}
		job.CreatedBy = userID
		if err := db.DB.Create(&job).Error; err != nil {
			return nil, err
		}
		res.Queued++
	}
	for _, id := range appIDs {
		if !found[id] {
			res.Skipped[id] = "application not found"
		}
	}
	if res.Queued > 0 {
		select {
		case settlementKick <- struct{}{}:
		default:
		}
	}
	return res, nil
}

// contractHasPayoutInFlight reports a payout that has not completed (queued, sending, M-Pesa error).
func contractHasPayoutInFlight(contractID uint64) bool {
	var n int64
	db.DB.Model(&models.LoanDisbursement{}).Where("loan_contract_id = ? AND status IN ?", contractID,
		[]string{models.DisbursementStatusPending, models.DisbursementStatusProcessing,
			models.DisbursementStatusMpesaError, models.DisbursementStatusMpesaFailed}).Count(&n)
	return n > 0
}

// settlementBlockedReason re-checks a queued settlement: the application must still be DISBURSED,
// the loan open and its payout complete. "" = go ahead.
func settlementBlockedReason(job *LoanSettlementJob) string {
	var a models.LoanApplication
	if err := db.DB.Select("id, status").First(&a, job.LoanApplicationID).Error; err != nil {
		return "loan application not found"
	}
	if a.Status != models.LoanAppDisbursed {
		return fmt.Sprintf("only disbursed loans are settled (status %s)", a.Status)
	}
	var c models.LoanContract
	if err := db.DB.Select("id, status").First(&c, job.LoanContractID).Error; err != nil {
		return "loan contract not found"
	}
	if c.Status == models.LoanContractSettled || c.Status == models.LoanContractClosed || c.Status == models.LoanContractWrittenOff {
		return "loan is already " + strings.ToLower(c.Status)
	}
	if contractHasPayoutInFlight(c.ID) {
		return "payout not completed yet"
	}
	return ""
}

// StartSettlementWorker settles queued loans one at a time (every 10 s, or at once when queued).
func StartSettlementWorker() {
	settlementWorkerOnce.Do(func() {
		ensureSettlementSchema()
		go func() {
			svc := NewLoanModuleService()
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				svc.processSettlementQueue()
				select {
				case <-ticker.C:
				case <-settlementKick:
				}
			}
		}()
	})
}

func (s *LoanModuleService) processSettlementQueue() {
	var jobs []LoanSettlementJob
	if err := db.DB.Where("status = ?", SettlementQueued).Order("id").Limit(50).Find(&jobs).Error; err != nil {
		return
	}
	for _, job := range jobs {
		claimed := db.DB.Model(&LoanSettlementJob{}).Where("id = ? AND status = ?", job.ID, SettlementQueued).
			Update("status", SettlementProcessing)
		if claimed.RowsAffected == 0 {
			continue
		}
		s.settleOne(&job)
	}
}

func (s *LoanModuleService) settleOne(job *LoanSettlementJob) {
	now := time.Now()
	var contract models.LoanContract
	db.DB.Select("id, contract_no, member_id").First(&contract, job.LoanContractID)
	member := memberFirstName(contract.MemberID)
	// Only loans that are DISBURSED when the job runs are settled; anything else is left untouched.
	if reason := settlementBlockedReason(job); reason != "" {
		db.DB.Model(&LoanSettlementJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
			"status": SettlementFailed, "error": reason, "processed_at": now})
		notifySettlement(job, "Loan not settled",
			fmt.Sprintf("Loan %s (%s) was not settled: %s", contract.ContractNo, member, reason), "WARNING")
		appendAppLog("disbursement.log", fmt.Sprintf("settlement job %d contract %s skipped: %s", job.ID, contract.ContractNo, reason))
		return
	}
	quote, err := s.SettleContract(job.LoanContractID, dtos.SettleLoanRequest{
		PostToGL: true,
		Channel:  models.LoanChannelMilkDeduction,
		Notes:    "Loan settled (selected on desktop)",
	}, job.RequestedBy)
	if err != nil {
		db.DB.Model(&LoanSettlementJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
			"status": SettlementFailed, "error": err.Error(), "processed_at": now})
		notifySettlement(job, "Loan not settled",
			fmt.Sprintf("Loan %s (%s) could not be settled: %s", contract.ContractNo, member, err.Error()), "ERROR")
		appendAppLog("disbursement.log", fmt.Sprintf("settlement job %d contract %s FAILED: %v", job.ID, contract.ContractNo, err))
		return
	}
	amount := 0.0
	if quote != nil {
		amount = quote.TotalSettlement
	}
	// Make sure the loan really is closed: pay any remainder the schedule still shows, then mark it
	// SETTLED with today's date if balances still disagree (the officer chose to settle it).
	var st struct {
		Status    string
		SettledAt *time.Time
	}
	db.DB.Raw(`SELECT status, settled_at FROM loan_contracts WHERE id = ?`, contract.ID).Scan(&st)
	if st.Status != models.LoanContractSettled {
		if rem := contractRemainingRepayable(db.DB, contract.ID); rem > 0.01 {
			if _, err := s.RecordRepayment(contract.ID, dtos.RecordLoanRepaymentRequest{
				Amount: rem, Channel: models.LoanChannelMilkDeduction, PaymentDate: now.Format("2006-01-02"),
				Reference: fmt.Sprintf("SETTLE-REM-%s-%d", contract.ContractNo, now.Unix()), Notes: "Settlement remainder",
			}, job.RequestedBy); err == nil {
				amount = roundMoney(amount + rem)
			}
		}
		db.DB.Raw(`SELECT status, settled_at FROM loan_contracts WHERE id = ?`, contract.ID).Scan(&st)
	}
	if st.Status != models.LoanContractSettled || st.SettledAt == nil {
		db.DB.Exec(`UPDATE loan_contracts SET status = ?, settled_at = COALESCE(settled_at, ?), updated_by = ?, updated_at = NOW()
			WHERE id = ?`, models.LoanContractSettled, now.Format("2006-01-02"), job.RequestedBy, contract.ID)
	}
	// The milk deduction for this loan stops: mark it settled.
	db.DB.Exec(`UPDATE recurrent_deductions SET settled = 1, updated_at = NOW()
		WHERE reference = ? AND customer_id = ? AND customer_type = 'member' AND settled = 0`,
		"LOAN-"+contract.ContractNo, contract.MemberID)
	NewMemberCreditService().RecalculateAsync([]uint64{contract.MemberID}, CreditChange{
		Action: models.CreditActionLoanSettled, SourceType: "loan_contract", SourceID: contract.ID,
		Reference: contract.ContractNo, Narration: "Loan settled", UserID: job.RequestedBy,
	})
	db.DB.Model(&LoanSettlementJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
		"status": SettlementDone, "settled_amount": amount, "processed_at": now, "error": ""})
	notifySettlement(job, "Loan settled",
		fmt.Sprintf("Loan %s (%s) has been settled on %s: KES %s.", contract.ContractNo, member, now.Format("02 Jan 2006"), formatKES(amount)), "SUCCESS")
	appendAppLog("disbursement.log", fmt.Sprintf("settlement job %d contract %s settled KES %s", job.ID, contract.ContractNo, formatKES(amount)))
}

func notifySettlement(job *LoanSettlementJob, title, message, kind string) {
	if job.RequestedBy == 0 {
		return
	}
	ref := job.LoanApplicationID
	refType := "loan_application"
	if _, err := NewUINotificationService().CreateNotification(job.RequestedBy, dtos.CreateUINotificationRequest{
		Title: title, Message: message, NotificationType: kind, ReferenceID: &ref, ReferenceType: &refType,
	}); err != nil {
		log.Printf("[LoanSettlement] notify user %d: %v", job.RequestedBy, err)
	}
}

// EnsureSettlementPermission lets roles that can settle contracts settle selected applications.
func EnsureSettlementPermission() {
	ensureSettlementSchema()
	for _, extra := range []string{"loan-module.applications.settle-month.create", "loan-module.applications.settlement-summary.view"} {
		ensureSettlementPermissionNamed(extra)
	}
	ensureSettlementPermissionNamed("loan-module.applications.settle.create")
	ensureDBUAccountantRole()
}

const (
	// PermClearDairyLoans allows "Settle all" (clearing a month's dairy loans).
	PermClearDairyLoans = "loan-module.dairy-loans.clear"
	// RoleDBUEdairyAccountant is the only role given PermClearDairyLoans.
	RoleDBUEdairyAccountant      = "DBU Edairy Accountant"
	ClearDairyLoansDeniedMessage = "Only the DBU Edairy Accountant is allowed to clear dairy loans."
)

// CanClearDairyLoans reports whether the caller's permissions include clearing dairy loans.
func CanClearDairyLoans(perms []string) bool {
	for _, p := range perms {
		if p == PermClearDairyLoans {
			return true
		}
	}
	return false
}

// ensureDBUAccountantRole creates the DBU Edairy Accountant role and the clear-dairy-loans permission,
// and gives the role what it needs to open the loan applications and settle a month.
func ensureDBUAccountantRole() {
	db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
		SELECT ?, 'web', NOW(), NOW() FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = ?)`,
		PermClearDairyLoans, PermClearDairyLoans)
	db.DB.Exec(`INSERT INTO roles (name, guard_name, created_at, updated_at)
		SELECT ?, 'web', NOW(), NOW() FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = ?)`,
		RoleDBUEdairyAccountant, RoleDBUEdairyAccountant)
	db.DB.Exec(`UPDATE roles SET deleted_at = NULL WHERE name = ?`, RoleDBUEdairyAccountant)
	for _, perm := range []string{
		PermClearDairyLoans,
		"loan-module.applications.settle-month.create",
		"loan-module.applications.settlement-summary.view",
		"loan-module.applications.view",
		"loan-module.products.view",
	} {
		db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
			SELECT r.id, p.id FROM roles r JOIN permissions p ON p.name = ? AND p.deleted_at IS NULL
			WHERE r.name = ? AND r.deleted_at IS NULL
			AND NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = r.id AND x.permission_id = p.id)`,
			perm, RoleDBUEdairyAccountant)
	}
}

func ensureSettlementPermissionNamed(name string) {
	db.DB.Exec(`INSERT INTO permissions (name, guard_name, created_at, updated_at)
		SELECT ?, 'web', NOW(), NOW() FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE name = ?)`, name, name)
	db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT DISTINCT rp.role_id, np.id FROM role_permissions rp
		JOIN permissions src ON src.id = rp.permission_id AND src.deleted_at IS NULL
			AND src.name IN ('loan-module.contracts.settle.create', 'loan-module.applications.approve.create')
		JOIN permissions np ON np.name = ? AND np.deleted_at IS NULL
		WHERE NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = rp.role_id AND x.permission_id = np.id)`, name)
	db.DB.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id FROM roles r JOIN permissions p ON p.name = ? AND p.deleted_at IS NULL
		WHERE r.name = 'Super Admin' AND NOT EXISTS (SELECT 1 FROM role_permissions x WHERE x.role_id = r.id AND x.permission_id = p.id)`, name)
}

// EnsureContractSettledDates fills settled_at for loans settled before the column existed
// (last payment date, else when the contract was last updated).
func EnsureContractSettledDates() {
	// Old rows saved an empty idempotency key; turn them into NULL so the unique index stops clashing.
	for _, t := range []string{"loan_repayments", "loan_disbursements"} {
		db.DB.Exec("UPDATE " + t + " SET idempotency_key = NULL WHERE idempotency_key = ''")
	}
	if !db.DB.Migrator().HasColumn(&models.LoanContract{}, "settled_at") {
		if err := db.DB.Migrator().AddColumn(&models.LoanContract{}, "SettledAt"); err != nil {
			log.Printf("[LoanSettlement] add settled_at: %v", err)
			return
		}
	}
	res := db.DB.Exec(`UPDATE loan_contracts SET settled_at = DATE(COALESCE(last_payment_date, updated_at))
		WHERE status = ? AND settled_at IS NULL AND deleted_at IS NULL`, models.LoanContractSettled)
	if res.Error == nil && res.RowsAffected > 0 {
		log.Printf("[LoanSettlement] filled settled_at for %d settled loans", res.RowsAffected)
	}
}

// MonthSettlementSummary describes disbursed loans of one month (by disbursement date) and all
// unsettled loans, for the "Settle for the month" panel.
type MonthSettlementSummary struct {
	Month            string  `json:"month"` // YYYY-MM
	UnsettledLoans   int64   `json:"unsettled_loans"`
	Borrowed         float64 `json:"borrowed"`          // principal disbursed (unsettled, this month)
	TotalRepayable   float64 `json:"total_repayable"`   // principal + interest + charges (unsettled, this month)
	Outstanding      float64 `json:"outstanding"`       // still to be repaid (unsettled, this month)
	SettledLoans     int64   `json:"settled_loans"`     // already settled, disbursed this month
	SettledAmount    float64 `json:"settled_amount"`    // their repayable
	QueuedLoans      int64   `json:"queued_loans"`      // waiting in the settlement queue
	AllUnsettled     int64   `json:"all_unsettled_loans"`
	AllOutstanding   float64 `json:"all_outstanding"`
	AllBorrowed      float64 `json:"all_borrowed"`
	// Product (product_id query): its DBU allocation quota and what is accumulated on it so far
	// that has not been paid (all unsettled disbursed loans of the product, any month).
	Product *ProductDBUSummary `json:"product,omitempty"`
}

// ProductDBUSummary compares a product's DBU allocation quota with its unpaid loans.
type ProductDBUSummary struct {
	ProductID       uint64  `json:"product_id"`
	DBUMinAmount    float64 `json:"dbu_min_amount"`
	DBUMaxAmount    float64 `json:"dbu_max_amount"`
	UnpaidLoans     int64   `json:"unpaid_loans"`
	UnpaidBorrowed  float64 `json:"unpaid_borrowed"`    // principal of unsettled loans (accumulated so far)
	UnpaidOutstand  float64 `json:"unpaid_outstanding"` // still to be repaid on them
	RemainingQuota  float64 `json:"remaining_quota"`    // dbu_max_amount - unpaid_borrowed (never below 0)
}

const unsettledDisbursedSQL = `FROM loan_applications a
	JOIN loan_contracts c ON c.id = a.loan_contract_id AND c.deleted_at IS NULL
	WHERE a.deleted_at IS NULL AND a.status = 'DISBURSED'
	AND c.status NOT IN ('SETTLED', 'CLOSED', 'WRITTEN_OFF')
	AND NOT EXISTS (SELECT 1 FROM loan_disbursements ld WHERE ld.loan_contract_id = c.id AND ld.deleted_at IS NULL
		AND ld.status IN ('PENDING', 'PROCESSING', 'MPESA_ERROR', 'MPESA_FAILED'))`

const contractRepayableSQL = `(SELECT COALESCE(SUM(i.principal_due + i.interest_due + i.fee_due + i.insurance_due + i.penalty_due), 0)
	FROM loan_schedule_installments i WHERE i.loan_contract_id = c.id AND i.deleted_at IS NULL)`

const contractOutstandingSQL = `(SELECT COALESCE(SUM(GREATEST((i.principal_due + i.interest_due + i.fee_due + i.insurance_due + i.penalty_due)
	- (i.principal_paid + i.interest_paid + i.fee_paid + i.insurance_paid + i.penalty_paid), 0)), 0)
	FROM loan_schedule_installments i WHERE i.loan_contract_id = c.id AND i.deleted_at IS NULL)`

// monthCondition matches loans disbursed in month (disbursement date, else application date).
const monthCondition = ` AND DATE_FORMAT(COALESCE(c.disbursement_date, a.date_applied), '%Y-%m') = ?`

// SettlementSummary returns the month's figures (month "" = current month).
func (s *LoanModuleService) SettlementSummary(month string, productID uint64) (*MonthSettlementSummary, error) {
	ensureSettlementSchema()
	month = normalizeMonth(month)
	out := &MonthSettlementSummary{Month: month}
	var r struct {
		N           int64
		Borrowed    float64
		Repayable   float64
		Outstanding float64
	}
	db.DB.Raw(`SELECT COUNT(*) AS n, COALESCE(SUM(c.principal_amount), 0) AS borrowed,
		COALESCE(SUM(`+contractRepayableSQL+`), 0) AS repayable, COALESCE(SUM(`+contractOutstandingSQL+`), 0) AS outstanding
		`+unsettledDisbursedSQL+monthCondition, month).Scan(&r)
	out.UnsettledLoans, out.Borrowed, out.TotalRepayable, out.Outstanding = r.N, roundMoney(r.Borrowed), roundMoney(r.Repayable), roundMoney(r.Outstanding)

	var st struct {
		N      int64
		Amount float64
	}
	db.DB.Raw(`SELECT COUNT(*) AS n, COALESCE(SUM(`+contractRepayableSQL+`), 0) AS amount
		FROM loan_applications a JOIN loan_contracts c ON c.id = a.loan_contract_id AND c.deleted_at IS NULL
		WHERE a.deleted_at IS NULL AND c.status = 'SETTLED'`+monthCondition, month).Scan(&st)
	out.SettledLoans, out.SettledAmount = st.N, roundMoney(st.Amount)

	db.DB.Model(&LoanSettlementJob{}).Where("status IN ?", []string{SettlementQueued, SettlementProcessing}).Count(&out.QueuedLoans)

	var all struct {
		N           int64
		Borrowed    float64
		Outstanding float64
	}
	db.DB.Raw(`SELECT COUNT(*) AS n, COALESCE(SUM(c.principal_amount), 0) AS borrowed,
		COALESCE(SUM(`+contractOutstandingSQL+`), 0) AS outstanding `+unsettledDisbursedSQL).Scan(&all)
	out.AllUnsettled, out.AllBorrowed, out.AllOutstanding = all.N, roundMoney(all.Borrowed), roundMoney(all.Outstanding)

	if productID > 0 {
		var p models.LoanProduct
		if err := db.DB.Select("id, dbu_min_amount, dbu_max_amount").First(&p, productID).Error; err == nil {
			var pr struct {
				N           int64
				Borrowed    float64
				Outstanding float64
			}
			db.DB.Raw(`SELECT COUNT(*) AS n, COALESCE(SUM(c.principal_amount), 0) AS borrowed,
				COALESCE(SUM(`+contractOutstandingSQL+`), 0) AS outstanding `+unsettledDisbursedSQL+
				` AND a.loan_product_id = ?`, productID).Scan(&pr)
			out.Product = &ProductDBUSummary{
				ProductID:      p.ID,
				DBUMinAmount:   p.DBUMinAmount,
				DBUMaxAmount:   p.DBUMaxAmount,
				UnpaidLoans:    pr.N,
				UnpaidBorrowed: roundMoney(pr.Borrowed),
				UnpaidOutstand: roundMoney(pr.Outstanding),
				RemainingQuota: roundMoney(math.Max(0, p.DBUMaxAmount-pr.Borrowed)),
			}
		}
	}
	return out, nil
}

// QueueMonthSettlement queues every unsettled disbursed loan of the month; each is settled separately.
func (s *LoanModuleService) QueueMonthSettlement(month string, userID uint64) (*SettleQueueResult, error) {
	month = normalizeMonth(month)
	var ids []uint64
	if err := db.DB.Raw(`SELECT a.id `+unsettledDisbursedSQL+monthCondition+` ORDER BY a.id`, month).Scan(&ids).Error; err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return &SettleQueueResult{Skipped: map[uint64]string{}}, nil
	}
	appendAppLog("disbursement.log", fmt.Sprintf("settle month %s: %d loan(s) requested by user %d", month, len(ids), userID))
	return s.QueueSettlements(ids, userID)
}

func normalizeMonth(month string) string {
	month = strings.TrimSpace(month)
	if t, err := time.Parse("2006-01", month); err == nil {
		return t.Format("2006-01")
	}
	return time.Now().Format("2006-01")
}
