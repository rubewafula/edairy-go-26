package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type SharePaymentService struct{}

func NewSharePaymentService() *SharePaymentService {
	return &SharePaymentService{}
}

func (s *SharePaymentService) CreateSharePayment(req dtos.CreateSharePaymentRequest) (*models.SharePayment, error) {
	status := req.Status
	if status == "" {
		status = "PENDING"
	}
	payment := &models.SharePayment{
		TransactionID:   req.TransactionID,
		MemberID:        req.MemberID,
		ShareAccountID:  req.ShareAccountID,
		AmountPaid:      req.AmountPaid,
		ShareUnits:      req.ShareUnits,
		PaymentModeID:   req.PaymentModeID,
		Description:     req.Description,
		Status:          status,
		TransactionDate: utils.ParseDate(req.TransactionDate),
		ApprovedBy:      req.ApprovedBy,
		DateApproved:    utils.ParseDate(req.DateApproved),
	}

	if err := db.DB.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *SharePaymentService) GetSharePayments() ([]dtos.SharePaymentResponse, int64, error) {
	var results []dtos.SharePaymentResponse
	var total int64
	db.DB.Model(&models.SharePayment{}).Count(&total)

	query := `
		SELECT 
			sp.id, sp.transaction_id, sp.member_id, m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			sp.share_account_id, sp.amount_paid, sp.share_units, sp.payment_mode_id, pm.name AS payment_mode_name,
			t.reference AS reference_no, sp.description, sp.status, sp.transaction_date, sp.approved_by, u.name AS approved_by_user_name,
			sp.date_approved, sp.created_at, sp.updated_at
		FROM share_payments sp
		LEFT JOIN member_registrations m ON sp.member_id = m.id
		LEFT JOIN payment_modes pm ON sp.payment_mode_id = pm.id
		LEFT JOIN users u ON sp.approved_by = u.id
		LEFT JOIN transactions t ON sp.transaction_id = t.id
		WHERE sp.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *SharePaymentService) GetSharePayment(id string) (*dtos.SharePaymentResponse, error) {
	var result dtos.SharePaymentResponse
	query := `
		SELECT 
			sp.id, sp.transaction_id, sp.member_id, m.member_no, m.first_name AS member_first_name, m.last_name AS member_last_name,
			sp.share_account_id, sp.amount_paid, sp.share_units, sp.payment_mode_id, pm.name AS payment_mode_name,
			t.reference AS reference_no, sp.description, sp.status, sp.transaction_date, sp.approved_by, u.name AS approved_by_user_name,
			sp.date_approved, sp.created_at, sp.updated_at
		FROM share_payments sp
		LEFT JOIN member_registrations m ON sp.member_id = m.id
		LEFT JOIN payment_modes pm ON sp.payment_mode_id = pm.id
		LEFT JOIN users u ON sp.approved_by = u.id
		LEFT JOIN transactions t ON sp.transaction_id = t.id
		WHERE sp.id = ? AND sp.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *SharePaymentService) UpdateSharePayment(id string, req dtos.UpdateSharePaymentRequest) error {
	var payment models.SharePayment
	if err := db.DB.First(&payment, id).Error; err != nil {
		return err
	}

	payment.TransactionID = req.TransactionID
	payment.ShareAccountID = req.ShareAccountID
	payment.AmountPaid = req.AmountPaid
	payment.ShareUnits = req.ShareUnits
	payment.PaymentModeID = req.PaymentModeID
	payment.Description = req.Description
	payment.Status = req.Status
	payment.TransactionDate = utils.ParseDate(req.TransactionDate)
	payment.ApprovedBy = req.ApprovedBy
	payment.DateApproved = utils.ParseDate(req.DateApproved)

	return db.DB.Save(&payment).Error
}

func (s *SharePaymentService) DeleteSharePayment(id string) error {
	return db.DB.Delete(&models.SharePayment{}, id).Error
}
