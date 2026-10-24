package services

import (
	"strings"
	"sort"
	"fmt"
	"encoding/json"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

func normalizeFeeCollectionMethod(method string) string {
	switch method {
	case models.LoanFeePayUpfrontCash, models.LoanFeeSpreadOverInstallments:
		return method
	default:
		return models.LoanFeeDeductFromProceeds
	}
}

func computeProductFees(product *models.LoanProduct, principal float64) (processingFee, insuranceFee float64) {
	if product == nil || principal <= 0 {
		return 0, 0
	}
	// The default member loan (mobile cashouts) uses the product's processing fee and insurance like any
	// other product; only its interest rate is ignored (replaced by the milk-rate spread, see spreadInterest).
	// Fixed fee per loan (KES); the percentage is only a fallback for products not yet given one.
	processing := product.ProcessingFeeAmount
	if processing <= 0 {
		processing = principal * product.ProcessingFeeRate / 100
	}
	if fee, ok := tierProcessingFee(product.FeeTiers, principal); ok {
		processing = fee
	}
	return roundMoney(processing), roundMoney(principal * product.InsuranceFeeRate / 100)
}

// tierProcessingFee returns the processing fee for the band that contains amount.
func tierProcessingFee(tiers []models.ProcessingFeeTier, amount float64) (float64, bool) {
	for _, t := range tiers {
		if amount+0.005 >= t.MinAmount && (t.MaxAmount <= 0 || amount-0.005 <= t.MaxAmount) {
			if strings.EqualFold(t.ChargeType, "PERCENT") {
				return amount * t.Charge / 100, true
			}
			return t.Charge, true
		}
	}
	return 0, false
}

// validateFeeTiers checks bands are well-formed and do not overlap; returns them sorted.
func validateFeeTiers(tiers []models.ProcessingFeeTier) ([]models.ProcessingFeeTier, error) {
	out := make([]models.ProcessingFeeTier, 0, len(tiers))
	for i, t := range tiers {
		t.ChargeType = strings.ToUpper(strings.TrimSpace(t.ChargeType))
		if t.ChargeType == "" {
			t.ChargeType = "FLAT"
		}
		if t.ChargeType != "FLAT" && t.ChargeType != "PERCENT" {
			return nil, fmt.Errorf("fee band %d: charge type must be FLAT or PERCENT", i+1)
		}
		if t.MinAmount < 0 || t.Charge < 0 {
			return nil, fmt.Errorf("fee band %d: amounts cannot be negative", i+1)
		}
		if t.MaxAmount > 0 && t.MaxAmount < t.MinAmount {
			return nil, fmt.Errorf("fee band %d: 'to' amount must be at least the 'from' amount", i+1)
		}
		t.MinAmount, t.MaxAmount, t.Charge = roundMoney(t.MinAmount), roundMoney(t.MaxAmount), roundMoney(t.Charge)
		out = append(out, t)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].MinAmount < out[b].MinAmount })
	for i := 1; i < len(out); i++ {
		prev := out[i-1]
		if prev.MaxAmount <= 0 || out[i].MinAmount <= prev.MaxAmount {
			return nil, fmt.Errorf("fee bands overlap: KES %s - %s and KES %s - %s",
				formatKES(prev.MinAmount), bandMax(prev.MaxAmount), formatKES(out[i].MinAmount), bandMax(out[i].MaxAmount))
		}
	}
	return out, nil
}

func bandMax(v float64) string {
	if v <= 0 {
		return "and above"
	}
	return formatKES(v)
}

func encodeFeeTiers(tiers []models.ProcessingFeeTier) string {
	if len(tiers) == 0 {
		return ""
	}
	b, _ := json.Marshal(tiers)
	return string(b)
}

// effectiveFeeCollectionMethod is how a product's fees are actually collected. The default member loan
// (mobile cashouts) always adds its charges to the Total Deductible as one charge on the first
// instalment (recovered from milk), whatever collection method is saved, so the member receives the
// full amount requested. Other products use their saved method.
func effectiveFeeCollectionMethod(product *models.LoanProduct) string {
	if product != nil && product.IsDefaultMemberLoan {
		return models.LoanFeeSpreadOverInstallments
	}
	if product == nil {
		return models.LoanFeeDeductFromProceeds
	}
	return normalizeFeeCollectionMethod(product.FeeCollectionMethod)
}

// feesAddedToRepayment: fees are part of what is repaid (not taken off the payout).
func feesAddedToRepayment(product *models.LoanProduct) bool {
	return effectiveFeeCollectionMethod(product) == models.LoanFeeSpreadOverInstallments
}

// feesOffPayout: fees are taken off the payout (deduct from proceeds).
func feesOffPayout(product *models.LoanProduct) bool {
	return effectiveFeeCollectionMethod(product) == models.LoanFeeDeductFromProceeds
}

func feesSpreadOverInstallments(method string) bool {
	return normalizeFeeCollectionMethod(method) == models.LoanFeeSpreadOverInstallments
}

func feesDeductedFromProceeds(method string) bool {
	return normalizeFeeCollectionMethod(method) == models.LoanFeeDeductFromProceeds
}

func (s *LoanModuleService) isContractFullyDisbursed(contract *models.LoanContract, product *models.LoanProduct) bool {
	if feesOffPayout(product) {
		return roundMoney(contract.DisbursedAmount+contract.FeesDeductedAtDisbursement) >= contract.ApprovedAmount
	}
	return contract.DisbursedAmount >= contract.ApprovedAmount
}

func (s *LoanModuleService) remainingDisbursementNet(contract *models.LoanContract, product *models.LoanProduct) float64 {
	if feesOffPayout(product) {
		return roundMoney(contract.ApprovedAmount - contract.DisbursedAmount - contract.FeesDeductedAtDisbursement)
	}
	return roundMoney(contract.ApprovedAmount - contract.DisbursedAmount)
}
