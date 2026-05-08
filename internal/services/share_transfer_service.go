package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type ShareTransferService struct{}

func NewShareTransferService() *ShareTransferService {
	return &ShareTransferService{}
}

func (s *ShareTransferService) CreateShareTransfer(req dtos.CreateShareTransferRequest) (*models.ShareTransfer, error) {
	status := req.Status
	if status == "" {
		status = "PENDING"
	}
	transfer := &models.ShareTransfer{
		TransactionID:   req.TransactionID,
		FromMemberID:    req.FromMemberID,
		ToMemberID:      req.ToMemberID,
		ShareUnits:      req.ShareUnits,
		TransferAmount:  req.TransferAmount,
		Status:          status,
		TransactionDate: utils.ParseDate(req.TransactionDate),
		ApprovedBy:      req.ApprovedBy,
		DateApproved:    utils.ParseDate(req.DateApproved),
	}

	if err := db.DB.Create(transfer).Error; err != nil {
		return nil, err
	}
	return transfer, nil
}

func (s *ShareTransferService) GetShareTransfers() ([]dtos.ShareTransferResponse, int64, error) {
	var results []dtos.ShareTransferResponse
	var total int64
	db.DB.Model(&models.ShareTransfer{}).Count(&total)

	query := `
		SELECT 
			st.id, st.transaction_id,
			st.from_member_id, fm.member_no AS from_member_no, fm.first_name AS from_member_first_name, fm.last_name AS from_member_last_name, 
			st.to_member_id, tm.member_no AS to_member_no, tm.first_name AS to_member_first_name, tm.last_name AS to_member_last_name, 
			st.share_units, st.transfer_amount, t.reference AS reference_no, st.status, st.transaction_date,
			st.approved_by, u.name AS approved_by_user_name, st.date_approved,
			st.created_at, st.updated_at
		FROM share_transfers st
		LEFT JOIN member_registrations fm ON st.from_member_id = fm.id
		LEFT JOIN member_registrations tm ON st.to_member_id = tm.id
		LEFT JOIN users u ON st.approved_by = u.id
		LEFT JOIN transactions t ON st.transaction_id = t.id
		WHERE st.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *ShareTransferService) GetShareTransfer(id string) (*dtos.ShareTransferResponse, error) {
	var result dtos.ShareTransferResponse
	query := `
		SELECT 
			st.id, st.transaction_id, 
			st.from_member_id, fm.member_no AS from_member_no, fm.first_name AS from_member_first_name, fm.last_name AS from_member_last_name, 
			st.to_member_id, tm.member_no AS to_member_no, tm.first_name AS to_member_first_name, tm.last_name AS to_member_last_name, 
			st.share_units, st.transfer_amount, t.reference AS reference_no, st.status, st.transaction_date,
			st.approved_by, u.name AS approved_by_user_name, st.date_approved,
			st.created_at, st.updated_at
		FROM share_transfers st
		LEFT JOIN member_registrations fm ON st.from_member_id = fm.id
		LEFT JOIN member_registrations tm ON st.to_member_id = tm.id
		LEFT JOIN users u ON st.approved_by = u.id
		LEFT JOIN transactions t ON st.transaction_id = t.id
		WHERE st.id = ? AND st.deleted_at IS NULL
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

func (s *ShareTransferService) UpdateShareTransfer(id string, req dtos.UpdateShareTransferRequest) error {
	var transfer models.ShareTransfer
	if err := db.DB.First(&transfer, id).Error; err != nil {
		return err
	}

	transfer.TransactionID = req.TransactionID
	transfer.ShareUnits = req.ShareUnits
	transfer.TransferAmount = req.TransferAmount
	transfer.Status = req.Status
	transfer.TransactionDate = utils.ParseDate(req.TransactionDate)
	transfer.ApprovedBy = req.ApprovedBy
	transfer.DateApproved = utils.ParseDate(req.DateApproved)

	return db.DB.Save(&transfer).Error
}

func (s *ShareTransferService) DeleteShareTransfer(id string) error {
	return db.DB.Delete(&models.ShareTransfer{}, id).Error
}
