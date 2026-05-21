package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type ActivityLogService struct{}

func NewActivityLogService() *ActivityLogService {
	return &ActivityLogService{}
}

func (s *ActivityLogService) GetLogs(page, limit int) ([]dtos.ActivityLogResponse, int64, error) {
	var results []dtos.ActivityLogResponse
	var total int64

	db.DB.Model(&models.ActivityLog{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT al.*, u.name as causer_name
		FROM activity_log al
		LEFT JOIN users u ON al.causer_id = u.id
		WHERE al.deleted_at IS NULL
		ORDER BY al.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *ActivityLogService) GetLog(id string) (*dtos.ActivityLogResponse, error) {
	var result dtos.ActivityLogResponse

	query := `
		SELECT al.*, u.name as causer_name
		FROM activity_log al
		LEFT JOIN users u ON al.causer_id = u.id
		WHERE al.id = ? AND al.deleted_at IS NULL
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
