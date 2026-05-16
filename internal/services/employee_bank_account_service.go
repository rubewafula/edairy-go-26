package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type EmployeeBankAccountService struct{}

func NewEmployeeBankAccountService() *EmployeeBankAccountService {
	return &EmployeeBankAccountService{}
}

func (s *EmployeeBankAccountService) CreateAccount(req dtos.CreateEmployeeBankAccountRequest, userID uint64) (*models.EmployeeBankAccount, error) {
	account := &models.EmployeeBankAccount{
		BaseModel:     models.BaseModel{CreatedBy: userID},
		EmployeeID:    req.EmployeeID,
		BankID:        req.BankID,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *EmployeeBankAccountService) GetAccounts(employeeID string, page, limit int) ([]dtos.EmployeeBankAccountResponse, int64, error) {
	var results []dtos.EmployeeBankAccountResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeBankAccount{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			eba.id, eba.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eba.bank_id, b.bank_name,
			eba.account_number, eba.account_name, eba.created_at, eba.updated_at,
			eba.created_by, eba.updated_by
		FROM employee_bank_accounts eba
		LEFT JOIN banks b ON eba.bank_id = b.id
		LEFT JOIN employees e ON eba.employee_id = e.id
		WHERE eba.deleted_at IS NULL AND (? = '' OR eba.employee_id = ?)
		ORDER BY eba.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, employeeID, employeeID, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *EmployeeBankAccountService) GetAccount(id string) (*dtos.EmployeeBankAccountResponse, error) {
	var result dtos.EmployeeBankAccountResponse
	query := `
		SELECT 
			eba.id, eba.employee_id, e.employee_no, CONCAT(e.first_name, ' ', e.surname) as employee_name,
			eba.bank_id, b.bank_name,
			eba.account_number, eba.account_name, eba.created_at, eba.updated_at
		FROM employee_bank_accounts eba
		LEFT JOIN banks b ON eba.bank_id = b.id
		LEFT JOIN employees e ON eba.employee_id = e.id
		WHERE eba.id = ? AND eba.deleted_at IS NULL
		LIMIT 1
	`
	err := db.DB.Raw(query, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (s *EmployeeBankAccountService) UpdateAccount(id string, req dtos.UpdateEmployeeBankAccountRequest, userID uint64) error {
	var account models.EmployeeBankAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"bank_id":        req.BankID,
		"account_number": req.AccountNumber,
		"account_name":   req.AccountName,
		"updated_by":     userID,
	}

	return db.DB.Model(&account).Updates(updates).Error
}

func (s *EmployeeBankAccountService) DeleteAccount(id string, userID uint64) error {
	// Audit the update before soft delete
	return db.DB.Model(&models.EmployeeBankAccount{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.EmployeeBankAccount{}).Error
}
