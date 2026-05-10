package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type CompanyTransporterService struct{}

func NewCompanyTransporterService() *CompanyTransporterService {
	return &CompanyTransporterService{}
}

func (s *CompanyTransporterService) GetCompanyTransporters() ([]dtos.CompanyTransporterResponse, int64, error) {
	var companies []models.CompanyTransporter
	var total int64
	db.DB.Model(&models.CompanyTransporter{}).Count(&total)
	err := db.DB.Preload("Transporter").Find(&companies).Error
	if err != nil {
		return nil, 0, err
	}

	var responses []dtos.CompanyTransporterResponse
	for _, c := range companies {
		responses = append(responses, s.toResponse(c))
	}
	return responses, total, nil
}

func (s *CompanyTransporterService) GetCompanyTransporter(id string) (*dtos.CompanyTransporterResponse, error) {
	var company models.CompanyTransporter
	if err := db.DB.Preload("Transporter").First(&company, id).Error; err != nil {
		return nil, err
	}
	response := s.toResponse(company)
	return &response, nil
}

func (s *CompanyTransporterService) UpdateCompanyTransporter(id string, req dtos.UpdateCompanyTransporterRequest) error {
	var company models.CompanyTransporter
	if err := db.DB.First(&company, id).Error; err != nil {
		return err
	}

	company.CompanyName = req.CompanyName
	company.RegistrationNo = req.RegistrationNo
	company.KraPin = req.KraPin
	company.ContactPersonName = req.ContactPersonName
	company.ContactPersonPhone = req.ContactPersonPhone
	company.PostalAddress = req.PostalAddress
	company.PostalCode = req.PostalCode
	company.Town = req.Town
	company.CertificateOfIncorporation = req.CertificateOfIncorporation

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&company).Error; err != nil {
			return err
		}
		return tx.Model(&models.Transporter{}).Where("id = ?", company.TransporterID).Update("route_id", req.RouteID).Error
	})
}

func (s *CompanyTransporterService) toResponse(c models.CompanyTransporter) dtos.CompanyTransporterResponse {
	routeID := uint64(0)
	if c.Transporter != nil {
		routeID = c.Transporter.RouteID
	}

	return dtos.CompanyTransporterResponse{
		ID:                         c.ID,
		TransporterID:              c.TransporterID,
		RouteID:                    routeID,
		CompanyName:                c.CompanyName,
		RegistrationNo:             c.RegistrationNo,
		KraPin:                     c.KraPin,
		ContactPersonName:          c.ContactPersonName,
		ContactPersonPhone:         c.ContactPersonPhone,
		PostalAddress:              c.PostalAddress,
		PostalCode:                 c.PostalCode,
		Town:                       c.Town,
		CertificateOfIncorporation: c.CertificateOfIncorporation,
	}
}
