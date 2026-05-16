package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type EmployeeContractService struct{}

func NewEmployeeContractService() *EmployeeContractService {
	return &EmployeeContractService{}
}

func (s *EmployeeContractService) CreateContract(req dtos.CreateEmployeeContractRequest, userID uint64) (*models.EmployeeContractDetail, error) {
	contract := &models.EmployeeContractDetail{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		EmployeeID:      req.EmployeeID,
		ContractType:    req.ContractType,
		ContractEndDate: utils.ParseDate(req.ContractEndDate),
		NoticePeriod:    req.NoticePeriod,
		RetirementDate:  utils.ParseDate(req.RetirementDate),
	}

	if err := db.DB.Create(contract).Error; err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *EmployeeContractService) GetContracts(employeeID string, page, limit int) ([]dtos.EmployeeContractResponse, int64, error) {
	var results []dtos.EmployeeContractResponse
	var total int64

	queryBuilder := db.DB.Model(&models.EmployeeContractDetail{})
	if employeeID != "" {
		queryBuilder = queryBuilder.Where("employee_id = ?", employeeID)
	}

	queryBuilder.Count(&total)
	offset := (page - 1) * limit

	err := db.DB.Model(&models.EmployeeContractDetail{}).
		Where("deleted_at IS NULL AND (? = '' OR employee_id = ?)", employeeID, employeeID).
		Limit(limit).Offset(offset).Order("id DESC").
		Scan(&results).Error

	return results, total, err
}

func (s *EmployeeContractService) GetContract(id string) (*dtos.EmployeeContractResponse, error) {
	var result dtos.EmployeeContractResponse
	err := db.DB.Model(&models.EmployeeContractDetail{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EmployeeContractService) UpdateContract(id string, req dtos.UpdateEmployeeContractRequest, userID uint64) error {
	var contract models.EmployeeContractDetail
	if err := db.DB.First(&contract, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"contract_type":     req.ContractType,
		"contract_end_date": utils.ParseDate(req.ContractEndDate),
		"notice_period":     req.NoticePeriod,
		"retirement_date":   utils.ParseDate(req.RetirementDate),
		"updated_by":        userID,
	}

	return db.DB.Model(&contract).Updates(updates).Error
}

func (s *EmployeeContractService) DeleteContract(id string, userID uint64) error {
	var contract models.EmployeeContractDetail
	if err := db.DB.First(&contract, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&contract).Update("updated_by", userID).Delete(&contract).Error
}
