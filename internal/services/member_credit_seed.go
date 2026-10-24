package services

import (
	"fmt"
	"math"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

// CreditSeedResult summarises a SeedMinimumCreditLimits run.
type CreditSeedResult struct {
	Members       int      // members processed
	Raised        int      // members topped up to the minimum
	AlreadyAbove  int      // members already at or above the minimum
	KilosSeeded   float64  // total unpaid milk seeded
	RateCreated   bool     // a global default milk rate had to be created
	Failures      []string // per-member errors (that member is rolled back)
	JournalPrefix string
}

// SeedMinimumCreditLimits gives every member a credit limit of at least minLimit. DEV/TEST DATA ONLY.
//
// The credit limit is derived (unpaid confirmed milk value - rejects - loan obligations - deductions)
// and recomputed on every credit query, so writing members.credit_limit directly would be overwritten.
// Instead, for each member below the minimum this seeds a confirmed milk journal entry dated today for
// the shortfall, then runs Recalculate, which stores the new limit and writes the
// credit_limit_change_log row (old limit, new limit, delta, action SEEDED, journal reference).
// Members already above the minimum are recalculated too, so stored limits and logs match source data.
//
// The seeded milk is unpaid and WILL be paid out by the next member payroll; never run this against
// a live dairy database.
func SeedMinimumCreditLimits(minLimit float64, fallbackRate float64, userID uint64) (*CreditSeedResult, error) {
	return SeedMinimumCreditLimitsFor(nil, minLimit, fallbackRate, userID)
}

// SeedMinimumCreditLimitsFor is SeedMinimumCreditLimits limited to memberIDs (nil or empty = all members).
func SeedMinimumCreditLimitsFor(memberIDs []uint64, minLimit float64, fallbackRate float64, userID uint64) (*CreditSeedResult, error) {
	if minLimit <= 0 {
		return nil, fmt.Errorf("minimum credit limit must be greater than zero")
	}
	if _, err := EnsureRouteCenters(); err != nil {
		return nil, fmt.Errorf("ensure route centers: %w", err)
	}

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	res := &CreditSeedResult{JournalPrefix: "CREDIT-SEED-" + today.Format("20060102")}

	shift, err := seedShift(userID)
	if err != nil {
		return nil, err
	}
	var fallbackRouteID uint64
	db.DB.Model(&models.Route{}).Order("id ASC").Limit(1).Pluck("id", &fallbackRouteID)

	var members []models.Member
	q := db.DB.Select("id, member_no, route_id").Order("id ASC")
	if len(memberIDs) > 0 {
		q = q.Where("id IN ?", memberIDs)
	}
	if err := q.Find(&members).Error; err != nil {
		return nil, err
	}

	svc := NewMemberCreditService()
	for _, m := range members {
		res.Members++
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			b, err := svc.Compute(tx, m.ID, m.RouteID)
			if err != nil {
				return err
			}
			if b.CreditLimit >= minLimit {
				res.AlreadyAbove++
				_, err := svc.Recalculate(tx, m.ID, CreditChange{
					Action:     models.CreditActionRecalculated,
					SourceType: "credit_seed",
					Narration:  "Credit limit recalculated during credit limit seed",
					UserID:     userID,
				})
				return err
			}

			routeID := m.RouteID
			if routeID == 0 {
				routeID = fallbackRouteID
			}
			if routeID == 0 {
				return fmt.Errorf("no route available")
			}

			rate, created, err := seedMilkRate(tx, m.ID, routeID, today, fallbackRate, userID)
			if err != nil {
				return err
			}
			if created {
				res.RateCreated = true
			}

			shortfall := minLimit - b.CreditLimit
			kilos := math.Ceil(shortfall/rate*100) / 100

			journal, batch, centerID, err := seedJournal(tx, res.JournalPrefix, today, routeID, shift.ID, userID)
			if err != nil {
				return err
			}
			if err := tx.Create(&models.MilkJournalEntry{
				BaseModel:          models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
				MemberID:           m.ID,
				MilkJournalID:      journal.ID,
				MilkJournalBatchID: batch.ID,
				Status:             "PENDING",
				Quantity:           kilos,
				RouteCenterID:      centerID,
			}).Error; err != nil {
				return err
			}

			after, err := svc.Recalculate(tx, m.ID, CreditChange{
				Action:     models.CreditActionSeeded,
				SourceType: "milk_journal",
				SourceID:   journal.ID,
				Reference:  journal.Journal,
				Narration:  fmt.Sprintf("Seeded %.2f kg unpaid milk @ KES %.2f to reach minimum credit limit KES %.2f", kilos, rate, minLimit),
				UserID:     userID,
			})
			if err != nil {
				return err
			}
			if after.CreditLimit < minLimit-0.01 {
				// e.g. today's date falls in a pay period whose payroll is already approved
				return fmt.Errorf("limit only reached KES %.2f after seeding; is today's pay period already in an approved payroll?", after.CreditLimit)
			}
			res.Raised++
			res.KilosSeeded += kilos
			return nil
		})
		if err != nil {
			res.Failures = append(res.Failures, fmt.Sprintf("member %s (id %d): %v", m.MemberNo, m.ID, err))
		}
	}
	return res, nil
}

