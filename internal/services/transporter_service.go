package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type TransporterService struct{}

func NewTransporterService() *TransporterService {
	return &TransporterService{}
}

func (s *TransporterService) CreateTransporter(req dtos.CreateTransporterRequest) (*models.Transporter, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	transporter := &models.Transporter{
		Name:         req.Name,
		Phone:        req.Phone,
		IDNumber:     req.IDNumber,
		VehicleRegNo: req.VehicleRegNo,
		Status:       status,
	}

	if err := db.DB.Create(transporter).Error; err != nil {
		return nil, err
	}
	return transporter, nil
}

func (s *TransporterService) GetTransporters() ([]models.Transporter, int64, error) {
	var transporters []models.Transporter
	var total int64
	db.DB.Model(&models.Transporter{}).Count(&total)
	err := db.DB.Find(&transporters).Error
	return transporters, total, err
}

func (s *TransporterService) GetTransporter(id string) (*models.Transporter, error) {
	var transporter models.Transporter
	if err := db.DB.First(&transporter, id).Error; err != nil {
		return nil, err
	}
	return &transporter, nil
}

func (s *TransporterService) UpdateTransporter(id string, req dtos.UpdateTransporterRequest) error {
	var transporter models.Transporter
	if err := db.DB.First(&transporter, id).Error; err != nil {
		return err
	}

	transporter.Name = req.Name
	transporter.Phone = req.Phone
	transporter.IDNumber = req.IDNumber
	transporter.VehicleRegNo = req.VehicleRegNo
	transporter.Status = req.Status

	return db.DB.Save(&transporter).Error
}

func (s *TransporterService) DeleteTransporter(id string) error {
	return db.DB.Delete(&models.Transporter{}, id).Error
}
