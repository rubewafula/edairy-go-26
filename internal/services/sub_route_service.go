package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type SubRouteService struct{}

func NewSubRouteService() *SubRouteService {
	return &SubRouteService{}
}

func (s *SubRouteService) CreateSubRoute(req dtos.CreateSubRouteRequest) (*models.SubRoute, error) {
	subRoute := &models.SubRoute{
		RouteID:     req.RouteID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(subRoute).Error; err != nil {
		return nil, err
	}
	return subRoute, nil
}

func (s *SubRouteService) GetSubRoutes() ([]models.SubRoute, int64, error) {
	var subRoutes []models.SubRoute
	var total int64
	db.DB.Model(&models.SubRoute{}).Count(&total)
	err := db.DB.Find(&subRoutes).Error
	return subRoutes, total, err
}

func (s *SubRouteService) GetSubRoute(id string) (*models.SubRoute, error) {
	var subRoute models.SubRoute
	if err := db.DB.First(&subRoute, id).Error; err != nil {
		return nil, err
	}
	return &subRoute, nil
}

func (s *SubRouteService) UpdateSubRoute(id string, req dtos.UpdateSubRouteRequest) error {
	var subRoute models.SubRoute
	if err := db.DB.First(&subRoute, id).Error; err != nil {
		return err
	}

	subRoute.RouteID = req.RouteID
	subRoute.Name = req.Name
	subRoute.Description = req.Description

	return db.DB.Save(&subRoute).Error
}

func (s *SubRouteService) DeleteSubRoute(id string) error {
	return db.DB.Delete(&models.SubRoute{}, id).Error
}
