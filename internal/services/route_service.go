package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type RouteService struct{}

func NewRouteService() *RouteService {
	return &RouteService{}
}

func (s *RouteService) CreateRoute(req dtos.CreateRouteRequest) (*models.Route, error) {
	route := &models.Route{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		LocationID:  req.LocationID,
	}

	if err := db.DB.Create(route).Error; err != nil {
		return nil, err
	}
	return route, nil
}

func (s *RouteService) GetRoutes() ([]dtos.RouteResponse, int64, error) {
	var results []dtos.RouteResponse
	var total int64
	db.DB.Model(&models.Route{}).Count(&total)

	query := `
		SELECT 
			r.id, r.route_name, r.description, r.route_code, r.location_id,
			l.location_name,
			r.created_at, r.updated_at
		FROM routes r
		LEFT JOIN locations l ON r.location_id = l.id
		WHERE r.deleted_at IS NULL
	`
	err := db.DB.Raw(query).Scan(&results).Error
	return results, total, err
}

func (s *RouteService) GetRoute(id string) (*dtos.RouteResponse, error) {
	var result dtos.RouteResponse
	query := `
		SELECT 
			r.id, r.route_name, r.description, r.route_code, r.location_id,
			l.location_name,
			r.created_at, r.updated_at
		FROM routes r
		LEFT JOIN locations l ON r.location_id = l.id
		WHERE r.id = ? AND r.deleted_at IS NULL
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

func (s *RouteService) UpdateRoute(id string, req dtos.UpdateRouteRequest) error {
	var route models.Route
	if err := db.DB.First(&route, id).Error; err != nil {
		return err
	}

	route.Name = req.Name
	route.Description = req.Description
	route.Code = req.Code
	route.LocationID = req.LocationID

	return db.DB.Save(&route).Error
}

func (s *RouteService) DeleteRoute(id string) error {
	return db.DB.Delete(&models.Route{}, id).Error
}
