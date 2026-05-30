package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type LocationService struct{}

func NewLocationService() *LocationService {
	return &LocationService{}
}

func (s *LocationService) CreateLocation(req dtos.CreateLocationRequest, userID uint64) (*models.Location, error) {
	location := &models.Location{
		BaseModel:    models.BaseModel{CreatedBy: userID},
		Code:         req.Code,
		LocationName: req.LocationName,
	}

	if err := db.DB.Create(location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (s *LocationService) GetLocations(page, limit int) ([]models.Location, int64, error) {
	var results []models.Location
	var total int64
	db.DB.Model(&models.Location{}).Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.Location{}).
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Find(&results).Error

	return results, total, err
}

func (s *LocationService) GetLocation(id string) (*models.Location, error) {
	var result models.Location
	if err := db.DB.First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *LocationService) UpdateLocation(id string, req dtos.UpdateLocationRequest, userID uint64) error {
	var location models.Location
	if err := db.DB.First(&location, id).Error; err != nil {
		return err
	}

	location.Code = req.Code
	location.LocationName = req.LocationName
	location.UpdatedBy = userID

	return db.DB.Save(&location).Error
}

func (s *LocationService) DeleteLocation(id string) error {
	return db.DB.Delete(&models.Location{}, id).Error
}

func (s *LocationService) GetCounties() ([]models.County, error) {
	var results []models.County
	err := db.DB.Find(&results).Error
	return results, err
}

func (s *LocationService) CreateCounty(name, code string) (*models.County, error) {
	county := &models.County{Name: name, Code: code}
	err := db.DB.Create(county).Error
	return county, err
}

func (s *LocationService) GetSubCounties(countyID string) ([]models.SubCounty, error) {
	var results []models.SubCounty
	query := db.DB.Model(&models.SubCounty{})
	if countyID != "" {
		query = query.Where("county_id = ?", countyID)
	}
	err := query.Find(&results).Error
	return results, err
}

func (s *LocationService) CreateSubCounty(countyID uint64, name string) (*models.SubCounty, error) {
	subCounty := &models.SubCounty{CountyID: countyID, Name: name}
	err := db.DB.Create(subCounty).Error
	return subCounty, err
}

func (s *LocationService) GetWards(subCountyID string) ([]models.Ward, error) {
	var results []models.Ward
	query := db.DB.Model(&models.Ward{})
	if subCountyID != "" {
		query = query.Where("sub_county_id = ?", subCountyID)
	}
	err := query.Find(&results).Error
	return results, err
}

func (s *LocationService) CreateWard(subCountyID uint64, name string) (*models.Ward, error) {
	ward := &models.Ward{SubCountyID: subCountyID, Name: name}
	err := db.DB.Create(ward).Error
	return ward, err
}
