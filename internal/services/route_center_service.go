package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type RouteCenterService struct{}

func NewRouteCenterService() *RouteCenterService {
	return &RouteCenterService{}
}

func (s *RouteCenterService) CreateCenter(req dtos.CreateRouteCenterRequest) (*models.RouteCenter, error) {
	center := &models.RouteCenter{
		RouteID: req.RouteID,
		Center:  req.Center,
	}

	if err := db.DB.Create(center).Error; err != nil {
		return nil, err
	}
	return center, nil
}

func (s *RouteCenterService) GetCenters() ([]dtos.RouteCenterResponse, int64, error) {
	var results []dtos.RouteCenterResponse
	var total int64
	db.DB.Model(&models.RouteCenter{}).Count(&total)

	query := `
		SELECT 
			rc.id, rc.route_id, r.route_name as route_name, rc.center, 
			rc.created_at, rc.updated_at
		FROM route_centers rc
		LEFT JOIN routes r ON rc.route_id = r.id
		WHERE rc.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *RouteCenterService) GetCenter(id string) (*dtos.RouteCenterResponse, error) {
	var result dtos.RouteCenterResponse
	query := `
		SELECT 
			rc.id, rc.route_id, r.route_name as route_name, rc.center, 
			rc.created_at, rc.updated_at
		FROM route_centers rc
		LEFT JOIN routes r ON rc.route_id = r.id
		WHERE rc.id = ? AND rc.deleted_at IS NULL
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

func (s *RouteCenterService) UpdateCenter(id string, req dtos.UpdateRouteCenterRequest) error {
	var center models.RouteCenter
	if err := db.DB.First(&center, id).Error; err != nil {
		return err
	}

	center.RouteID = req.RouteID
	center.Center = req.Center

	return db.DB.Save(&center).Error
}

func (s *RouteCenterService) DeleteCenter(id string) error {
	return db.DB.Delete(&models.RouteCenter{}, id).Error
}
