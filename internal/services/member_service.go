package services

import (
	"context"
	"mime/multipart"
	"strings"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	models "github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
	"gorm.io/gorm"
)

type MemberService struct{}

func NewMemberService() *MemberService {
	return &MemberService{}
}

// Create user
func (s *MemberService) CreateMember(
	ctx context.Context,
	req dtos.CreateMemberRequest) (*models.Member, error) {

	memberNo := req.MemberNo
	if memberNo == "" {
		memberNo = utils.GenerateMemberNo()
	}

	idFrontPath, err := utils.SaveFile(req.IDFrontPhoto, "members")
	if err != nil {
		return nil, err
	}

	idBackPath, err := utils.SaveFile(req.IDBackPhoto, "members")
	if err != nil {
		return nil, err
	}

	passportPath, err := utils.SaveFile(req.PassportPhoto, "members")
	if err != nil {
		return nil, err
	}

	newMember := &models.Member{
		MemberNo:       memberNo,
		MemberTypeID:   req.MemberTypeID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		OtherNames:     req.OtherNames,
		RouteID:        req.RouteID,
		DateOfBirth:    req.DOB,
		IDNumber:       req.IDNo,
		Gender:         req.Gender,
		BirthCity:      req.BirthCity,
		PrimaryPhone:   utils.NormalizePhone(req.PrimaryPhone),
		SecondaryPhone: utils.NormalizePhone(req.SecondaryPhone),
		Email:          req.Email,
		NumberOfCows:   req.NumberOfCows,
		IdFrontPhoto:   idFrontPath,
		IdBackPhoto:    idBackPath,
		PassportPhoto:  passportPath,
		IdDateOfIssue:  utils.ParseDate(req.IDDateOfIssue),
		TaxNumber:      req.TaxNumber,
		MaritalStatus:  req.MaritalStatus,
		Title:          req.Title,

		NextOfKinFullName: req.NextOfKinFullName,
		NextOfKinPhone:    utils.NormalizePhone(req.NextOfKinPhone),
		Status:            "PENDING",
		DateRegistered:    time.Now(),
	}

	if err := db.DB.Create(newMember).Error; err != nil {
		return nil, err
	}

	return newMember, nil
}

// Get all users
func (s *MemberService) GetMembers(page, limit int, memberNo, primaryPhone, memberTypeID, routeID string) ([]dtos.MemberResponse, int64, error) {
	var results []dtos.MemberResponse
	var total int64

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Base query for data
	baseQuery := `
		SELECT 
			m.id, m.member_no, m.member_type_id, mt.name AS member_type_name,
			m.first_name, m.last_name, m.other_names, m.route_id, r.route_name,
			m.date_of_birth, m.id_no, m.gender, m.birth_city,
			m.primary_phone, m.secondary_phone, m.email, m.number_of_cows,
			m.id_front_photo, m.id_back_photo, m.passport_photo, m.id_date_of_issue,
			m.tax_number, m.marital_status, m.title,
			m.next_of_kin_full_name, m.next_of_kin_phone, m.status, m.date_registered,
			m.created_at, m.updated_at
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
	`
	// Base query for count
	baseCountQuery := `
		SELECT COUNT(*)
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
	`

	var args []interface{}

	whereClauses := []string{"m.deleted_at IS NULL"}

	if memberNo != "" {
		whereClauses = append(whereClauses, "m.member_no LIKE ?")
		args = append(args, "%"+memberNo+"%")
	}
	if primaryPhone != "" {
		primaryPhone = utils.NormalizePhone(primaryPhone)
		whereClauses = append(whereClauses, "m.primary_phone LIKE ?")
		args = append(args, "%"+primaryPhone+"%")
	}
	if memberTypeID != "" {
		whereClauses = append(whereClauses, "m.member_type_id = ?")
		args = append(args, memberTypeID)
	}
	if routeID != "" {
		whereClauses = append(whereClauses, "m.route_id = ?")
		args = append(args, routeID)
	}

	whereSql := " WHERE " + strings.Join(whereClauses, " AND ")

	// Execute count query first to get filtered total
	if err := db.DB.Raw(baseCountQuery+whereSql, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Execute main data query
	fullQuery := baseQuery + whereSql + " ORDER BY m.id DESC LIMIT ? OFFSET ?"
	if err := db.DB.Raw(fullQuery, append(args, limit, offset)...).Scan(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

// Get single user
func (s *MemberService) GetMember(id string) (*dtos.MemberResponse, error) {
	var result dtos.MemberResponse
	query := `
		SELECT 
			m.id, m.member_no, m.member_type_id, mt.name AS member_type_name,
			m.first_name, m.last_name, m.other_names, m.route_id, r.route_name,
			m.date_of_birth, m.id_no, m.gender, m.birth_city,
			m.primary_phone, m.secondary_phone, m.email, m.number_of_cows,
			m.id_front_photo, m.id_back_photo, m.passport_photo, m.id_date_of_issue,
			m.tax_number, m.marital_status, m.title,
			m.next_of_kin_full_name, m.next_of_kin_phone, m.status, m.date_registered,
			m.created_at, m.updated_at
		FROM member_registrations m
		LEFT JOIN member_types mt ON m.member_type_id = mt.id
		LEFT JOIN routes r ON m.route_id = r.id
		WHERE m.id = ? AND m.deleted_at IS NULL
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

func (s *MemberService) UpdateMember(
	id string,
	req dtos.UpdateMemberRequest,
	idFront *multipart.FileHeader,
	idBack *multipart.FileHeader,
	passport *multipart.FileHeader,
) error {
	var member models.Member

	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}

	// Handle File Updates (only if provided)
	if idFront != nil {
		path, err := utils.SaveFile(idFront, "members")
		if err == nil {
			member.IdFrontPhoto = path
		}
	}
	if idBack != nil {
		path, err := utils.SaveFile(idBack, "members")
		if err == nil {
			member.IdBackPhoto = path
		}
	}
	if passport != nil {
		path, err := utils.SaveFile(passport, "members")
		if err == nil {
			member.PassportPhoto = path
		}
	}

	// Update Fields mapping according to CreateMember
	member.MemberTypeID = req.MemberTypeID
	member.FirstName = req.FirstName
	member.LastName = req.LastName
	member.OtherNames = req.OtherNames
	member.RouteID = req.RouteID
	member.DateOfBirth = req.DOB
	member.IDNumber = req.IDNo
	member.Gender = req.Gender
	member.BirthCity = req.BirthCity
	member.PrimaryPhone = utils.NormalizePhone(req.PrimaryPhone)
	member.SecondaryPhone = utils.NormalizePhone(req.SecondaryPhone)
	member.Email = req.Email
	member.NumberOfCows = req.NumberOfCows
	member.IdDateOfIssue = utils.ParseDate(req.IDDateOfIssue)
	member.TaxNumber = req.TaxNumber
	member.MaritalStatus = req.MaritalStatus
	member.Title = req.Title
	member.NextOfKinFullName = req.NextOfKinFullName
	member.NextOfKinPhone = utils.NormalizePhone(req.NextOfKinPhone)

	return db.DB.Save(&member).Error
}

func (s *MemberService) DeleteMember(id string) error {
	var member models.Member

	// ensure record exists
	if err := db.DB.First(&member, id).Error; err != nil {
		return err
	}

	return db.DB.Delete(&member).Error
}
