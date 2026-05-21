package services

import (
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type AssetAssignmentService struct{}

func NewAssetAssignmentService() *AssetAssignmentService {
	return &AssetAssignmentService{}
}

func (s *AssetAssignmentService) CreateAssignment(req dtos.CreateAssetAssignmentRequest) (*models.AssetAssignment, error) {
	status := req.Status
	if status == "" {
		status = "ASSIGNED"
	}

	var returnedAt *time.Time

	if req.ReturnedAt != "" {
		t := utils.ParseDate(req.ReturnedAt)
		returnedAt = &t
	} else {
		returnedAt = nil
	}

	assignment := &models.AssetAssignment{
		AssetID:        req.AssetID,
		AssignedToID:   req.AssignedToID,
		AssignedAt:     utils.ParseDate(req.AssignedAt),
		ReturnedAt:     returnedAt,
		ConditionNotes: req.ConditionNotes,
		Status:         status,
	}

	if err := db.DB.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *AssetAssignmentService) GetAssignments(page, limit int) ([]dtos.AssetAssignmentResponse, int64, error) {
	var results []dtos.AssetAssignmentResponse
	var total int64
	db.DB.Model(&models.AssetAssignment{}).Count(&total)

	offset := (page - 1) * limit

	query := `
		SELECT 
			aa.id, aa.asset_id, fa.asset_name, fa.asset_code,
			aa.assigned_to_id, CONCAT(e.first_name, ' ', e.surname) AS assigned_to_name,
			aa.assigned_at, aa.returned_at, aa.condition_notes, aa.status,
			aa.created_at, aa.updated_at
		FROM asset_assignments aa
		LEFT JOIN fixed_assets fa ON aa.asset_id = fa.id
		LEFT JOIN employees e ON aa.assigned_to_id = e.id
		WHERE aa.deleted_at IS NULL
		ORDER BY aa.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *AssetAssignmentService) GetAssignment(id string) (*dtos.AssetAssignmentResponse, error) {
	var result dtos.AssetAssignmentResponse
	query := `
		SELECT 
			aa.id, aa.asset_id, fa.asset_name, fa.asset_code,
			aa.assigned_to_id, CONCAT(e.first_name, ' ', e.surname) AS assigned_to_name,
			aa.assigned_at, aa.returned_at, aa.condition_notes, aa.status,
			aa.created_at, aa.updated_at
		FROM asset_assignments aa
		LEFT JOIN fixed_assets fa ON aa.asset_id = fa.id
		LEFT JOIN employees e ON aa.assigned_to_id = e.id
		WHERE aa.id = ? AND aa.deleted_at IS NULL
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

func (s *AssetAssignmentService) UpdateAssignment(id string, req dtos.UpdateAssetAssignmentRequest) error {
	var assignment models.AssetAssignment
	if err := db.DB.First(&assignment, id).Error; err != nil {
		return err
	}

	assignment.AssetID = req.AssetID
	assignment.AssignedToID = req.AssignedToID
	assignment.AssignedAt = utils.ParseDate(req.AssignedAt)
	if req.ReturnedAt != "" {
		t := utils.ParseDate(req.ReturnedAt)
		assignment.ReturnedAt = &t
	}
	assignment.ConditionNotes = req.ConditionNotes
	assignment.Status = req.Status

	return db.DB.Save(&assignment).Error
}

func (s *AssetAssignmentService) DeleteAssignment(id string) error {
	return db.DB.Delete(&models.AssetAssignment{}, id).Error
}
