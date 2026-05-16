package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type LoanManagementService struct{}

func NewLoanManagementService() *LoanManagementService {
	return &LoanManagementService{}
}

// --- LoanAccount CRUD ---
func (s *LoanManagementService) CreateLoanAccount(req dtos.CreateLoanAccountRequest, userID uint64) (*models.LoanAccount, error) {
	account := &models.LoanAccount{
		BaseModel:      models.BaseModel{CreatedBy: userID},
		CustomerID:     req.CustomerID,
		CustomerType:   req.CustomerType,
		ManuallyRatify: req.ManuallyRatify,
		NextLevel:      req.NextLevel,
		Status:         req.Status,
		LinkStatus:     req.LinkStatus,
		LivenessPassed: req.LivenessPassed,
		AstraRemarks:   req.AstraRemarks,
		AuthCreated:    req.AuthCreated,
		Locale:         req.Locale,
	}

	if req.AstraID != "" {
		account.AstraID = &req.AstraID
	}
	if req.CreditLimit != 0 {
		account.CreditLimit = &req.CreditLimit
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *LoanManagementService) GetLoanAccounts(page, limit int) ([]models.LoanAccount, int64, error) {
	var accounts []models.LoanAccount
	var total int64
	db.DB.Model(&models.LoanAccount{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&accounts).Error
	return accounts, total, err
}

func (s *LoanManagementService) GetLoanAccount(id string) (*models.LoanAccount, error) {
	var account models.LoanAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *LoanManagementService) UpdateLoanAccount(id string, req dtos.UpdateLoanAccountRequest, userID uint64) error {
	var account models.LoanAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"updated_by": userID,
	}
	return db.DB.Model(&account).Updates(req).Updates(updates).Error
}

func (s *LoanManagementService) DeleteLoanAccount(id string, userID uint64) error {
	return db.DB.Model(&models.LoanAccount{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.LoanAccount{}).Error
}

// --- LoanCallback CRUD ---
func (s *LoanManagementService) CreateLoanCallback(req dtos.CreateLoanCallbackRequest, userID uint64) (*models.LoanCallback, error) {
	callback := &models.LoanCallback{
		BaseModel: models.BaseModel{CreatedBy: userID},
		Detail:    req.Detail,
		LoanID:    req.LoanID,
		Type:      req.Type,
	}
	if err := db.DB.Create(callback).Error; err != nil {
		return nil, err
	}
	return callback, nil
}

func (s *LoanManagementService) GetLoanCallbacks(page, limit int) ([]models.LoanCallback, int64, error) {
	var callbacks []models.LoanCallback
	var total int64
	db.DB.Model(&models.LoanCallback{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&callbacks).Error
	return callbacks, total, err
}

func (s *LoanManagementService) GetLoanCallback(id string) (*models.LoanCallback, error) {
	var callback models.LoanCallback
	if err := db.DB.First(&callback, id).Error; err != nil {
		return nil, err
	}
	return &callback, nil
}

func (s *LoanManagementService) UpdateLoanCallback(id string, req dtos.UpdateLoanCallbackRequest, userID uint64) error {
	var callback models.LoanCallback
	if err := db.DB.First(&callback, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"detail":     req.Detail,
		"loan_id":    req.LoanID,
		"type":       req.Type,
		"updated_by": userID,
	}
	return db.DB.Model(&callback).Updates(updates).Error
}

func (s *LoanManagementService) DeleteLoanCallback(id string, userID uint64) error {
	var callback models.LoanCallback
	if err := db.DB.First(&callback, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&callback).Update("updated_by", userID).Delete(&callback).Error
}

// --- LoanOrganizationProfile CRUD ---
func (s *LoanManagementService) CreateLoanOrganizationProfile(req dtos.CreateLoanOrganizationProfileRequest, userID uint64) (*models.LoanOrganizationProfile, error) {
	profile := &models.LoanOrganizationProfile{
		BaseModel:       models.BaseModel{CreatedBy: userID},
		NextLevel:       req.NextLevel,
		AstraID:         req.AstraID,
		LinkStatus:      req.LinkStatus,
		UUID:            req.UUID,
		Version:         req.Version,
		ProductID:       req.ProductID,
		CompanyDetailID: req.CompanyDetailID,
		ManuallyRatify:  req.ManuallyRatify,
	}
	if err := db.DB.Create(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *LoanManagementService) GetLoanOrganizationProfiles(page, limit int) ([]models.LoanOrganizationProfile, int64, error) {
	var profiles []models.LoanOrganizationProfile
	var total int64
	db.DB.Model(&models.LoanOrganizationProfile{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&profiles).Error
	return profiles, total, err
}

func (s *LoanManagementService) GetLoanOrganizationProfile(id string) (*models.LoanOrganizationProfile, error) {
	var profile models.LoanOrganizationProfile
	if err := db.DB.First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *LoanManagementService) UpdateLoanOrganizationProfile(id string, req dtos.UpdateLoanOrganizationProfileRequest, userID uint64) error {
	var profile models.LoanOrganizationProfile
	if err := db.DB.First(&profile, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"next_level":        req.NextLevel,
		"astra_id":          req.AstraID,
		"link_status":       req.LinkStatus,
		"uuid":              req.UUID,
		"version":           req.Version,
		"product_id":        req.ProductID,
		"company_detail_id": req.CompanyDetailID,
		"manually_ratify":   req.ManuallyRatify,
		"updated_by":        userID,
	}
	return db.DB.Model(&profile).Updates(updates).Error
}

func (s *LoanManagementService) DeleteLoanOrganizationProfile(id string, userID uint64) error {
	var profile models.LoanOrganizationProfile
	if err := db.DB.First(&profile, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&profile).Update("updated_by", userID).Delete(&profile).Error
}

// --- LoanOriginationCallbackLog CRUD ---
func (s *LoanManagementService) CreateLoanOriginationCallbackLog(req dtos.CreateLoanOriginationCallbackLogRequest, userID uint64) (*models.LoanOriginationCallbackLog, error) {
	log := &models.LoanOriginationCallbackLog{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		AstraDetail: req.AstraDetail,
		SyncAttempt: req.SyncAttempt,
	}
	if err := db.DB.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (s *LoanManagementService) GetLoanOriginationCallbackLogs(page, limit int) ([]models.LoanOriginationCallbackLog, int64, error) {
	var logs []models.LoanOriginationCallbackLog
	var total int64
	db.DB.Model(&models.LoanOriginationCallbackLog{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

func (s *LoanManagementService) GetLoanOriginationCallbackLog(id string) (*models.LoanOriginationCallbackLog, error) {
	var log models.LoanOriginationCallbackLog
	if err := db.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (s *LoanManagementService) UpdateLoanOriginationCallbackLog(id string, req dtos.UpdateLoanOriginationCallbackLogRequest, userID uint64) error {
	var log models.LoanOriginationCallbackLog
	if err := db.DB.First(&log, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"astra_detail": req.AstraDetail,
		"sync_attempt": req.SyncAttempt,
		"updated_by":   userID,
	}
	return db.DB.Model(&log).Updates(updates).Error
}

func (s *LoanManagementService) DeleteLoanOriginationCallbackLog(id string, userID uint64) error {
	var log models.LoanOriginationCallbackLog
	if err := db.DB.First(&log, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&log).Update("updated_by", userID).Delete(&log).Error
}

// --- LoanTransaction CRUD ---
func (s *LoanManagementService) CreateLoanTransaction(req dtos.CreateLoanTransactionRequest, userID uint64) (*models.LoanTransaction, error) {
	transaction := &models.LoanTransaction{
		BaseModel:   models.BaseModel{CreatedBy: userID},
		LoanID:      req.LoanID,
		Amount:      req.Amount,
		Type:        req.Type,
		Reference:   req.Reference,
		Description: req.Description,
		Date:        utils.ParseDate(req.Date),
	}
	if err := db.DB.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *LoanManagementService) GetLoanTransactions(page, limit int) ([]models.LoanTransaction, int64, error) {
	var transactions []models.LoanTransaction
	var total int64
	db.DB.Model(&models.LoanTransaction{}).Count(&total)
	offset := (page - 1) * limit
	err := db.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&transactions).Error
	return transactions, total, err
}

func (s *LoanManagementService) GetLoanTransaction(id string) (*models.LoanTransaction, error) {
	var transaction models.LoanTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *LoanManagementService) GetLoanTransactionsByLoanID(loanID string) ([]models.LoanTransaction, error) {
	var transactions []models.LoanTransaction
	err := db.DB.Where("loan_id = ?", loanID).Find(&transactions).Error
	return transactions, err
}

func (s *LoanManagementService) UpdateLoanTransaction(id string, req dtos.UpdateLoanTransactionRequest, userID uint64) error {
	var transaction models.LoanTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return err
	}
	updates := map[string]interface{}{
		"loan_id":     req.LoanID,
		"amount":      req.Amount,
		"type":        req.Type,
		"reference":   req.Reference,
		"description": req.Description,
		"date":        utils.ParseDate(req.Date),
		"updated_by":  userID,
	}
	return db.DB.Model(&transaction).Updates(updates).Error
}

func (s *LoanManagementService) DeleteLoanTransaction(id string, userID uint64) error {
	var transaction models.LoanTransaction
	if err := db.DB.First(&transaction, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&transaction).Update("updated_by", userID).Delete(&transaction).Error
}
