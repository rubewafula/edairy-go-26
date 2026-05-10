package services

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type TransporterDriverAssignmentService struct{}

func NewTransporterDriverAssignmentService() *TransporterDriverAssignmentService {
	return &TransporterDriverAssignmentService{}
}

func (s *TransporterDriverAssignmentService) CreateAssignment(req dtos.CreateTransporterDriverAssignmentRequest) (*models.TransporterDriverAssignment, error) {
	assignedFrom, _ := time.Parse(time.RFC3339, req.AssignedFrom)
	var assignedTo *time.Time
	if req.AssignedTo != "" {
		t, _ := time.Parse(time.RFC3339, req.AssignedTo)
		assignedTo = &t
	}

	assignment := &models.TransporterDriverAssignment{
		TransporterDriverID:  req.TransporterDriverID,
		TransporterVehicleID: req.TransporterVehicleID,
		AssignedFrom:         assignedFrom,
		AssignedTo:           assignedTo,
		AssignmentType:       req.AssignmentType,
		Active:               req.Active,
		Notes:                req.Notes,
	}

	if err := db.DB.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *TransporterDriverAssignmentService) GetAssignments() ([]dtos.TransporterDriverAssignmentResponse, int64, error) {
	var results []dtos.TransporterDriverAssignmentResponse
	var total int64
	db.DB.Model(&models.TransporterDriverAssignment{}).Count(&total)

	query := `
		SELECT 
			tda.id, tda.transporter_driver_id, CONCAT(td.first_name, ' ', td.last_name) AS driver_name, td.driver_no,
			tda.transporter_vehicle_id, tv.registration_no AS vehicle_reg_no,
			tda.assigned_from, tda.assigned_to, tda.assignment_type, tda.active, tda.notes,
			tda.created_at
		FROM transporter_driver_assignments tda
		LEFT JOIN transporter_drivers td ON tda.transporter_driver_id = td.id
		LEFT JOIN transporter_vehicles tv ON tda.transporter_vehicle_id = tv.id
		WHERE tda.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *TransporterDriverAssignmentService) GetAssignment(id string) (*dtos.TransporterDriverAssignmentResponse, error) {
	var result dtos.TransporterDriverAssignmentResponse
	query := `
		SELECT 
			tda.id, tda.transporter_driver_id, CONCAT(td.first_name, ' ', td.last_name) AS driver_name, td.driver_no,
			tda.transporter_vehicle_id, tv.registration_no AS vehicle_reg_no,
			tda.assigned_from, tda.assigned_to, tda.assignment_type, tda.active, tda.notes,
			tda.created_at
		FROM transporter_driver_assignments tda
		LEFT JOIN transporter_drivers td ON tda.transporter_driver_id = td.id
		LEFT JOIN transporter_vehicles tv ON tda.transporter_vehicle_id = tv.id
		WHERE tda.id = ? AND tda.deleted_at IS NULL
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

func (s *TransporterDriverAssignmentService) UpdateAssignment(id string, req dtos.UpdateTransporterDriverAssignmentRequest) error {
	var assignment models.TransporterDriverAssignment
	if err := db.DB.First(&assignment, id).Error; err != nil {
		return err
	}

	if req.AssignedTo != "" {
		t, _ := time.Parse(time.RFC3339, req.AssignedTo)
		assignment.AssignedTo = &t
	}
	assignment.AssignmentType = req.AssignmentType
	assignment.Active = req.Active
	assignment.Notes = req.Notes
	return db.DB.Save(&assignment).Error
}

func (s *TransporterDriverAssignmentService) DeleteAssignment(id string) error {
	return db.DB.Delete(&models.TransporterDriverAssignment{}, id).Error
}
