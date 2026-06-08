package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type IndividualTransporterService struct{}

func NewIndividualTransporterService() *IndividualTransporterService {
	return &IndividualTransporterService{}
}

func (s *IndividualTransporterService) GetIndividualTransporters() ([]dtos.IndividualTransporterResponse, int64, error) {
	var individuals []models.IndividualTransporter
	var total int64
	db.DB.Model(&models.IndividualTransporter{}).Count(&total)
	err := db.DB.Preload("Transporter").Find(&individuals).Error
	if err != nil {
		return nil, 0, err
	}

	var responses []dtos.IndividualTransporterResponse
	for _, i := range individuals {
		responses = append(responses, s.toResponse(i))
	}
	return responses, total, nil
}

func (s *IndividualTransporterService) GetIndividualTransporter(id string) (*dtos.IndividualTransporterResponse, error) {
	var individual models.IndividualTransporter
	if err := db.DB.Preload("Transporter").First(&individual, id).Error; err != nil {
		return nil, err
	}
	response := s.toResponse(individual)
	return &response, nil
}

func (s *IndividualTransporterService) UpdateIndividualTransporter(id string, req dtos.UpdateIndividualTransporterRequest) error {
	var individual models.IndividualTransporter
	if err := db.DB.First(&individual, id).Error; err != nil {
		return err
	}

	individual.FirstName = req.FirstName
	individual.LastName = req.LastName
	individual.OtherNames = utils.StringPtr(req.OtherNames)
	individual.Gender = utils.StringPtr(req.Gender)

	if req.DateOfBirth != "" {
		t := utils.ParseDate(req.DateOfBirth)
		if !t.IsZero() {
			individual.DateOfBirth = &t
		}
	}

	individual.NationalIDNo = utils.StringPtr(req.NationalIDNo)
	individual.KraPin = utils.StringPtr(req.KraPin)
	individual.MaritalStatus = utils.StringPtr(req.MaritalStatus)
	individual.NextOfKinFullName = utils.StringPtr(req.NextOfKinFullName)
	individual.NextOfKinPhone = utils.StringPtr(req.NextOfKinPhone)
	individual.PassportPhoto = utils.StringPtr(req.PassportPhoto)
	individual.IDFrontPhoto = utils.StringPtr(req.IDFrontPhoto)
	individual.IDBackPhoto = utils.StringPtr(req.IDBackPhoto)

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&individual).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *IndividualTransporterService) toResponse(i models.IndividualTransporter) dtos.IndividualTransporterResponse {

	return dtos.IndividualTransporterResponse{
		ID:                i.ID,
		TransporterID:     i.TransporterID,
		FirstName:         i.FirstName,
		LastName:          i.LastName,
		OtherNames:        utils.StringValue(i.OtherNames),
		Gender:            utils.StringValue(i.Gender),
		DateOfBirth:       i.DateOfBirth,
		NationalIDNo:      utils.StringValue(i.NationalIDNo),
		KraPin:            utils.StringValue(i.KraPin),
		MaritalStatus:     utils.StringValue(i.MaritalStatus),
		NextOfKinFullName: utils.StringValue(i.NextOfKinFullName),
		NextOfKinPhone:    utils.StringValue(i.NextOfKinPhone),
		PassportPhoto:     utils.StringValue(i.PassportPhoto),
		IDFrontPhoto:      utils.StringValue(i.IDFrontPhoto),
		IDBackPhoto:       utils.StringValue(i.IDBackPhoto),
	}
}
