package services

import (
	"github.com/google/uuid"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type LoanService struct{}

func NewLoanService() *LoanService {
	return &LoanService{}
}

func (s *LoanService) CreateLoan(req dtos.CreateLoanRequest) (*models.Loan, error) {
	status := req.Status
	if status == "" {
		status = "PENDING"
	}
	loan := &models.Loan{
		MemberID:              req.MemberID,
		Amount:                req.Amount,
		Interest:              req.Interest,
		TotalPayable:          req.TotalPayable,
		Status:                status,
		UUID:                  uuid.New().String(),
		RequestID:             req.RequestID,
		RepaymentAmount:       req.RepaymentAmount,
		WithdrawalRequestUUID: req.WithdrawalRequestUUID,
	}

	if err := db.DB.Create(loan).Error; err != nil {
		return nil, err
	}
	return loan, nil
}

func (s *LoanService) GetLoans() ([]models.Loan, int64, error) {
	var loans []models.Loan
	var total int64
	db.DB.Model(&models.Loan{}).Count(&total)
	err := db.DB.Find(&loans).Error
	return loans, total, err
}

func (s *LoanService) GetLoan(id string) (*models.Loan, error) {
	var loan models.Loan
	if err := db.DB.First(&loan, id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

func (s *LoanService) UpdateLoan(id string, req dtos.UpdateLoanRequest) error {
	var loan models.Loan
	if err := db.DB.First(&loan, id).Error; err != nil {
		return err
	}

	loan.Amount = req.Amount
	loan.Interest = req.Interest
	loan.TotalPayable = req.TotalPayable
	loan.Status = req.Status
	loan.ApprovedAmt = req.ApprovedAmt
	loan.ProcessedBy = req.ProcessedBy
	loan.LoanLimitBy = req.LoanLimitBy
	loan.ReviewAccepted = req.ReviewAccepted
	loan.CreditLimit = req.CreditLimit
	loan.DisbursedAt = utils.ParseDate(req.DisbursedAt)
	loan.ProcessedAt = utils.ParseDate(req.ProcessedAt)
	loan.TotalDue = req.TotalDue
	loan.RepaymentAmount = req.RepaymentAmount
	loan.RequestID = req.RequestID
	loan.WithdrawalRequestUUID = req.WithdrawalRequestUUID

	return db.DB.Save(&loan).Error
}

func (s *LoanService) DeleteLoan(id string) error {
	return db.DB.Delete(&models.Loan{}, id).Error
}
