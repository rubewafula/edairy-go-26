package services

import (
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
	return roundMoney(principal * product.ProcessingFeeRate / 100),
		roundMoney(principal * product.InsuranceFeeRate / 100)
}

func feesSpreadOverInstallments(method string) bool {
	return normalizeFeeCollectionMethod(method) == models.LoanFeeSpreadOverInstallments
}

func feesDeductedFromProceeds(method string) bool {
	return normalizeFeeCollectionMethod(method) == models.LoanFeeDeductFromProceeds
}

func (s *LoanModuleService) isContractFullyDisbursed(contract *models.LoanContract, product *models.LoanProduct) bool {
	if feesDeductedFromProceeds(product.FeeCollectionMethod) {
		return roundMoney(contract.DisbursedAmount+contract.FeesDeductedAtDisbursement) >= contract.ApprovedAmount
	}
	return contract.DisbursedAmount >= contract.ApprovedAmount
}

func (s *LoanModuleService) remainingDisbursementNet(contract *models.LoanContract, product *models.LoanProduct) float64 {
	if feesDeductedFromProceeds(product.FeeCollectionMethod) {
		return roundMoney(contract.ApprovedAmount - contract.DisbursedAmount - contract.FeesDeductedAtDisbursement)
	}
	return roundMoney(contract.ApprovedAmount - contract.DisbursedAmount)
}
