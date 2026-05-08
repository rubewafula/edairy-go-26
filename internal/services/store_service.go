package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type StoreService struct{}

func NewStoreService() *StoreService {
	return &StoreService{}
}

func (s *StoreService) CreateStore(req dtos.CreateStoreRequest) (*models.Store, error) {
	store := &models.Store{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.DB.Create(store).Error; err != nil {
		return nil, err
	}
	return store, nil
}

func (s *StoreService) GetStores() ([]models.Store, int64, error) {
	var stores []models.Store
	var total int64
	db.DB.Model(&models.Store{}).Count(&total)
	err := db.DB.Find(&stores).Error
	return stores, total, err
}

func (s *StoreService) GetStore(id string) (*models.Store, error) {
	var store models.Store
	if err := db.DB.First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *StoreService) UpdateStore(id string, req dtos.UpdateStoreRequest) error {
	var store models.Store
	if err := db.DB.First(&store, id).Error; err != nil {
		return err
	}

	store.Name = req.Name
	store.Description = req.Description

	return db.DB.Save(&store).Error
}

func (s *StoreService) DeleteStore(id string) error {
	return db.DB.Delete(&models.Store{}, id).Error
}
