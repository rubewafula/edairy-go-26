package services

import (
	"strings"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"gorm.io/gorm"
)

type MemberNextOfKinService struct{}

func NewMemberNextOfKinService() *MemberNextOfKinService {
	return &MemberNextOfKinService{}
}

func (s *MemberNextOfKinService) CreateMemberNextOfKin(req dtos.CreateMemberNextOfKinRequest, userID uint64) (*models.MemberNextOfKin, error) {
	nextOfKin := &models.MemberNextOfKin{
		BaseModel:              models.BaseModel{CreatedBy: userID},
		MemberID:               req.MemberID,
		FullName:               req.FullName,
		Relationship:           req.Relationship,
		PhoneNumber:            req.PhoneNumber,
		AlternativePhoneNumber: req.AlternativePhoneNumber,
		EmailAddress:           req.EmailAddress,
		NationalIDNo:           req.NationalIDNo,
		PostalAddress:          req.PostalAddress,
		PhysicalAddress:        req.PhysicalAddress,
		Occupation:             req.Occupation,
		IsPrimary:              req.IsPrimary,
		Status:                 req.Status,
		Remarks:                req.Remarks,
	}

	if err := db.DB.Create(nextOfKin).Error; err != nil {
		return nil, err
	}
	return nextOfKin, nil
}

func (s *MemberNextOfKinService) GetMemberNextOfKins(page, limit int, memberID string) ([]dtos.MemberNextOfKinResponse, int64, error) {
	var results []dtos.MemberNextOfKinResponse
	var total int64

	offset := (page - 1) * limit

	baseQuery := `
		SELECT 
			mnok.*, 
			m.member_no, 
			CONCAT(m.first_name, ' ', m.last_name, ' ', COALESCE(m.other_names, '')) AS member_full_name
		FROM member_next_of_kins mnok
		JOIN member_registrations m ON mnok.member_id = m.id
	`
	baseCountQuery := `
		SELECT COUNT(*)
		FROM member_next_of_kins mnok
		JOIN member_registrations m ON mnok.member_id = m.id
	`

	var args []interface{}
	whereClauses := []string{"mnok.deleted_at IS NULL"}

	if memberID != "" {
		whereClauses = append(whereClauses, "mnok.member_id = ?")
		args = append(args, memberID)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	if err := db.DB.Raw(baseCountQuery+whereSql, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	fullQuery := baseQuery + whereSql + " ORDER BY mnok.id DESC LIMIT ? OFFSET ?"
	if err := db.DB.Raw(fullQuery, append(args, limit, offset)...).Scan(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (s *MemberNextOfKinService) GetMemberNextOfKin(id string) (*dtos.MemberNextOfKinResponse, error) {
	var result dtos.MemberNextOfKinResponse
	query := `
		SELECT 
			mnok.*, 
			m.member_no, 
			CONCAT(m.first_name, ' ', m.last_name, ' ', COALESCE(m.other_names, '')) AS member_full_name
		FROM member_next_of_kins mnok
		JOIN member_registrations m ON mnok.member_id = m.id
		WHERE mnok.id = ? AND mnok.deleted_at IS NULL
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

func (s *MemberNextOfKinService) UpdateMemberNextOfKin(id string, req dtos.UpdateMemberNextOfKinRequest, userID uint64) error {
	var nextOfKin models.MemberNextOfKin
	if err := db.DB.First(&nextOfKin, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"updated_by": userID,
	}

	if req.MemberID != 0 {
		updates["member_id"] = req.MemberID
	}
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.Relationship != nil {
		updates["relationship"] = req.Relationship
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = req.PhoneNumber
	}
	if req.AlternativePhoneNumber != nil {
		updates["alternative_phone_number"] = req.AlternativePhoneNumber
	}
	if req.EmailAddress != nil {
		updates["email_address"] = req.EmailAddress
	}
	if req.NationalIDNo != nil {
		updates["national_id_no"] = req.NationalIDNo
	}
	if req.PostalAddress != nil {
		updates["postal_address"] = req.PostalAddress
	}
	if req.PhysicalAddress != nil {
		updates["physical_address"] = req.PhysicalAddress
	}
	if req.Occupation != nil {
		updates["occupation"] = req.Occupation
	}
	updates["is_primary"] = req.IsPrimary
	updates["status"] = req.Status
	if req.Remarks != nil {
		updates["remarks"] = req.Remarks
	}

	return db.DB.Model(&nextOfKin).Updates(updates).Error
}

func (s *MemberNextOfKinService) DeleteMemberNextOfKin(id string, userID uint64) error {
	var nextOfKin models.MemberNextOfKin
	if err := db.DB.First(&nextOfKin, id).Error; err != nil {
		return err
	}
	return db.DB.Model(&nextOfKin).Update("updated_by", userID).Delete(&nextOfKin).Error
}
