package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type CanMovementService struct{}

func NewCanMovementService() *CanMovementService {
	return &CanMovementService{}
}

func (s *CanMovementService) CreateMovement(req dtos.CreateCanMovementRequest) (*models.CanMovement, error) {
	movement := &models.CanMovement{
		CanID:             req.CanID,
		MovementType:      req.MovementType,
		Quantity:          req.Quantity,
		Remarks:           req.Remarks,
		ShiftID:           req.ShiftID,
		TransporterID:     req.TransporterID,
		RouteID:           req.RouteID,
		MovementDate:      utils.ParseDate(req.MovementDate),
		ConditionOnReturn: req.ConditionOnReturn,
	}

	if err := db.DB.Create(movement).Error; err != nil {
		return nil, err
	}
	return movement, nil
}

func (s *CanMovementService) GetMovements(page, limit int) ([]dtos.CanMovementResponse, int64, error) {
	var results []dtos.CanMovementResponse
	var total int64
	db.DB.Model(&models.CanMovement{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			cm.id, cm.can_id, mc.can_id AS can_code, cm.movement_type, cm.quantity, cm.remarks,
			cm.shift_id, mds.name AS shift_name,
			cm.transporter_id, t.transporter_no,
			cm.route_id, r.route_name,
			cm.movement_date, cm.condition_on_return,
			cm.created_at, cm.updated_at
		FROM can_movements cm
		LEFT JOIN milk_cans mc ON cm.can_id = mc.id
		LEFT JOIN milk_delivery_shifts mds ON cm.shift_id = mds.id
		LEFT JOIN transporters t ON cm.transporter_id = t.id
		LEFT JOIN routes r ON cm.route_id = r.id
		WHERE cm.deleted_at IS NULL
		ORDER BY cm.movement_date DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *CanMovementService) GetMovement(id string) (*dtos.CanMovementResponse, error) {
	var result dtos.CanMovementResponse
	query := `
		SELECT 
			cm.id, cm.can_id, mc.can_id AS can_code, cm.movement_type, cm.quantity, cm.remarks,
			cm.shift_id, mds.name AS shift_name,
			cm.transporter_id, t.transporter_no,
			cm.route_id, r.route_name,
			cm.movement_date, cm.condition_on_return,
			cm.created_at, cm.updated_at
		FROM can_movements cm
		LEFT JOIN milk_cans mc ON cm.can_id = mc.id
		LEFT JOIN milk_delivery_shifts mds ON cm.shift_id = mds.id
		LEFT JOIN transporters t ON cm.transporter_id = t.id
		LEFT JOIN routes r ON cm.route_id = r.id
		WHERE cm.id = ? AND cm.deleted_at IS NULL
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

func (s *CanMovementService) UpdateMovement(id string, req dtos.UpdateCanMovementRequest) error {
	var movement models.CanMovement
	if err := db.DB.First(&movement, id).Error; err != nil {
		return err
	}

	movement.CanID = req.CanID
	movement.MovementType = req.MovementType
	movement.Quantity = req.Quantity
	movement.Remarks = req.Remarks
	movement.ShiftID = req.ShiftID
	movement.TransporterID = req.TransporterID
	movement.RouteID = req.RouteID
	movement.MovementDate = utils.ParseDate(req.MovementDate)
	movement.ConditionOnReturn = req.ConditionOnReturn

	return db.DB.Save(&movement).Error
}

func (s *CanMovementService) DeleteMovement(id string) error {
	return db.DB.Delete(&models.CanMovement{}, id).Error
}
