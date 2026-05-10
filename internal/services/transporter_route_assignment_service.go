package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type TransporterRouteAssignmentService struct{}

func NewTransporterRouteAssignmentService() *TransporterRouteAssignmentService {
	return &TransporterRouteAssignmentService{}
}

func (s *TransporterRouteAssignmentService) CreateAssignment(req dtos.CreateTransporterRouteAssignmentRequest) (*models.TransporterRouteAssignment, error) {
	assignment := &models.TransporterRouteAssignment{
		TransporterID: req.TransporterID,
		RouteID:       req.RouteID,
		StartDate:     utils.ParseDate(req.StartDate),
		EndDate:       utils.ParseDate(req.EndDate),
		Active:        req.Active,
	}

	if err := db.DB.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *TransporterRouteAssignmentService) GetAssignments() ([]dtos.TransporterRouteAssignmentResponse, int64, error) {
	var results []dtos.TransporterRouteAssignmentResponse
	var total int64
	db.DB.Model(&models.TransporterRouteAssignment{}).Count(&total)

	query := `
		SELECT 
			tra.id, tra.transporter_id, t.transporter_no,
			tra.route_id, r.route_name,
			tra.start_date, tra.end_date, tra.active,
			tra.created_at, tra.updated_at
		FROM transporter_route_assignments tra
		LEFT JOIN transporters t ON tra.transporter_id = t.id
		LEFT JOIN routes r ON tra.route_id = r.id
		WHERE tra.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (s *TransporterRouteAssignmentService) GetAssignment(id string) (*dtos.TransporterRouteAssignmentResponse, error) {
	var result dtos.TransporterRouteAssignmentResponse
	query := `
		SELECT 
			tra.id, tra.transporter_id, t.transporter_no,
			tra.route_id, r.route_name,
			tra.start_date, tra.end_date, tra.active,
			tra.created_at, tra.updated_at
		FROM transporter_route_assignments tra
		LEFT JOIN transporters t ON tra.transporter_id = t.id
		LEFT JOIN routes r ON tra.route_id = r.id
		WHERE tra.id = ? AND tra.deleted_at IS NULL
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

func (s *TransporterRouteAssignmentService) UpdateAssignment(id string, req dtos.UpdateTransporterRouteAssignmentRequest) error {
	var assignment models.TransporterRouteAssignment
	if err := db.DB.First(&assignment, id).Error; err != nil {
		return err
	}

	assignment.TransporterID = req.TransporterID
	assignment.RouteID = req.RouteID
	assignment.StartDate = utils.ParseDate(req.StartDate)
	assignment.EndDate = utils.ParseDate(req.EndDate)
	assignment.Active = req.Active
	return db.DB.Save(&assignment).Error
}

func (s *TransporterRouteAssignmentService) DeleteAssignment(id string) error {
	return db.DB.Delete(&models.TransporterRouteAssignment{}, id).Error
}
