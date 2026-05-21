package services

import (
	"log"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type TransporterService struct{}

func NewTransporterService() *TransporterService {
	return &TransporterService{}
}

func (s *TransporterService) CreateTransporter(req dtos.CreateTransporterRequest) (*dtos.TransporterResponse, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	transporter := &models.Transporter{
		TransporterNo: req.TransporterNo,
		Category:      req.Category,
		PrimaryPhone:  req.PrimaryPhone,
		EmailAddress:  req.EmailAddress,
		Status:        status,
		Restricted:    req.Restricted,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transporter).Error; err != nil {
			return err
		}

		if req.Category == "INDIVIDUAL" {
			passportPath, _ := utils.SaveFile(req.PassportPhoto, "transporters")
			idFrontPath, _ := utils.SaveFile(req.IDFrontPhoto, "transporters")
			idBackPath, _ := utils.SaveFile(req.IDBackPhoto, "transporters")

			individual := &models.IndividualTransporter{
				TransporterID:     transporter.ID,
				FirstName:         req.FirstName,
				LastName:          req.LastName,
				OtherNames:        req.OtherNames,
				Gender:            req.Gender,
				NationalIDNo:      req.NationalIDNo,
				KraPin:            req.KraPin,
				MaritalStatus:     req.MaritalStatus,
				NextOfKinFullName: req.NextOfKinFullName,
				NextOfKinPhone:    req.NextOfKinPhone,
				PassportPhoto:     passportPath,
				IDFrontPhoto:      idFrontPath,
				IDBackPhoto:       idBackPath,
			}

			if req.DateOfBirth != "" {
				t := utils.ParseDate(req.DateOfBirth)
				individual.DateOfBirth = &t
			}

			return tx.Create(individual).Error
		} else {
			certificatePath := ""
			if req.CertificateOfIncorporation != nil {
				certificatePath, _ = utils.SaveFile(req.CertificateOfIncorporation, "transporters")
			}

			company := &models.CompanyTransporter{
				TransporterID:              transporter.ID,
				CompanyName:                req.CompanyName,
				RegistrationNo:             req.RegistrationNo,
				KraPin:                     req.KraPin,
				ContactPersonName:          req.ContactPersonName,
				ContactPersonPhone:         req.ContactPersonPhone,
				PostalAddress:              req.PostalAddress,
				PostalCode:                 req.PostalCode,
				Town:                       req.Town,
				CertificateOfIncorporation: certificatePath,
			}
			return tx.Create(company).Error
		}
	})

	if err != nil {
		log.Printf("Transporter: Error creating: %s", err.Error())
		return nil, err
	}

	var created models.Transporter
	if err := db.DB.Preload("Individual").Preload("Company").First(&created, transporter.ID).Error; err != nil {
		return nil, err
	}
	res := s.toTransporterResponse(created)
	return &res, nil
}

func (s *TransporterService) GetTransporters() ([]dtos.TransporterResponse, int64, error) {
	var transporters []models.Transporter
	var total int64
	db.DB.Model(&models.Transporter{}).Count(&total)
	err := db.DB.Preload("Individual").Preload("Company").Find(&transporters).Error
	if err != nil {
		return nil, 0, err
	}

	var responses []dtos.TransporterResponse
	for _, t := range transporters {
		responses = append(responses, s.toTransporterResponse(t))
	}
	return responses, total, nil
}

func (s *TransporterService) GetTransporter(id string) (*dtos.TransporterResponse, error) {
	var transporter models.Transporter
	if err := db.DB.Preload("Individual").Preload("Company").First(&transporter, id).Error; err != nil {
		return nil, err
	}
	response := s.toTransporterResponse(transporter)
	return &response, nil
}

func (s *TransporterService) toTransporterResponse(t models.Transporter) dtos.TransporterResponse {
	response := dtos.TransporterResponse{
		ID:            t.ID,
		TransporterNo: t.TransporterNo,
		Category:      t.Category,
		PrimaryPhone:  t.PrimaryPhone,
		EmailAddress:  t.EmailAddress,
		Status:        t.Status,
		Restricted:    t.Restricted,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	if t.Individual != nil {
		response.Individual = &dtos.IndividualTransporterResponse{
			ID: t.Individual.ID, TransporterID: t.Individual.TransporterID, FirstName: t.Individual.FirstName, LastName: t.Individual.LastName, OtherNames: t.Individual.OtherNames, Gender: t.Individual.Gender, DateOfBirth: t.Individual.DateOfBirth, NationalIDNo: t.Individual.NationalIDNo, KraPin: t.Individual.KraPin, MaritalStatus: t.Individual.MaritalStatus, NextOfKinFullName: t.Individual.NextOfKinFullName, NextOfKinPhone: t.Individual.NextOfKinPhone, PassportPhoto: t.Individual.PassportPhoto, IDFrontPhoto: t.Individual.IDFrontPhoto, IDBackPhoto: t.Individual.IDBackPhoto,
		}
	}
	if t.Company != nil {
		response.Company = &dtos.CompanyTransporterResponse{
			ID: t.Company.ID, TransporterID: t.Company.TransporterID, CompanyName: t.Company.CompanyName, RegistrationNo: t.Company.RegistrationNo, KraPin: t.Company.KraPin, ContactPersonName: t.Company.ContactPersonName, ContactPersonPhone: t.Company.ContactPersonPhone, PostalAddress: t.Company.PostalAddress, PostalCode: t.Company.PostalCode, Town: t.Company.Town, CertificateOfIncorporation: t.Company.CertificateOfIncorporation,
		}
	}
	return response
}

func (s *TransporterService) UpdateTransporter(id string, req dtos.UpdateTransporterRequest) error {
	var transporter models.Transporter
	if err := db.DB.First(&transporter, id).Error; err != nil {
		return err
	}

	transporter.PrimaryPhone = req.PrimaryPhone
	transporter.EmailAddress = req.EmailAddress
	transporter.Status = req.Status
	transporter.Restricted = req.Restricted

	return db.DB.Save(&transporter).Error
}

func (s *TransporterService) DeleteTransporter(id string) error {
	return db.DB.Delete(&models.Transporter{}, id).Error
}