func seedShift(userID uint64) (*models.MilkDeliveryShift, error) {
	var shift models.MilkDeliveryShift
	if err := db.DB.Order("id ASC").Limit(1).Find(&shift).Error; err != nil {
		return nil, err
	}
	if shift.ID != 0 {
		return &shift, nil
	}
	shift = models.MilkDeliveryShift{
		BaseModel:   models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		Name:        "Morning",
		Description: "Morning milk collection shift",
	}
	return &shift, db.DB.Create(&shift).Error
}

// seedMilkRate returns the rate the credit calculation will apply to today's milk for this member.
// If no rate resolves, it creates a global default milk rate of fallbackRate.
func seedMilkRate(tx *gorm.DB, memberID, routeID uint64, day time.Time, fallbackRate float64, userID uint64) (float64, bool, error) {
	var pdrID uint64
	tx.Raw(`SELECT id FROM member_pay_date_ranges WHERE ? BETWEEN start_date AND end_date AND deleted_at IS NULL ORDER BY id DESC LIMIT 1`,
		day.Format("2006-01-02")).Scan(&pdrID)
	var ids []uint64
	if pdrID > 0 {
		ids = []uint64{pdrID}
	}
	rates, err := loadAdvanceRateResolver(tx)
	_ = ids
	if err != nil {
		return 0, false, err
	}
	if rate := rates.Rate(pdrID, memberID, routeID); rate > 0 {
		return rate, false, nil
	}
	if fallbackRate <= 0 {
		return 0, false, fmt.Errorf("no milk rate configured for this member/route; set a default milk rate first")
	}
	if err := tx.Create(&models.DefaultMilkRate{
		BaseModel: models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
		Rate:      fallbackRate,
	}).Error; err != nil {
		return 0, false, err
	}
	return fallbackRate, true, nil
}

// seedJournal finds or creates today's confirmed seed journal for the route, its batch, and the route's first center.
func seedJournal(tx *gorm.DB, prefix string, day time.Time, routeID, shiftID, userID uint64) (*models.MilkJournal, *models.MilkJournalBatch, uint64, error) {
	name := fmt.Sprintf("%s-R%d", prefix, routeID)
	var journal models.MilkJournal
	if err := tx.Where("journal = ? AND route_id = ?", name, routeID).Limit(1).Find(&journal).Error; err != nil {
		return nil, nil, 0, err
	}
	if journal.ID == 0 {
		journal = models.MilkJournal{
			BaseModel:           models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			Journal:             name,
			JournalDate:         &day,
			MilkDeliveryShiftID: shiftID,
			RouteID:             routeID,
			UserID:              userID,
			Confirmed:           true,
		}
		if err := tx.Create(&journal).Error; err != nil {
			return nil, nil, 0, err
		}
	}
	var batch models.MilkJournalBatch
	if err := tx.Where("milk_journal_id = ? AND batch_no = ?", journal.ID, "CREDIT-SEED").Limit(1).Find(&batch).Error; err != nil {
		return nil, nil, 0, err
	}
	if batch.ID == 0 {
		batch = models.MilkJournalBatch{
			BaseModel:     models.BaseModel{CreatedBy: userID, UpdatedBy: userID},
			MilkJournalID: journal.ID,
			BatchNo:       "CREDIT-SEED",
		}
		if err := tx.Create(&batch).Error; err != nil {
			return nil, nil, 0, err
		}
	}
	var center models.RouteCenter
	if err := tx.Where("route_id = ?", routeID).Order("id ASC").First(&center).Error; err != nil {
		return nil, nil, 0, fmt.Errorf("route %d has no center", routeID)
	}
	return &journal, &batch, center.ID, nil
}
