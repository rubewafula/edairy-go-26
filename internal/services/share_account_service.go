package services

import (
	"fmt"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type ShareAccountService struct{}

func NewShareAccountService() *ShareAccountService {
	return &ShareAccountService{}
}

func (s *ShareAccountService) CreateAccount(req dtos.CreateShareAccountRequest) (*models.ShareAccount, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	openedAt := time.Now()
	if req.OpenedAt != "" {
		openedAt = utils.ParseDate(req.OpenedAt)
	}

	account := &models.ShareAccount{
		MemberID:    req.MemberID,
		ShareTypeID: req.ShareTypeID,
		Status:      status,
		OpenedAt:    openedAt,
	}

	if err := db.DB.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (s *ShareAccountService) GetShareAccounts(memberID, shareTypeID string) ([]dtos.ShareAccountResponse, int64, error) {
	var results []dtos.ShareAccountResponse
	var total int64

	countQuery := db.DB.Model(&models.ShareAccount{})
	if memberID != "" {
		countQuery = countQuery.Where("member_id = ?", memberID)
	}
	if shareTypeID != "" {
		countQuery = countQuery.Where("share_type_id = ?", shareTypeID)
	}
	countQuery.Count(&total)

	whereClause := "sa.deleted_at IS NULL"
	var args []interface{}
	if memberID != "" {
		whereClause += " AND sa.member_id = ?"
		args = append(args, memberID)
	}
	if shareTypeID != "" {
		whereClause += " AND sa.share_type_id = ?"
		args = append(args, shareTypeID)
	}

	query := fmt.Sprintf(`
		SELECT 
			sa.id, sa.member_id, m.member_no, m.first_name, m.last_name,
			sa.share_type_id, st.share_code, st.share_type AS share_type_name, st.description,
			sa.status, sa.share_units, sa.share_amount, sa.opened_at, sa.created_at, sa.updated_at
		FROM share_accounts sa
		LEFT JOIN member_registrations m ON sa.member_id = m.id
		LEFT JOIN share_types st ON sa.share_type_id = st.id
		WHERE %s
	`, whereClause)

	err := db.DB.Raw(query, args...).Scan(&results).Error
	return results, total, err
}

func (s *ShareAccountService) GetShareAccount(id string) (*dtos.ShareAccountResponse, error) {
	var result dtos.ShareAccountResponse
	query := `
		SELECT 
			sa.id, sa.member_id, m.member_no, m.first_name, m.last_name,
			sa.share_type_id, st.share_code, st.share_type AS share_type_name, st.description,
			sa.status, sa.share_units, sa.share_amount, sa.opened_at, sa.created_at, sa.updated_at
		FROM share_accounts sa
		LEFT JOIN member_registrations m ON sa.member_id = m.id
		LEFT JOIN share_types st ON sa.share_type_id = st.id
		WHERE sa.id = ? AND sa.deleted_at IS NULL
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

func (s *ShareAccountService) UpdateAccount(id string, req dtos.UpdateShareAccountRequest) error {
	var account models.ShareAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		return err
	}

	account.MemberID = req.MemberID
	account.ShareTypeID = req.ShareTypeID
	account.Status = req.Status
	if req.OpenedAt != "" {
		account.OpenedAt = utils.ParseDate(req.OpenedAt)
	}

	return db.DB.Save(&account).Error
}

func (s *ShareAccountService) DeleteAccount(id string) error {
	return db.DB.Delete(&models.ShareAccount{}, id).Error
}
