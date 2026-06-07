package services

import (
	"fmt"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type UserStoreAssignmentService struct{}

func NewUserStoreAssignmentService() *UserStoreAssignmentService {
	return &UserStoreAssignmentService{}
}

func (s *UserStoreAssignmentService) CreateUserStoreAssignment(req dtos.CreateUserStoreAssignmentRequest, createdBy uint64) (*models.UserStoreAssignment, error) {
	assignment := &models.UserStoreAssignment{
		BaseModel: models.BaseModel{CreatedBy: createdBy},
		UserID:    req.UserID,
		StoreID:   req.StoreID,
	}

	// Check for existing assignment to prevent duplicates
	var existingAssignment models.UserStoreAssignment
	err := db.DB.Where("user_id = ? AND store_id = ?", req.UserID, req.StoreID).First(&existingAssignment).Error
	if err == nil {
		return nil, fmt.Errorf("assignment for user %d to store %d already exists", req.UserID, req.StoreID)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err := db.DB.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *UserStoreAssignmentService) GetUserStoreAssignments(page, limit int, userID, storeID string) ([]dtos.UserStoreAssignmentResponse, int64, error) {
	var results []dtos.UserStoreAssignmentResponse
	var total int64

	queryBuilder := db.DB.Model(&models.UserStoreAssignment{})
	if userID != "" {
		queryBuilder = queryBuilder.Where("user_id = ?", userID)
	}
	if storeID != "" {
		queryBuilder = queryBuilder.Where("store_id = ?", storeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT usa.id, usa.user_id, u.name as user_name, usa.store_id, s.store as store_name, usa.created_at, usa.updated_at
		FROM user_store_assignments usa
		LEFT JOIN users u ON usa.user_id = u.id
		LEFT JOIN stores s ON usa.store_id = s.id
		WHERE usa.deleted_at IS NULL
	`
	var args []interface{}
	if userID != "" {
		query += " AND usa.user_id = ?"
		args = append(args, userID)
	}
	if storeID != "" {
		query += " AND usa.store_id = ?"
		args = append(args, storeID)
	}
	query += " ORDER BY usa.id DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := db.DB.Raw(query, args...).Scan(&results).Error
	return results, total, err
}

func (s *UserStoreAssignmentService) GetUserStoreAssignment(id string) (*dtos.UserStoreAssignmentResponse, error) {
	var result dtos.UserStoreAssignmentResponse
	query := `
		SELECT usa.id, usa.user_id, u.name as user_name, usa.store_id, s.store as store_name, usa.created_at, usa.updated_at
		FROM user_store_assignments usa
		LEFT JOIN users u ON usa.user_id = u.id
		LEFT JOIN stores s ON usa.store_id = s.id
		WHERE usa.id = ? AND usa.deleted_at IS NULL
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

func (s *UserStoreAssignmentService) UpdateUserStoreAssignment(id string, req dtos.UpdateUserStoreAssignmentRequest, updatedBy uint64) error {
	var assignment models.UserStoreAssignment
	if err := db.DB.First(&assignment, id).Error; err != nil {
		return err
	}

	// Check for existing assignment to prevent duplicates if user_id or store_id is changed
	var existingAssignment models.UserStoreAssignment
	err := db.DB.Where("user_id = ? AND store_id = ? AND id != ?", req.UserID, req.StoreID, id).First(&existingAssignment).Error
	if err == nil {
		return fmt.Errorf("assignment for user %d to store %d already exists", req.UserID, req.StoreID)
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	updates := map[string]interface{}{
		"user_id":    req.UserID,
		"store_id":   req.StoreID,
		"updated_by": updatedBy,
	}

	return db.DB.Model(&assignment).Updates(updates).Error
}

func (s *UserStoreAssignmentService) DeleteUserStoreAssignment(id string, deletedBy uint64) error {
	var assignment models.UserStoreAssignment
	if err := db.DB.First(&assignment, id).Error; err != nil {
		return err
	}
	// Soft delete and update deleted_by
	return db.DB.Model(&assignment).Update("deleted_by", deletedBy).Delete(&assignment).Error
}
