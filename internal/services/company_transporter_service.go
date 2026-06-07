package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
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
	company.RegistrationNo = utils.StringPtr(req.RegistrationNo)
	company.KraPin = utils.StringPtr(req.KraPin)
	company.ContactPersonName = utils.StringPtr(req.ContactPersonName)
	company.ContactPersonPhone = utils.StringPtr(req.ContactPersonPhone)
	company.PostalAddress = utils.StringPtr(req.PostalAddress)
	company.PostalCode = utils.StringPtr(req.PostalCode)
	company.Town = utils.StringPtr(req.Town)
	company.CertificateOfIncorporation = utils.StringPtr(req.CertificateOfIncorporation)

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&company).Error; err != nil {
			return err
		}
		return tx.Model(&models.Transporter{}).Where("id = ?", company.TransporterID).Update("route_id", req.RouteID).Error
	})
}

func (s *CompanyTransporterService) toResponse(c models.CompanyTransporter) dtos.CompanyTransporterResponse {

	return dtos.CompanyTransporterResponse{
		ID:                         c.ID,
		TransporterID:              c.TransporterID,
		CompanyName:                c.CompanyName,
		RegistrationNo:             utils.StringValue(c.RegistrationNo),
		KraPin:                     utils.StringValue(c.KraPin),
		ContactPersonName:          utils.StringValue(c.ContactPersonName),
		ContactPersonPhone:         utils.StringValue(c.ContactPersonPhone),
		PostalAddress:              utils.StringValue(c.PostalAddress),
		PostalCode:                 utils.StringValue(c.PostalCode),
		Town:                       utils.StringValue(c.Town),
		CertificateOfIncorporation: utils.StringValue(c.CertificateOfIncorporation),
	}
}
