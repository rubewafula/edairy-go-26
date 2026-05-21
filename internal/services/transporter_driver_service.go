package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type TransporterDriverService struct{}

func NewTransporterDriverService() *TransporterDriverService {
	return &TransporterDriverService{}
}

func (s *TransporterDriverService) CreateDriver(req dtos.CreateTransporterDriverRequest) (*models.TransporterDriver, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	driver := &models.TransporterDriver{
		TransporterID:            req.TransporterID,
		DriverNo:                 req.DriverNo,
		FirstName:                req.FirstName,
		LastName:                 req.LastName,
		OtherNames:               req.OtherNames,
		Gender:                   req.Gender,
		NationalIDNo:             req.NationalIDNo,
		KraPin:                   req.KraPin,
		PrimaryPhone:             req.PrimaryPhone,
		SecondaryPhone:           req.SecondaryPhone,
		EmailAddress:             req.EmailAddress,
		DrivingLicenseNo:         req.DrivingLicenseNo,
		DrivingLicenseClass:      req.DrivingLicenseClass,
		Status:                   status,
		NextOfKinFullName:        req.NextOfKinFullName,
		NextOfKinPhone:           req.NextOfKinPhone,
		PassportPhoto:            req.PassportPhoto,
		DrivingLicenseFrontPhoto: req.DrivingLicenseFrontPhoto,
		DrivingLicenseBackPhoto:  req.DrivingLicenseBackPhoto,
	}

	if req.DateOfBirth != "" {
		t := utils.ParseDate(req.DateOfBirth)
		driver.DateOfBirth = &t
	}
	if req.DrivingLicenseExpiry != "" {
		t := utils.ParseDate(req.DrivingLicenseExpiry)
		driver.DrivingLicenseExpiry = &t
	}
	if req.EmploymentDate != "" {
		t := utils.ParseDate(req.EmploymentDate)
		driver.EmploymentDate = &t
	}

	if err := db.DB.Create(driver).Error; err != nil {
		return nil, err
	}
	return driver, nil
}

func (s *TransporterDriverService) GetDrivers() ([]dtos.TransporterDriverResponse, int64, error) {
	var drivers []models.TransporterDriver
	var total int64
	db.DB.Model(&models.TransporterDriver{}).Count(&total)
	err := db.DB.Find(&drivers).Error
	if err != nil {
		return nil, 0, err
	}
	var responses []dtos.TransporterDriverResponse
	for _, d := range drivers {
		responses = append(responses, s.toTransporterDriverResponse(d))
	}
	return responses, total, nil
}

func (s *TransporterDriverService) GetDriver(id string) (*dtos.TransporterDriverResponse, error) {
	var driver models.TransporterDriver
	if err := db.DB.First(&driver, id).Error; err != nil {
		return nil, err
	}
	response := s.toTransporterDriverResponse(driver)
	return &response, nil
}

func (s *TransporterDriverService) toTransporterDriverResponse(d models.TransporterDriver) dtos.TransporterDriverResponse {
	return dtos.TransporterDriverResponse{
		ID: d.ID, TransporterID: d.TransporterID, DriverNo: d.DriverNo, FirstName: d.FirstName, LastName: d.LastName, OtherNames: d.OtherNames, Gender: d.Gender, DateOfBirth: d.DateOfBirth, NationalIDNo: d.NationalIDNo, KraPin: d.KraPin, PrimaryPhone: d.PrimaryPhone, SecondaryPhone: d.SecondaryPhone, EmailAddress: d.EmailAddress, DrivingLicenseNo: d.DrivingLicenseNo, DrivingLicenseClass: d.DrivingLicenseClass, DrivingLicenseExpiry: d.DrivingLicenseExpiry, EmploymentDate: d.EmploymentDate, Status: d.Status, PassportPhoto: d.PassportPhoto, DrivingLicenseFrontPhoto: d.DrivingLicenseFrontPhoto, DrivingLicenseBackPhoto: d.DrivingLicenseBackPhoto, NextOfKinFullName: d.NextOfKinFullName, NextOfKinPhone: d.NextOfKinPhone, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func (s *TransporterDriverService) UpdateDriver(id string, req dtos.UpdateTransporterDriverRequest) error {
	var driver models.TransporterDriver
	// The UpdateTransporterDriverRequest DTO is missing photo fields.
	// If these fields are meant to be updatable, they should be added to the DTO.
	// For now, I'll update only the fields present in the DTO.
	// Also, NationalIDNo and DrivingLicenseNo are unique and usually not updated.
	// Assuming they are not intended to be updated via this endpoint.
	if err := db.DB.First(&driver, id).Error; err != nil {
		return err
	}

	driver.FirstName = req.FirstName
	driver.LastName = req.LastName
	driver.OtherNames = req.OtherNames
	driver.Gender = req.Gender
	if req.DateOfBirth != "" {
		t := utils.ParseDate(req.DateOfBirth)
		driver.DateOfBirth = &t
	}
	driver.PrimaryPhone = req.PrimaryPhone
	if req.DrivingLicenseExpiry != "" {
		t := utils.ParseDate(req.DrivingLicenseExpiry)
		driver.DrivingLicenseExpiry = &t
	}
	driver.Status = req.Status
	driver.NextOfKinFullName = req.NextOfKinFullName
	driver.NextOfKinPhone = req.NextOfKinPhone
	driver.PassportPhoto = req.PassportPhoto
	driver.DrivingLicenseFrontPhoto = req.DrivingLicenseFrontPhoto
	driver.DrivingLicenseBackPhoto = req.DrivingLicenseBackPhoto

	return db.DB.Save(&driver).Error
}

func (s *TransporterDriverService) DeleteDriver(id string) error {
	return db.DB.Delete(&models.TransporterDriver{}, id).Error
}
