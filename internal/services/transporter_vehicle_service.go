package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type TransporterVehicleService struct{}

func NewTransporterVehicleService() *TransporterVehicleService {
	return &TransporterVehicleService{}
}

func (s *TransporterVehicleService) CreateVehicle(req dtos.CreateTransporterVehicleRequest) (*models.TransporterVehicle, error) {
	vehicle := &models.TransporterVehicle{
		TransporterID:  req.TransporterID,
		RouteID:        req.RouteID,
		RegistrationNo: req.RegistrationNo,
		VehicleType:    req.VehicleType,
		CapacityLitres: req.CapacityLitres,
		Active:         req.Active,
	}

	if err := db.DB.Create(vehicle).Error; err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *TransporterVehicleService) GetVehicles(transporterID string) ([]dtos.TransporterVehicleResponse, int64, error) {
	var vehicles []models.TransporterVehicle
	var total int64

	query := db.DB.Model(&models.TransporterVehicle{})
	if transporterID != "" {
		query = query.Where("transporter_id = ?", transporterID)
	}

	query.Count(&total)
	err := query.Find(&vehicles).Error
	if err != nil {
		return nil, 0, err
	}
	var responses []dtos.TransporterVehicleResponse
	for _, v := range vehicles {
		responses = append(responses, s.toTransporterVehicleResponse(v))
	}
	return responses, total, nil
}

func (s *TransporterVehicleService) GetVehicle(id string) (*dtos.TransporterVehicleResponse, error) {
	var vehicle models.TransporterVehicle
	if err := db.DB.First(&vehicle, id).Error; err != nil {
		return nil, err
	}
	response := s.toTransporterVehicleResponse(vehicle)
	return &response, nil
}

func (s *TransporterVehicleService) toTransporterVehicleResponse(v models.TransporterVehicle) dtos.TransporterVehicleResponse {
	return dtos.TransporterVehicleResponse{
		ID:             v.ID,
		TransporterID:  v.TransporterID,
		RouteID:        v.RouteID,
		RegistrationNo: v.RegistrationNo,
		VehicleType:    v.VehicleType,
		CapacityLitres: v.CapacityLitres,
		Active:         v.Active,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
	}
}

func (s *TransporterVehicleService) UpdateVehicle(id string, req dtos.UpdateTransporterVehicleRequest) error {
	var vehicle models.TransporterVehicle
	if err := db.DB.First(&vehicle, id).Error; err != nil {
		return err
	}

	vehicle.RouteID = req.RouteID
	vehicle.VehicleType = req.VehicleType
	vehicle.CapacityLitres = req.CapacityLitres
	vehicle.Active = req.Active

	return db.DB.Save(&vehicle).Error
}

func (s *TransporterVehicleService) DeleteVehicle(id string) error {
	return db.DB.Delete(&models.TransporterVehicle{}, id).Error
}
