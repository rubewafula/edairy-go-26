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

// BankBranchListItem is one row of GET /bank-branches (mobile common data + admin list).
type BankBranchListItem struct {
	ID         uint64 `json:"id"`
	BankID     uint64 `json:"bank_id"`
	BankName   string `json:"bank_name"`
	BankCode   string `json:"bank_code"`
	Name       string `json:"name"`
	BranchName string `json:"branch_name"`
	Location   string `json:"location"`
}

// GetBankBranches lists branches with their bank's name; bankID > 0 limits it to one bank.
func (s *BankBranchService) GetBankBranches(bankID ...uint64) ([]BankBranchListItem, int64, error) {
	q := db.DB.Table("bank_branches bb").
		Select("bb.id, bb.bank_id, COALESCE(b.bank_name, '') AS bank_name, COALESCE(b.bank_code, '') AS bank_code, " +
			"COALESCE(bb.branch_name, '') AS name, COALESCE(bb.branch_name, '') AS branch_name, COALESCE(bb.location, '') AS location").
		Joins("LEFT JOIN banks b ON b.id = bb.bank_id AND b.deleted_at IS NULL").
		Where("bb.deleted_at IS NULL")
	if len(bankID) > 0 && bankID[0] > 0 {
		q = q.Where("bb.bank_id = ?", bankID[0])
	}
	var branches []BankBranchListItem
	err := q.Order("b.bank_name, bb.branch_name").Scan(&branches).Error
	if branches == nil {
		branches = []BankBranchListItem{}
	}
	return branches, int64(len(branches)), err
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
