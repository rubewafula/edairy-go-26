package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MilkRejectService struct{}

func NewMilkRejectService() *MilkRejectService {
	return &MilkRejectService{}
}

func (s *MilkRejectService) CreateReject(req dtos.CreateMilkRejectRequest) (*models.MilkReject, error) {
	reject := &models.MilkReject{
		RouteID:             req.RouteID,
		Quantity:            req.Quantity,
		TransactionDate:     utils.ParseDate(req.TransactionDate),
		Confirmed:           req.Confirmed,
		Reason:              req.Reason,
		Description:         req.Description,
		TransporterID:       req.TransporterID,
		CanID:               req.CanID,
		MemberID:            req.MemberID,
		MilkDeliveryShiftID: req.MilkDeliveryShiftID,
	}

	if err := db.DB.Create(reject).Error; err != nil {
		return nil, err
	}
	return reject, nil
}

func (s *MilkRejectService) GetRejects(page, limit int) ([]dtos.MilkRejectResponse, int64, error) {
	var results []dtos.MilkRejectResponse
	var total int64
	db.DB.Model(&models.MilkReject{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			mr.id, r.route_name, mr.quantity, mr.transaction_date, 
			mr.confirmed, mr.reason, mr.description, 
			t.transporter_no AS transporter_name, 
			CONCAT(m.first_name, ' ', m.last_name) AS member_name, 
			mds.name AS milk_delivery_shift,
			mr.created_at
		FROM milk_rejects mr
		LEFT JOIN routes r ON mr.route_id = r.id
		LEFT JOIN transporters t ON mr.transporter_id = t.id
		LEFT JOIN member_registrations m ON mr.member_id = m.id
		LEFT JOIN milk_delivery_shifts mds ON mr.milk_delivery_shift_id = mds.id
		WHERE mr.deleted_at IS NULL
		ORDER BY mr.created_at DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *MilkRejectService) GetReject(id string) (*dtos.MilkRejectResponse, error) {
	var result dtos.MilkRejectResponse
	query := `
		SELECT 
			mr.id, r.route_name, mr.quantity, mr.transaction_date, 
			mr.confirmed, mr.reason, mr.description, 
			t.transporter_no AS transporter_name, 
			CONCAT(m.first_name, ' ', m.last_name) AS member_name, 
			mds.name AS milk_delivery_shift,
			mr.created_at
		FROM milk_rejects mr
		LEFT JOIN routes r ON mr.route_id = r.id
		LEFT JOIN transporters t ON mr.transporter_id = t.id
		LEFT JOIN member_registrations m ON mr.member_id = m.id
		LEFT JOIN milk_delivery_shifts mds ON mr.milk_delivery_shift_id = mds.id
		WHERE mr.id = ? AND mr.deleted_at IS NULL
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

func (s *MilkRejectService) DeleteReject(id string) error {
	return db.DB.Delete(&models.MilkReject{}, id).Error
}
