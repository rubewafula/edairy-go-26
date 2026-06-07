package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
)

type BankBranchService struct{}

func NewBankBranchService() *BankBranchService {
	return &BankBranchService{}
}

func (s *BankBranchService) CreateBankBranch(req dtos.CreateBankBranchRequest) (*models.BankBranch, error) {
	branch := &models.BankBranch{
		Name:     req.Name,
		BankID:   req.BankID,
		Location: req.Location,
	}

	if err := db.DB.Create(branch).Error; err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *BankBranchService) GetBankBranches() ([]models.BankBranch, int64, error) {
	var branches []models.BankBranch
	var total int64
	db.DB.Model(&models.BankBranch{}).Count(&total)
	err := db.DB.Find(&branches).Error
	return branches, total, err
}

func (s *BankBranchService) GetBankBranch(id string) (*models.BankBranch, error) {
	var branch models.BankBranch
	if err := db.DB.First(&branch, id).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

func (s *BankBranchService) UpdateBankBranch(id string, req dtos.UpdateBankBranchRequest) error {
	var branch models.BankBranch
	if err := db.DB.First(&branch, id).Error; err != nil {
		return err
	}

	branch.Name = req.Name
	branch.BankID = req.BankID
	branch.Location = req.Location

	return db.DB.Save(&branch).Error
}

func (s *BankBranchService) DeleteBankBranch(id string) error {
	return db.DB.Delete(&models.BankBranch{}, id).Error
}
