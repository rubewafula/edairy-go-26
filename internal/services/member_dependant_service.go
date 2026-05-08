package services

import (
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/dtos"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/utils"
)

type MemberDependantService struct{}

func NewMemberDependantService() *MemberDependantService {
	return &MemberDependantService{}
}

func (s *MemberDependantService) CreateMemberDependant(req dtos.CreateMemberDependantRequest) (*models.MemberDependant, error) {
	dependant := &models.MemberDependant{
		MemberID:     req.MemberID,
		Name:         req.Name,
		NationalID:   req.NationalID,
		Relationship: req.Relationship,
		MobileNo:     req.MobileNo,
		Gender:       req.Gender,
		DateOfBirth:  utils.ParseDate(req.DateOfBirth),
		BirthCertNo:  req.BirthCertNo,
		Email:        req.Email,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		Town:         req.Town,
	}

	if err := db.DB.Create(dependant).Error; err != nil {
		return nil, err
	}
	return dependant, nil
}

func (s *MemberDependantService) GetMemberDependants() ([]models.MemberDependant, int64, error) {
	var dependants []models.MemberDependant
	var total int64
	db.DB.Model(&models.MemberDependant{}).Count(&total)
	err := db.DB.Find(&dependants).Error
	return dependants, total, err
}

func (s *MemberDependantService) GetMemberDependant(id string) (*models.MemberDependant, error) {
	var dependant models.MemberDependant
	if err := db.DB.First(&dependant, id).Error; err != nil {
		return nil, err
	}
	return &dependant, nil
}

func (s *MemberDependantService) UpdateMemberDependant(id string, req dtos.UpdateMemberDependantRequest) error {
	var dependant models.MemberDependant
	if err := db.DB.First(&dependant, id).Error; err != nil {
		return err
	}

	dependant.MemberID = req.MemberID
	dependant.Name = req.Name
	dependant.NationalID = req.NationalID
	dependant.Relationship = req.Relationship
	dependant.MobileNo = req.MobileNo
	dependant.Gender = req.Gender
	dependant.DateOfBirth = utils.ParseDate(req.DateOfBirth)
	dependant.BirthCertNo = req.BirthCertNo
	dependant.Email = req.Email
	dependant.Address = req.Address
	dependant.PostalCode = req.PostalCode
	dependant.Town = req.Town

	return db.DB.Save(&dependant).Error
}

func (s *MemberDependantService) DeleteMemberDependant(id string) error {
	return db.DB.Delete(&models.MemberDependant{}, id).Error
}
