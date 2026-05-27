package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type LocationService struct{}

func NewLocationService() *LocationService {
	return &LocationService{}
}

// County Methods
func (s *LocationService) GetCounties() ([]models.County, error) {
	var counties []models.County
	err := db.DB.Find(&counties).Error
	return counties, err
}

func (s *LocationService) CreateCounty(name, code string) (*models.County, error) {
	county := &models.County{Name: name, Code: code}
	err := db.DB.Create(county).Error
	return county, err
}

// SubCounty Methods
func (s *LocationService) GetSubCounties(countyID string) ([]models.SubCounty, error) {
	var subCounties []models.SubCounty
	query := db.DB
	if countyID != "" {
		query = query.Where("county_id = ?", countyID)
	}
	err := query.Find(&subCounties).Error
	return subCounties, err
}

func (s *LocationService) CreateSubCounty(countyID uint64, name string) (*models.SubCounty, error) {
	subCounty := &models.SubCounty{CountyID: countyID, Name: name}
	err := db.DB.Create(subCounty).Error
	return subCounty, err
}

// Ward Methods
func (s *LocationService) GetWards(subCountyID string) ([]models.Ward, error) {
	var wards []models.Ward
	query := db.DB
	if subCountyID != "" {
		query = query.Where("sub_county_id = ?", subCountyID)
	}
	err := query.Find(&wards).Error
	return wards, err
}

func (s *LocationService) CreateWard(subCountyID uint64, name string) (*models.Ward, error) {
	ward := &models.Ward{SubCountyID: subCountyID, Name: name}
	err := db.DB.Create(ward).Error
	return ward, err
}
