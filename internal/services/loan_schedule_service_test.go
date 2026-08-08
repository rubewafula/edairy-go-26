package services

import (
	"testing"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/models"
)

func TestAmortizingSchedulePrincipalSum(t *testing.T) {
	svc := NewLoanScheduleService()
	lines := svc.Generate(scheduleInput{
		Principal:  100000,
		AnnualRate: 12,
		TermMonths: 12,
		Method:     models.LoanInterestEqualInstallments,
		StartDate:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
	})
	if len(lines) != 12 {
		t.Fatalf("expected 12 installments, got %d", len(lines))
	}
	var principalSum float64
	for _, l := range lines {
		principalSum += l.PrincipalDue
	}
	if diff := principalSum - 100000; diff > 1 || diff < -1 {
		t.Fatalf("principal sum %f expected ~100000", principalSum)
	}
}

func TestFlatScheduleTotalInterest(t *testing.T) {
	svc := NewLoanScheduleService()
	lines := svc.Generate(scheduleInput{
		Principal:  100000,
		AnnualRate: 12,
		TermMonths: 12,
		Method:     models.LoanInterestFlat,
		StartDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	var interestSum float64
	for _, l := range lines {
		interestSum += l.InterestDue
	}
	// flat monthly rate 1% * 12 months on 100k = 12k interest
	if interestSum < 11000 || interestSum > 13000 {
		t.Fatalf("unexpected flat interest total %f", interestSum)
	}
}
