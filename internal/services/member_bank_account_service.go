package services

import (
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type MemberBankAccountService struct{}

func NewMemberBankAccountService() *MemberBankAccountService {
	return &MemberBankAccountService{}
}

func (s *MemberBankAccountService) CreateAccount(req dtos.CreateMemberBankAccountRequest) (*models.MemberBankAccount, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	account := &models.MemberBankAccount{
		MemberID:      req.MemberID,
		BankID:        req.BankID,
		BankBranchId:  req.BankBranchId,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		Status:        status,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *MemberBankAccountService) GetAccounts(page, limit int, memberNo, firstName, lastName, bankName, accountNo string) ([]dtos.MemberBankAccountResponse, int64, error) {
	var accounts []dtos.MemberBankAccountResponse
	var total int64

	offset := (page - 1) * limit

	baseQuery := `
		SELECT 
			mba.id, mba.member_id, m.member_no, m.first_name, m.last_name,
			mba.bank_id, b.bank_name, mba.bank_branch_id,
			mba.account_number, mba.account_name, mba.status,
			mba.created_at, mba.updated_at
		FROM member_bank_accounts mba
		LEFT JOIN member_registrations m ON mba.member_id = m.id
		LEFT JOIN banks b ON mba.bank_id = b.id
	`

	baseCountQuery := `
		SELECT COUNT(*)
		FROM member_bank_accounts mba
		LEFT JOIN member_registrations m ON mba.member_id = m.id
		LEFT JOIN banks b ON mba.bank_id = b.id
	`

	var args []interface{}
	whereClauses := []string{"mba.deleted_at IS NULL"}

	if memberNo != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+memberNo+"%")
	}
	if firstName != "" {
		whereClauses = append(whereClauses, "m.first_name LIKE ?")
		args = append(args, "%"+firstName+"%")
	}
	if lastName != "" {
		whereClauses = append(whereClauses, "m.last_name LIKE ?")
		args = append(args, "%"+lastName+"%")
	}
	if bankName != "" {
		whereClauses = append(whereClauses, "b.bank_name LIKE ?")
		args = append(args, "%"+bankName+"%")
	}
	if accountNo != "" {
		whereClauses = append(whereClauses, "mba.account_number LIKE ?")
		args = append(args, "%"+accountNo+"%")
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	if err := db.DB.Raw(baseCountQuery+whereSql, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	fullQuery := baseQuery + whereSql + " ORDER BY mba.id DESC LIMIT ? OFFSET ?"
	if err := db.DB.Raw(fullQuery, append(args, limit, offset)...).Scan(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (s *MemberBankAccountService) GetAccount(id string) (*dtos.MemberBankAccountResponse, error) {
	var account dtos.MemberBankAccountResponse
	query := `
		SELECT 
			mba.id, mba.member_id, m.member_no, m.first_name, m.last_name,
			mba.bank_id, b.bank_name, mba.bank_branch_id,
			mba.account_number, mba.account_name, mba.status,
			mba.created_at, mba.updated_at
		FROM member_bank_accounts mba
		LEFT JOIN member_registrations m ON mba.member_id = m.id
		LEFT JOIN banks b ON mba.bank_id = b.id
		WHERE mba.id = ? AND mba.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&account).Error
	if err != nil {
		return nil, err
	}
	if account.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &account, nil
}

func (s *MemberBankAccountService) UpdateAccount(id string, req dtos.UpdateMemberBankAccountRequest) error {
	var account models.MemberBankAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	account.MemberID = req.MemberID
	account.BankID = req.BankID
	account.BankBranchId = req.BankBranchId
	account.AccountNumber = req.AccountNumber
	account.AccountName = req.AccountName
	account.Status = req.Status

	return db.DB.Save(&account).Error
}

func (s *MemberBankAccountService) DeleteAccount(id string) error {
	return db.DB.Delete(&models.MemberBankAccount{}, id).Error
}
