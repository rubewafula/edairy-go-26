package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type CattleBreedService struct{}

func NewCattleBreedService() *CattleBreedService {
	return &CattleBreedService{}
}

func (s *CattleBreedService) CreateCattleBreed(req dtos.CreateCattleBreedRequest) (*models.CattleBreed, error) {
	cattleBreed := &models.CattleBreed{
		Code: req.Code,
		Name: req.Name,
	}

	if err := db.DB.Create(cattleBreed).Error; err != nil {
		return nil, err
	}
	return cattleBreed, nil
}

func (s *CattleBreedService) GetCattleBreeds() ([]models.CattleBreed, int64, error) {
	var cattleBreeds []models.CattleBreed
	var total int64
	db.DB.Model(&models.CattleBreed{}).Count(&total)
	err := db.DB.Find(&cattleBreeds).Error
	return cattleBreeds, total, err
}

func (s *CattleBreedService) GetCattleBreed(id string) (*models.CattleBreed, error) {
	var cattleBreed models.CattleBreed
	if err := db.DB.First(&cattleBreed, id).Error; err != nil {
		return nil, err
	}
	return &cattleBreed, nil
}

func (s *CattleBreedService) UpdateCattleBreed(id string, req dtos.UpdateCattleBreedRequest) error {
	var cattleBreed models.CattleBreed
	if err := db.DB.First(&cattleBreed, id).Error; err != nil {
		return err
	}

	cattleBreed.Code = req.Code
	cattleBreed.Name = req.Name

	return db.DB.Save(&cattleBreed).Error
}

func (s *CattleBreedService) DeleteCattleBreed(id string) error {
	return db.DB.Delete(&models.CattleBreed{}, id).Error
}
