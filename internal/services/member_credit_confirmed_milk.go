package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const confirmedMilkCreditMigration = "member_credit_confirmed_milk_v1"

// migrateToConfirmedMilkCredit runs once. Credit now comes from milk on CONFIRMED journals only, so:
//   - every member's baseline (credit_limit_computed) moves to the new formula without changing the
//     stored credit limit, and
//   - milk on journals that are already confirmed but was never credited (no credit-limit log for it,
//     e.g. recorded before credit tracking started) is credited now, with a MILK_CONFIRMED log line.
func (s *MemberCreditService) migrateToConfirmedMilkCredit() {
	if err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS app_data_migrations (
		name VARCHAR(128) NOT NULL PRIMARY KEY, applied_at DATETIME NOT NULL)`).Error; err != nil {
		log.Printf("[MemberCredit] confirmed-milk migration: %v", err)
		return
	}
	var done int64
	db.DB.Raw(`SELECT COUNT(*) FROM app_data_migrations WHERE name = ?`, confirmedMilkCreditMigration).Scan(&done)
	if done > 0 {
		return
	}

	var ids []uint64
	if err := db.DB.Model(&models.Member{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("[MemberCredit] confirmed-milk migration: %v", err)
		return
	}
	rebased, credited := 0, 0
	var creditedTotal float64
	for _, id := range ids {
		amount, err := s.rebaseMemberOnConfirmedMilk(id)
		if err != nil {
			log.Printf("[MemberCredit] confirmed-milk migration member %d: %v", id, err)
			continue
		}
		rebased++
		if amount > 0.005 {
			credited++
			creditedTotal += amount
		}
	}
	db.DB.Exec(`INSERT IGNORE INTO app_data_migrations (name, applied_at) VALUES (?, NOW())`, confirmedMilkCreditMigration)
	log.Printf("[MemberCredit] credit now follows confirmed milk: %d members re-based, %d credited KES %.2f for confirmed milk not credited before",
		rebased, credited, creditedTotal)
}

// uncreditedConfirmedJournals lists the member's confirmed journals whose milk has no credit-limit log yet.
func uncreditedConfirmedJournals(tx *gorm.DB, memberID uint64) ([]uint64, []string, error) {
	var rows []struct {
		ID      uint64
		Journal string
	}
	err := tx.Raw(`SELECT DISTINCT j.id, j.journal FROM milk_journals j
		JOIN milk_journal_entries e ON e.milk_journal_id = j.id AND e.deleted_at IS NULL AND e.member_id = ?
		WHERE j.deleted_at IS NULL AND j.confirmed = 1
		AND NOT EXISTS (SELECT 1 FROM credit_limit_change_log l
			WHERE l.customer_id = ? AND l.customer_type = ? AND l.amount > 0 AND l.deleted_at IS NULL
			AND ((l.source_type = 'milk_journal' AND l.source_id = j.id)
				OR (l.source_type = 'milk_journal_entry' AND l.source_id IN
					(SELECT e2.id FROM milk_journal_entries e2 WHERE e2.milk_journal_id = j.id))
				OR (l.source_type = 'milk_journal_import' AND l.created_at >= j.created_at)))`,
		memberID, memberID, creditCustomerTypeMember).Scan(&rows).Error
	ids := make([]uint64, 0, len(rows))
	names := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
		names = append(names, firstNonEmptyString(r.Journal, fmt.Sprintf("#%d", r.ID)))
	}
	return ids, names, err
}

func (s *MemberCreditService) rebaseMemberOnConfirmedMilk(memberID uint64) (float64, error) {
	var credited float64
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var member models.Member
		if err := tx.Select("id, route_id").First(&member, memberID).Error; err != nil {
			return err
		}
		journalIDs, names, err := uncreditedConfirmedJournals(tx, memberID)
		if err != nil {
			return err
		}
		// Baseline = new formula without the never-credited journals, so the next step credits exactly them.
		without := &MemberCreditService{excludeJournalIDs: journalIDs}
		b, err := without.Compute(tx, member.ID, member.RouteID)
		if err != nil {
			return err
		}
		if err := tx.Model(&models.Member{}).Where("id = ?", member.ID).
			Update("credit_limit_computed", b.CreditLimit).Error; err != nil {
			return err
		}
		if len(journalIDs) == 0 {
			return nil
		}
		ref := "MILK-JOURNAL"
		var sourceID uint64
		if len(journalIDs) == 1 {
			sourceID = journalIDs[0]
			ref = fmt.Sprintf("MILK-JOURNAL-%d", journalIDs[0])
		}
		var before float64
		tx.Raw(`SELECT credit_limit FROM member_registrations WHERE id = ?`, member.ID).Scan(&before)
		after, err := s.Recalculate(tx, member.ID, CreditChange{
			Action:     models.CreditActionMilkConfirmed,
			SourceType: "milk_journal",
			SourceID:   sourceID,
			Reference:  ref,
			Narration:  "Confirmed milk journal(s) " + strings.Join(names, ", ") + " credited",
		})
		if err != nil {
			return err
		}
		credited = roundMoney(after.CreditLimit - before)
		return nil
	})
	return credited, err
}

const creditReconcileMigration = "member_credit_reconcile_confirmed_milk_v1"

// ReconcileLine is one journal's correction for one member.
type ReconcileLine struct {
	MemberID  uint64  `json:"member_id"`
	JournalID uint64  `json:"journal_id"`
	Journal   string  `json:"journal"`
	Confirmed bool    `json:"confirmed"`
	Kilos     float64 `json:"kilos"`
	Rate      float64 `json:"rate"`
	Expected  float64 `json:"expected"`
	Credited  float64 `json:"credited_before"`
	Adjusted  float64 `json:"adjusted"`
}

// ReconcileResult summarises a reconciliation run.
type ReconcileResult struct {
	Members       int             `json:"members_checked"`
	MembersFixed  int             `json:"members_adjusted"`
	TotalAdjusted float64         `json:"total_adjusted"`
	Lines         []ReconcileLine `json:"lines"`
}

var creditedMilkActions = []string{
	models.CreditActionMilkConfirmed, models.CreditActionMilkUnconfirmed, models.CreditActionMilkEntryChanged,
	models.CreditActionSeeded, "MILK_RECONCILED",
}

// ReconcileConfirmedMilkCredit makes member_registrations.credit_limit agree with milk on journals:
// for every unpaid journal with the member's milk, the credit given so far (sum of the member's
// credit-limit log lines for that journal or its entries) must equal the milk's value if the journal
// is confirmed, or 0 if it is not. Any difference is added to credit_limit with a MILK_RECONCILED log
// line per journal. memberID 0 = all members.
func (s *MemberCreditService) ReconcileConfirmedMilkCredit(memberID uint64, userID uint64) (*ReconcileResult, error) {
	var ids []uint64
	q := db.DB.Model(&models.Member{})
	if memberID > 0 {
		q = q.Where("id = ?", memberID)
	}
	if err := q.Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	res := &ReconcileResult{Lines: []ReconcileLine{}}
	for _, id := range ids {
		lines, err := s.reconcileMember(id, userID)
		if err != nil {
			log.Printf("[MemberCredit] reconcile member %d: %v", id, err)
			continue
		}
		res.Members++
		if len(lines) > 0 {
			res.MembersFixed++
			for _, l := range lines {
				res.TotalAdjusted += l.Adjusted
			}
			res.Lines = append(res.Lines, lines...)
		}
	}
	res.TotalAdjusted = roundMoney(res.TotalAdjusted)
	appendAppLog("credit-changes.log", fmt.Sprintf("RECONCILE confirmed milk: %d members checked, %d adjusted, total %+.2f",
		res.Members, res.MembersFixed, res.TotalAdjusted))
	return res, nil
}

func (s *MemberCreditService) reconcileMember(memberID, userID uint64) ([]ReconcileLine, error) {
	var lines []ReconcileLine
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var member models.Member
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, route_id, credit_limit, status").First(&member, memberID).Error; err != nil {
			return err
		}
		type jrow struct {
			ID        uint64
			Journal   string
			Confirmed bool
			Kilos     float64
			PdrID     uint64
			RouteID   uint64
		}
		var journals []jrow
		if err := tx.Raw(fmt.Sprintf(`SELECT j.id, j.journal, j.confirmed, SUM(e.quantity) AS kilos, `+payDateRangeForDate+` AS pdr_id, COALESCE(j.route_id, 0) AS route_id
			FROM milk_journal_entries e JOIN milk_journals j ON j.id = e.milk_journal_id
			WHERE e.member_id = ? AND e.deleted_at IS NULL AND j.deleted_at IS NULL AND `+notInApprovedPayroll+`
			GROUP BY j.id, j.journal, j.confirmed, j.journal_date, j.route_id`, "j.journal_date"), memberID).Scan(&journals).Error; err != nil {
			return err
		}
		if len(journals) == 0 {
			return nil
		}
		pdrSet := map[uint64]bool{}
		for _, j := range journals {
			if j.PdrID > 0 {
				pdrSet[j.PdrID] = true
			}
		}
		pdrIDs := make([]uint64, 0, len(pdrSet))
		for id := range pdrSet {
			pdrIDs = append(pdrIDs, id)
		}
		rates, err := loadAdvanceRateResolver(tx)
		if err != nil {
			return err
		}
		_ = pdrIDs
		// Credit already given per journal.
		credited := map[uint64]float64{}
		for _, j := range journals {
			var sum float64
			tx.Raw(`SELECT COALESCE(SUM(l.amount), 0) FROM credit_limit_change_log l
				WHERE l.customer_id = ? AND l.customer_type = ? AND l.deleted_at IS NULL AND l.action IN ?
				AND ((l.source_type = 'milk_journal' AND l.source_id = ?)
					OR (l.source_type = 'milk_journal_entry' AND l.source_id IN
						(SELECT e2.id FROM milk_journal_entries e2 WHERE e2.milk_journal_id = ?)))`,
				memberID, creditCustomerTypeMember, creditedMilkActions, j.ID, j.ID).Scan(&sum)
			credited[j.ID] = sum
		}
		// Credits from the one-off migration that named several journals in one log line.
		var multi []struct {
			Narration string
			Amount    float64
		}
		tx.Raw(`SELECT narration, amount FROM credit_limit_change_log WHERE customer_id = ? AND customer_type = ?
			AND deleted_at IS NULL AND action = ? AND source_type = 'milk_journal' AND source_id = 0
			AND narration LIKE 'Confirmed milk journal(s) % credited'`,
			memberID, creditCustomerTypeMember, models.CreditActionMilkConfirmed).Scan(&multi)
		byName := map[string]uint64{}
		value := map[uint64]float64{}
		for _, j := range journals {
			byName[strings.TrimSpace(j.Journal)] = j.ID
			value[j.ID] = roundMoney(j.Kilos * rates.RateForJournal(j.PdrID, memberID, member.RouteID, j.RouteID))
		}
		for _, m := range multi {
			list := strings.TrimSuffix(strings.TrimPrefix(m.Narration, "Confirmed milk journal(s) "), " credited")
			for _, name := range strings.Split(list, ", ") {
				if id, ok := byName[strings.TrimSpace(name)]; ok {
					credited[id] += value[id] // that line credited each named journal's value
				}
			}
		}

		// Credit from rate changes (MILK_RATE_CHANGED) is not tied to one journal: share it over the
		// confirmed journals by kg, so a later reconcile neither re-credits nor re-debits it.
		var ratePool float64
		tx.Raw(`SELECT COALESCE(SUM(amount), 0) FROM credit_limit_change_log WHERE customer_id = ? AND customer_type = ?
			AND deleted_at IS NULL AND action = 'MILK_RATE_CHANGED'`, memberID, creditCustomerTypeMember).Scan(&ratePool)
		var confirmedKg float64
		for _, j := range journals {
			if j.Confirmed {
				confirmedKg += j.Kilos
			}
		}
		if ratePool > 0.005 || ratePool < -0.005 {
			for _, j := range journals {
				if j.Confirmed && confirmedKg > 0 {
					credited[j.ID] += ratePool * j.Kilos / confirmedKg
				}
			}
		}

		total := 0.0
		for _, j := range journals {
			rate := rates.RateForJournal(j.PdrID, memberID, member.RouteID, j.RouteID)
			expected := 0.0
			if j.Confirmed {
				expected = value[j.ID]
			}
			adj := roundMoney(expected - credited[j.ID])
			if adj > -0.01 && adj < 0.01 {
				continue
			}
			total += adj
			lines = append(lines, ReconcileLine{MemberID: memberID, JournalID: j.ID, Journal: j.Journal, Confirmed: j.Confirmed,
				Kilos: j.Kilos, Rate: rate, Expected: expected, Credited: roundMoney(credited[j.ID]), Adjusted: adj})
		}
		if len(lines) == 0 {
			return nil
		}

		now := time.Now()
		limit := roundMoney(member.CreditLimit)
		for _, l := range lines {
			newLimit := roundMoney(limit + l.Adjusted)
			state := "not confirmed"
			if l.Confirmed {
				state = "confirmed"
			}
			if err := tx.Create(&models.CreditLimitChangeLog{
				CustomerID: memberID, CustomerType: creditCustomerTypeMember, OldLimit: limit, CreditLimit: newLimit,
				Amount: l.Adjusted, Action: "MILK_RECONCILED", Reference: fmt.Sprintf("MILK-JOURNAL-%d", l.JournalID),
				SourceType: "milk_journal", SourceID: l.JournalID,
				Narration: fmt.Sprintf("Reconciled milk journal %s (%s): %.2f kg x KES %.2f = KES %.2f due, KES %.2f credited before",
					firstNonEmptyString(l.Journal, fmt.Sprintf("#%d", l.JournalID)), state, l.Kilos, l.Rate, l.Expected, l.Credited),
				CreatedBy: userID, UpdatedBy: userID, CreatedAt: &now, UpdatedAt: now.Format("2006-01-02 15:04:05"),
			}).Error; err != nil {
				return err
			}
			appendAppLog("credit-changes.log", fmt.Sprintf("member %d | MILK_RECONCILED journal %d %s | due %.2f, credited %.2f | %+.2f | limit %.2f -> %.2f",
				memberID, l.JournalID, l.Journal, l.Expected, l.Credited, l.Adjusted, limit, newLimit))
			limit = newLimit
		}
		// New stored limit, and a fresh baseline so the next event only adds what changes after this.
		b, err := s.Compute(tx, member.ID, member.RouteID)
		if err != nil {
			return err
		}
		borrow := s.EvaluateBorrowStatus(tx, member.ID, member.Status, limit)
		return tx.Model(&models.Member{}).Where("id = ?", member.ID).Updates(map[string]interface{}{
			"credit_limit": limit, "credit_limit_updated_at": now, "credit_limit_computed": b.CreditLimit,
			"can_borrow": borrow.CanBorrow, "cannot_borrow_code": borrow.Code, "cannot_borrow_reason": borrow.Reason,
		}).Error
	})
	return lines, err
}

// runCreditReconcileOnce reconciles every member once (recorded in app_data_migrations).
func (s *MemberCreditService) runCreditReconcileOnce() {
	var done int64
	db.DB.Raw(`SELECT COUNT(*) FROM app_data_migrations WHERE name = ?`, creditReconcileMigration).Scan(&done)
	if done > 0 {
		return
	}
	res, err := s.ReconcileConfirmedMilkCredit(0, 0)
	if err != nil {
		log.Printf("[MemberCredit] reconcile: %v", err)
		return
	}
	db.DB.Exec(`INSERT IGNORE INTO app_data_migrations (name, applied_at) VALUES (?, NOW())`, creditReconcileMigration)
	log.Printf("[MemberCredit] reconciled credit limits with confirmed milk: %d members adjusted, total %+.2f", res.MembersFixed, res.TotalAdjusted)
}

// recalcCreditAfterMilkRateChange re-values every member's confirmed, unpaid milk at the new rates,
// so a rate set after milk was confirmed still reaches the credit limit (log action MILK_RATE_CHANGED).
func recalcCreditAfterMilkRateChange(narration string) {
	ForgetMemberRates()
	var ids []uint64
	if err := db.DB.Model(&models.Member{}).Pluck("id", &ids).Error; err != nil {
		log.Printf("[MemberCredit] rate change recalculation: %v", err)
		return
	}
	appendAppLog("credit-changes.log", narration+fmt.Sprintf(" - recalculating %d members", len(ids)))
	NewMemberCreditService().RecalculateAsync(ids, CreditChange{
		Action:     "MILK_RATE_CHANGED",
		SourceType: "milk_rate",
		Reference:  "MILK-RATE",
		Narration:  narration,
	})
}

// loadAdvanceRateResolver loads the advance rates (default_milk_rates.advance_rate): member default,
// then route default, then the all-members default. Special (period) rates do not apply to advances.
func loadAdvanceRateResolver(tx *gorm.DB) (*milkRateResolver, error) {
	r := newMilkRateResolver()
	var defaults []models.DefaultMilkRate
	if err := tx.Where("deleted_at IS NULL AND advance_rate IS NOT NULL AND advance_rate > 0").Find(&defaults).Error; err != nil {
		return nil, err
	}
	for _, dr := range defaults {
		r.addDefault(dr.MemberID, dr.RouteID, *dr.AdvanceRate)
	}
	return r, nil
}

// MemberAdvanceRate is the advance rate (KES/kg) that applies to the member, 0 when none is set.
func MemberAdvanceRate(tx *gorm.DB, memberID uint64) float64 {
	var routeID uint64
	tx.Raw(`SELECT COALESCE(route_id, 0) FROM member_registrations WHERE id = ?`, memberID).Scan(&routeID)
	rates, err := loadAdvanceRateResolver(tx)
	if err != nil {
		return 0
	}
	return rates.Rate(0, memberID, routeID)
}

const creditAdvanceBasisMigration = "member_credit_two_rate_basis_v2" // v2 also fills member_balance

// runAdvanceBasisSwitchOnce re-values everyone's confirmed milk at the advance rate once (credit was
// valued at the payout milk rate before). Logged as MILK_RATE_CHANGED.
func (s *MemberCreditService) runAdvanceBasisSwitchOnce() {
	var done int64
	db.DB.Raw(`SELECT COUNT(*) FROM app_data_migrations WHERE name = ?`, creditAdvanceBasisMigration).Scan(&done)
	if done > 0 {
		return
	}
	if !db.DB.Migrator().HasColumn(&models.DefaultMilkRate{}, "advance_rate") {
		log.Printf("[MemberCredit] advance_rate column missing; advance basis switch postponed")
		return
	}
	var ids []uint64
	if err := db.DB.Model(&models.Member{}).Pluck("id", &ids).Error; err != nil {
		return
	}
	for _, id := range ids {
		if _, err := s.RecalculateNow(id, CreditChange{
			Action: "MILK_RATE_CHANGED", SourceType: "milk_rate", Reference: "ADVANCE-RATE",
			Narration: "Credit now valued on two rates: milk at the credit rate, loans/deductions converted at credit/default",
		}); err != nil {
			log.Printf("[MemberCredit] advance basis member %d: %v", id, err)
		}
	}
	db.DB.Exec(`INSERT IGNORE INTO app_data_migrations (name, applied_at) VALUES (?, NOW())`, creditAdvanceBasisMigration)
	appendAppLog("credit-changes.log", fmt.Sprintf("credit now uses the ADVANCE rate: %d members re-valued", len(ids)))
}

const creditLoanDeductibleMigration = "member_credit_loan_total_deductible_v1"

// runLoanDeductibleSwitchOnce re-applies outstanding member loans once: their Total Deductible now
// comes off the credit limit at full value (it was converted at credit/default before). The member
// balance does not change. Logged as LOAN_BASIS_CHANGED.
func (s *MemberCreditService) runLoanDeductibleSwitchOnce() {
	var done int64
	db.DB.Raw(`SELECT COUNT(*) FROM app_data_migrations WHERE name = ?`, creditLoanDeductibleMigration).Scan(&done)
	if done > 0 {
		return
	}
	var ids []uint64
	if err := db.DB.Model(&models.LoanContract{}).Where("status IN ?", creditObligationContractStatuses).
		Distinct("member_id").Pluck("member_id", &ids).Error; err != nil {
		log.Printf("[MemberCredit] loan deductible switch: %v", err)
		return
	}
	for _, id := range ids {
		if _, err := s.RecalculateNow(id, CreditChange{
			Action: "LOAN_BASIS_CHANGED", SourceType: "loan", Reference: "TOTAL-DEDUCTIBLE",
			Narration: "Loans now reduce the credit limit by their full Total Deductible",
		}); err != nil {
			log.Printf("[MemberCredit] loan deductible switch member %d: %v", id, err)
		}
	}
	db.DB.Exec(`INSERT IGNORE INTO app_data_migrations (name, applied_at) VALUES (?, NOW())`, creditLoanDeductibleMigration)
	appendAppLog("credit-changes.log", fmt.Sprintf("loans now reduce the credit limit by their Total Deductible: %d members re-applied", len(ids)))
}
