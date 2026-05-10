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
	individual.OtherNames = req.OtherNames
	individual.Gender = req.Gender
	individual.DateOfBirth = utils.ParseDate(req.DateOfBirth)
	individual.NationalIDNo = req.NationalIDNo
	individual.KraPin = req.KraPin
	individual.MaritalStatus = req.MaritalStatus
	individual.NextOfKinFullName = req.NextOfKinFullName
	individual.NextOfKinPhone = req.NextOfKinPhone
	individual.PassportPhoto = req.PassportPhoto
	individual.IDFrontPhoto = req.IDFrontPhoto
	individual.IDBackPhoto = req.IDBackPhoto

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&individual).Error; err != nil {
			return err
		}
		// Update the RouteID on the parent Transporter record
		return tx.Model(&models.Transporter{}).Where("id = ?", individual.TransporterID).Update("route_id", req.RouteID).Error
	})
}

func (s *IndividualTransporterService) toResponse(i models.IndividualTransporter) dtos.IndividualTransporterResponse {
	routeID := uint64(0)
	if i.Transporter != nil {
		routeID = i.Transporter.RouteID
	}

	return dtos.IndividualTransporterResponse{
		ID:                i.ID,
		TransporterID:     i.TransporterID,
		RouteID:           routeID,
		FirstName:         i.FirstName,
		LastName:          i.LastName,
		OtherNames:        i.OtherNames,
		Gender:            i.Gender,
		DateOfBirth:       i.DateOfBirth,
		NationalIDNo:      i.NationalIDNo,
		KraPin:            i.KraPin,
		MaritalStatus:     i.MaritalStatus,
		NextOfKinFullName: i.NextOfKinFullName,
		NextOfKinPhone:    i.NextOfKinPhone,
		PassportPhoto:     i.PassportPhoto,
		IDFrontPhoto:      i.IDFrontPhoto,
		IDBackPhoto:       i.IDBackPhoto,
	}
}
