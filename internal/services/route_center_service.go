package services

import (
	"fmt"
	"strings"

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

// GetCenters lists route centers, optionally only those belonging to routeID.
func (s *RouteCenterService) GetCenters(routeID string) ([]dtos.RouteCenterResponse, int64, error) {
	var results []dtos.RouteCenterResponse
	var total int64

	countQuery := db.DB.Model(&models.RouteCenter{})
	query := `
		SELECT 
			rc.id, rc.route_id, r.route_name as route_name, rc.center, 
			rc.created_at, rc.updated_at
		FROM route_centers rc
		LEFT JOIN routes r ON rc.route_id = r.id
		WHERE rc.deleted_at IS NULL
	`
	var args []interface{}
	if routeID != "" {
		countQuery = countQuery.Where("route_id = ?", routeID)
		query += " AND rc.route_id = ?"
		args = append(args, routeID)
	}
	query += " ORDER BY rc.center ASC"

	countQuery.Count(&total)
	err := db.DB.Raw(query, args...).Scan(&results).Error
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

// defaultCenterName builds the name of the default center seeded for a route.
func defaultCenterName(routeID uint64, routeName, routeCode string) string {
	base := strings.TrimSpace(routeName)
	if base == "" {
		base = strings.TrimSpace(routeCode)
	}
	if base == "" {
		base = fmt.Sprintf("Route %d", routeID)
	}
	name := base + " Center"
	if r := []rune(name); len(r) > 96 { // route_centers.center is varchar(96)
		name = string(r[:96])
	}
	return name
}

// EnsureRouteCenters guarantees every route has at least one center: milk is always
// recorded against a center. Routes without one get a default "<route name> Center".
// Safe to run repeatedly; returns the number of centers created.
func EnsureRouteCenters() (int64, error) {
	type routeRow struct {
		ID        uint64
		RouteName string
		RouteCode string
	}
	var missing []routeRow
	err := db.DB.Raw(`
		SELECT r.id, COALESCE(r.route_name, '') AS route_name, COALESCE(r.route_code, '') AS route_code
		FROM routes r
		WHERE r.deleted_at IS NULL
		  AND NOT EXISTS (
			SELECT 1 FROM route_centers rc
			WHERE rc.route_id = r.id AND rc.deleted_at IS NULL
		  )`).Scan(&missing).Error
	if err != nil {
		return 0, err
	}
	if len(missing) == 0 {
		return 0, nil
	}
	centers := make([]models.RouteCenter, 0, len(missing))
	for _, r := range missing {
		centers = append(centers, models.RouteCenter{
			RouteID: r.ID,
			Center:  defaultCenterName(r.ID, r.RouteName, r.RouteCode),
		})
	}
	if err := db.DB.CreateInBatches(&centers, 200).Error; err != nil {
		return 0, err
	}
	return int64(len(centers)), nil
}
