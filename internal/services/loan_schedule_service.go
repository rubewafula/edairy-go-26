package services

import (
	"encoding/json"
	"math"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/models"
)

type scheduleInput struct {
	Principal       float64
	AnnualRate      float64
	TermMonths      int
	Method          string
	GracePeriodDays int
	StartDate       time.Time
	ProcessingFee   float64
	InsuranceFee    float64
}

type scheduleLine struct {
	InstallmentNo int
	DueDate       time.Time
	PrincipalDue  float64
	InterestDue   float64
	FeeDue        float64
	InsuranceDue  float64
	TotalDue      float64
	BalanceAfter  float64
}

type LoanScheduleService struct{}

func NewLoanScheduleService() *LoanScheduleService {
	return &LoanScheduleService{}
}

func (s *LoanScheduleService) Generate(input scheduleInput) []scheduleLine {
	switch input.Method {
	case models.LoanInterestFlat:
		return s.flatSchedule(input)
	case models.LoanInterestOnly:
		return s.interestOnlySchedule(input)
	case models.LoanInterestBalloon:
		return s.balloonSchedule(input)
	case models.LoanInterestEqualInstallments, models.LoanInterestReducingBalance:
		return s.amortizingSchedule(input)
	default:
		return s.amortizingSchedule(input)
	}
}

func (s *LoanScheduleService) firstDueDate(start time.Time, graceDays int) time.Time {
	d := start
	if graceDays > 0 {
		d = d.AddDate(0, 0, graceDays)
	}
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location()).AddDate(0, 1, 0)
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}

func (s *LoanScheduleService) flatSchedule(input scheduleInput) []scheduleLine {
	monthlyRate := input.AnnualRate / 100 / 12
	totalInterest := roundMoney(input.Principal * monthlyRate * float64(input.TermMonths))
	totalPrincipal := input.Principal
	installmentPrincipal := roundMoney(totalPrincipal / float64(input.TermMonths))
	installmentInterest := roundMoney(totalInterest / float64(input.TermMonths))
	feePer := roundMoney(input.ProcessingFee / float64(input.TermMonths))
	insPer := roundMoney(input.InsuranceFee / float64(input.TermMonths))

	lines := make([]scheduleLine, 0, input.TermMonths)
	balance := totalPrincipal
	due := s.firstDueDate(input.StartDate, input.GracePeriodDays)
	for i := 1; i <= input.TermMonths; i++ {
		principal := installmentPrincipal
		if i == input.TermMonths {
			principal = roundMoney(balance)
		}
		balance = roundMoney(balance - principal)
		if balance < 0 {
			balance = 0
		}
		lines = append(lines, scheduleLine{
			InstallmentNo: i,
			DueDate:       due,
			PrincipalDue:  principal,
			InterestDue:   installmentInterest,
			FeeDue:        feePer,
			InsuranceDue:  insPer,
			TotalDue:      roundMoney(principal + installmentInterest + feePer + insPer),
			BalanceAfter:  balance,
		})
		due = due.AddDate(0, 1, 0)
	}
	return lines
}

func (s *LoanScheduleService) amortizingSchedule(input scheduleInput) []scheduleLine {
	monthlyRate := input.AnnualRate / 100 / 12
	pmt := input.Principal
	if monthlyRate > 0 {
		pmt = input.Principal * (monthlyRate * math.Pow(1+monthlyRate, float64(input.TermMonths))) / (math.Pow(1+monthlyRate, float64(input.TermMonths)) - 1)
	} else {
		pmt = input.Principal / float64(input.TermMonths)
	}
	pmt = roundMoney(pmt)

	lines := make([]scheduleLine, 0, input.TermMonths)
	balance := input.Principal
	due := s.firstDueDate(input.StartDate, input.GracePeriodDays)
	for i := 1; i <= input.TermMonths; i++ {
		interest := roundMoney(balance * monthlyRate)
		principal := roundMoney(pmt - interest)
		if i == input.TermMonths {
			principal = roundMoney(balance)
			pmt = roundMoney(principal + interest)
		}
		balance = roundMoney(balance - principal)
		if balance < 0 {
			balance = 0
		}
		lines = append(lines, scheduleLine{
			InstallmentNo: i,
			DueDate:       due,
			PrincipalDue:  principal,
			InterestDue:   interest,
			TotalDue:      pmt,
			BalanceAfter:  balance,
		})
		due = due.AddDate(0, 1, 0)
	}
	return lines
}

func (s *LoanScheduleService) interestOnlySchedule(input scheduleInput) []scheduleLine {
	monthlyRate := input.AnnualRate / 100 / 12
	monthlyInterest := roundMoney(input.Principal * monthlyRate)
	lines := make([]scheduleLine, 0, input.TermMonths)
	due := s.firstDueDate(input.StartDate, input.GracePeriodDays)
	for i := 1; i <= input.TermMonths; i++ {
		principal := 0.0
		if i == input.TermMonths {
			principal = input.Principal
		}
		lines = append(lines, scheduleLine{
			InstallmentNo: i,
			DueDate:       due,
			PrincipalDue:  principal,
			InterestDue:   monthlyInterest,
			TotalDue:      roundMoney(principal + monthlyInterest),
			BalanceAfter:  roundMoney(input.Principal - principal),
		})
		due = due.AddDate(0, 1, 0)
	}
	return lines
}

func (s *LoanScheduleService) balloonSchedule(input scheduleInput) []scheduleLine {
	monthlyRate := input.AnnualRate / 100 / 12
	monthlyInterest := roundMoney(input.Principal * monthlyRate)
	smallPrincipal := roundMoney(input.Principal * 0.05)
	lines := make([]scheduleLine, 0, input.TermMonths)
	balance := input.Principal
	due := s.firstDueDate(input.StartDate, input.GracePeriodDays)
	for i := 1; i <= input.TermMonths; i++ {
		principal := smallPrincipal
		if i == input.TermMonths {
			principal = roundMoney(balance)
		}
		balance = roundMoney(balance - principal)
		lines = append(lines, scheduleLine{
			InstallmentNo: i,
			DueDate:       due,
			PrincipalDue:  principal,
			InterestDue:   monthlyInterest,
			TotalDue:      roundMoney(principal + monthlyInterest),
			BalanceAfter:  balance,
		})
		due = due.AddDate(0, 1, 0)
	}
	return lines
}

func defaultAllocationOrder(product *models.LoanProduct) []string {
	if product != nil && len(product.AllocationPriority) > 0 {
		var order []string
		if err := json.Unmarshal(product.AllocationPriority, &order); err == nil && len(order) > 0 {
			return order
		}
	}
	return []string{"PENALTY", "FEE", "INTEREST", "PRINCIPAL", "INSURANCE"}
}
