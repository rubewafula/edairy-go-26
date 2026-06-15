package services

import (
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) CreateAccount(req dtos.CreateAccountRequest, userID uint64) (*models.Account, error) {
	account := &models.Account{
		BaseModel:         models.BaseModel{CreatedBy: userID},
		AccountCode:       req.AccountCode,
		Name:              req.Name,
		Description:       req.Description,
		AccountCategoryID: req.AccountCategoryID,
		ParentAccountID:   req.ParentAccountID,
		IsPostable:        req.IsPostable,
		IsActive:          req.IsActive,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAccounts(page, limit int) ([]dtos.AccountResponse, int64, error) {
	var results []dtos.AccountResponse
	var total int64

	db.DB.Model(&models.Account{}).Count(&total)
	offset := (page - 1) * limit

	query := `
		SELECT 
			a.id, a.account_code, a.name, a.description, a.account_category_id,
			ac.name as account_category_name, a.parent_account_id, pa.name as parent_account_name,
			a.is_postable, a.is_active, a.created_at, a.updated_at, a.created_by, a.updated_by
		FROM accounts a
		LEFT JOIN account_categories ac ON a.account_category_id = ac.id
		LEFT JOIN accounts pa ON a.parent_account_id = pa.id
		WHERE a.deleted_at IS NULL
		ORDER BY a.id DESC
		LIMIT ? OFFSET ?
	`
	err := db.DB.Raw(query, limit, offset).Scan(&results).Error
	return results, total, err
}

func (s *AccountService) GetAccount(id string) (*dtos.AccountResponse, error) {
	var result dtos.AccountResponse
	query := `
		SELECT 
			a.id, a.account_code, a.name, a.description, a.account_category_id,
			ac.name as account_category_name, a.parent_account_id, pa.name as parent_account_name,
			a.is_postable, a.is_active, a.created_at, a.updated_at, a.created_by, a.updated_by
		FROM accounts a
		LEFT JOIN account_categories ac ON a.account_category_id = ac.id
		LEFT JOIN accounts pa ON a.parent_account_id = pa.id
		WHERE a.id = ? AND a.deleted_at IS NULL
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

func (s *AccountService) UpdateAccount(id string, req dtos.UpdateAccountRequest, userID uint64) error {
	var account models.Account
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"account_code":        req.AccountCode,
		"name":                req.Name,
		"description":         req.Description,
		"account_category_id": req.AccountCategoryID,
		"parent_account_id":   req.ParentAccountID,
		"is_postable":         req.IsPostable,
		"is_active":           req.IsActive,
		"updated_by":          userID,
	}

	return db.DB.Model(&account).Updates(updates).Error
}

func (s *AccountService) DeleteAccount(id string, userID uint64) error {
	return db.DB.Model(&models.Account{}).Where("id = ?", id).Update("updated_by", userID).Delete(&models.Account{}).Error
}

func (s *AccountService) GetTrialBalance() (*dtos.TrialBalanceResponse, error) {
	var items []dtos.TrialBalanceItem

	query := `
		SELECT 
			a.id as account_id, 
			a.account_code, 
			a.name as account_name,
			SUM(COALESCE(gle.debit, 0)) as total_debit,
			SUM(COALESCE(gle.credit, 0)) as total_credit,
			(SUM(COALESCE(gle.debit, 0)) - SUM(COALESCE(gle.credit, 0))) as balance
		FROM accounts a
		INNER JOIN general_ledger_entries gle ON a.id = gle.account_id
		WHERE a.deleted_at IS NULL AND gle.deleted_at IS NULL
		GROUP BY a.id, a.account_code, a.name
		ORDER BY a.account_code ASC
	`

	if err := db.DB.Raw(query).Scan(&items).Error; err != nil {
		return nil, err
	}

	var totalDebits, totalCredits float64
	for _, item := range items {
		totalDebits += item.TotalDebit
		totalCredits += item.TotalCredit
	}

	return &dtos.TrialBalanceResponse{
		Items:        items,
		TotalDebits:  totalDebits,
		TotalCredits: totalCredits,
	}, nil
}

func (s *AccountService) GetProfitLoss() (*dtos.ProfitLossResponse, error) {
	var rawItems []struct {
		AccountID       uint64
		AccountCode     string
		AccountName     string
		CategoryName    string
		AccountTypeName string
		TotalDebit      float64
		TotalCredit     float64
	}

	query := `
		SELECT 
			a.id as account_id, 
			a.account_code, 
			a.name as account_name,
			ac.name as category_name,
			at.name as account_type_name,
			SUM(COALESCE(gle.debit, 0)) as total_debit,
			SUM(COALESCE(gle.credit, 0)) as total_credit
		FROM accounts a
		INNER JOIN account_categories ac ON a.account_category_id = ac.id
		INNER JOIN account_types at ON ac.account_type_id = at.id
		INNER JOIN general_ledger_entries gle ON a.id = gle.account_id
		WHERE a.deleted_at IS NULL AND gle.deleted_at IS NULL
		  AND (at.name LIKE '%Revenue%' OR at.name LIKE '%Income%' OR at.name LIKE '%Expense%')
		GROUP BY a.id, a.account_code, a.name, ac.name, at.name
		ORDER BY at.name, a.account_code ASC
	`

	if err := db.DB.Raw(query).Scan(&rawItems).Error; err != nil {
		return nil, err
	}

	resp := &dtos.ProfitLossResponse{
		RevenueItems: []dtos.ProfitLossItem{},
		ExpenseItems: []dtos.ProfitLossItem{},
	}

	for _, item := range rawItems {
		typeName := strings.ToLower(item.AccountTypeName)
		isRevenue := strings.Contains(typeName, "revenue") || strings.Contains(typeName, "income")

		amount := 0.0
		pLItem := dtos.ProfitLossItem{
			AccountID:    item.AccountID,
			AccountCode:  item.AccountCode,
			AccountName:  item.AccountName,
			CategoryName: item.CategoryName,
			TypeName:     item.AccountTypeName,
		}

		if isRevenue {
			amount = item.TotalCredit - item.TotalDebit
			pLItem.Amount = amount
			resp.RevenueItems = append(resp.RevenueItems, pLItem)
			resp.TotalRevenue += amount
		} else {
			amount = item.TotalDebit - item.TotalCredit
			pLItem.Amount = amount
			resp.ExpenseItems = append(resp.ExpenseItems, pLItem)
			resp.TotalExpenses += amount
		}
	}

	resp.NetProfit = resp.TotalRevenue - resp.TotalExpenses
	return resp, nil
}

func (s *AccountService) GetBalanceSheet() (*dtos.BalanceSheetResponse, error) {
	// 1. Get Net Profit from P&L to account for current period earnings in Equity
	pl, err := s.GetProfitLoss()
	if err != nil {
		return nil, err
	}

	var rawItems []struct {
		AccountID       uint64
		AccountCode     string
		AccountName     string
		CategoryName    string
		AccountTypeName string
		TotalDebit      float64
		TotalCredit     float64
	}

	query := `
		SELECT 
			a.id as account_id, 
			a.account_code, 
			a.name as account_name,
			ac.name as category_name,
			at.name as account_type_name,
			SUM(COALESCE(gle.debit, 0)) as total_debit,
			SUM(COALESCE(gle.credit, 0)) as total_credit
		FROM accounts a
		INNER JOIN account_categories ac ON a.account_category_id = ac.id
		INNER JOIN account_types at ON ac.account_type_id = at.id
		INNER JOIN general_ledger_entries gle ON a.id = gle.account_id
		WHERE a.deleted_at IS NULL AND gle.deleted_at IS NULL
		  AND (at.name LIKE '%Asset%' OR at.name LIKE '%Liability%' OR at.name LIKE '%Equity%' OR at.name LIKE '%Capital%')
		GROUP BY a.id, a.account_code, a.name, ac.name, at.name
		ORDER BY at.name, a.account_code ASC
	`

	if err := db.DB.Raw(query).Scan(&rawItems).Error; err != nil {
		return nil, err
	}

	resp := &dtos.BalanceSheetResponse{
		AssetItems:     []dtos.BalanceSheetItem{},
		LiabilityItems: []dtos.BalanceSheetItem{},
		EquityItems:    []dtos.BalanceSheetItem{},
	}

	for _, item := range rawItems {
		typeName := strings.ToLower(item.AccountTypeName)
		amount := 0.0
		bsItem := dtos.BalanceSheetItem{
			AccountID:    item.AccountID,
			AccountCode:  item.AccountCode,
			AccountName:  item.AccountName,
			CategoryName: item.CategoryName,
			TypeName:     item.AccountTypeName,
		}

		if strings.Contains(typeName, "asset") {
			amount = item.TotalDebit - item.TotalCredit
			bsItem.Amount = amount
			resp.AssetItems = append(resp.AssetItems, bsItem)
			resp.TotalAssets += amount
		} else if strings.Contains(typeName, "liability") {
			amount = item.TotalCredit - item.TotalDebit
			bsItem.Amount = amount
			resp.LiabilityItems = append(resp.LiabilityItems, bsItem)
			resp.TotalLiabilities += amount
		} else if strings.Contains(typeName, "equity") || strings.Contains(typeName, "capital") {
			amount = item.TotalCredit - item.TotalDebit
			bsItem.Amount = amount
			resp.EquityItems = append(resp.EquityItems, bsItem)
			resp.TotalEquity += amount
		}
	}

	// 2. Add Net Profit to Equity (Retained Earnings for current period)
	if pl.NetProfit != 0 {
		resp.EquityItems = append(resp.EquityItems, dtos.BalanceSheetItem{
			AccountName: "Retained Earnings (Current Period Profit/Loss)",
			Amount:      pl.NetProfit,
			TypeName:    "Equity",
		})
		resp.TotalEquity += pl.NetProfit
	}

	resp.TotalLiabilitiesEquity = resp.TotalLiabilities + resp.TotalEquity
	return resp, nil
}
