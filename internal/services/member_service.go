package services

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	models "github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type MemberService struct{}

func NewMemberService() *MemberService {
	return &MemberService{}
}

// Create user
func (s *MemberService) CreateMember(
	ctx context.Context,
	req dtos.CreateMemberRequest,
	idFront *multipart.FileHeader,
	idBack *multipart.FileHeader,
	passport *multipart.FileHeader,
) (*models.Member, error) {

	memberNo := req.MemberNo
	if memberNo == "" {
		memberNo = utils.GenerateMemberNo()
	}

	idFrontPath, err := utils.SaveFile(idFront, "members")
	if err != nil {
		return nil, err
	}

	idBackPath, err := utils.SaveFile(idBack, "members")
	if err != nil {
		return nil, err
	}

	passportPath, err := utils.SaveFile(passport, "members")
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
func (s *MemberService) GetMembers() ([]models.Member, int64, error) {
	var members []models.Member
	var total int64
	db.DB.Model(&models.Member{}).Count(&total)
	err := db.DB.Find(&members).Error
	return members, total, err
}

// Get single user
func (s *MemberService) GetMember(id string) (models.Member, error) {
	var member models.Member
	err := db.DB.First(&member, id).Error
	return member, err
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
