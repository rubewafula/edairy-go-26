package dtos

import (
	"mime/multipart"
)

type CreateMemberRequest struct {
	MemberTypeID uint64 `json:"member_type_id" form:"member_type_id" validate:"required"`
	FirstName    string `json:"first_name" form:"first_name" validate:"required,max=128"`
	LastName     string `json:"last_name" form:"last_name" validate:"required,max=128"`
	OtherNames   string `json:"other_names" form:"other_names" validate:"max=128"`
	RouteID      uint64 `json:"route_id" form:"route_id"`
	DOB          string `json:"date_of_birth" form:"date_of_birth" validate:"required"`

	IDNo      string `json:"id_no" form:"id_no" validate:"required,max=25"`
	MemberNo  string `json:"member_no" form:"member_no"`
	BirthCity string `json:"birth_city" form:"birth_city"`
	Gender    string `json:"gender" form:"gender" validate:"required,oneof=MALE FEMALE"`

	PrimaryPhone   string `json:"primary_phone" form:"primary_phone" validate:"required,max=15"`
	SecondaryPhone string `json:"secondary_phone" form:"secondary_phone"`
	Email          string `json:"email" form:"email" validate:"omitempty,email"`
	NumberOfCows   int    `json:"number_of_cows" form:"number_of_cows" validate:"required,min=0"`

	IDFrontPhoto  *multipart.FileHeader `form:"id_front_photo" validate:"required"`
	IDBackPhoto   *multipart.FileHeader `form:"id_back_photo" validate:"required"`
	PassportPhoto *multipart.FileHeader `form:"passport_photo" validate:"required"`
	IDDateOfIssue string                `json:"id_date_of_issue" form:"id_date_of_issue" validate:"required"`

	TaxNumber     string                         `json:"tax_number" form:"tax_number"`
	MaritalStatus string                         `json:"marital_status" form:"marital_status"` //
	Title         string                         `json:"title" form:"title"`                   //
	NextOfKins    []CreateMemberNextOfKinRequest `json:"next_of_kins" form:"next_of_kins"`

	NextOfKinFullName string `json:"next_of_kin_full_name" form:"next_of_kin_full_name"`
	NextOfKinPhone    string `json:"next_of_kin_phone" form:"next_of_kin_phone"`

	BankID      uint64 `json:"bank_id" form:"bank_id"`
	BankBranch  string `json:"bank_branch" form:"bank_branch"`
	AccountNo   string `json:"account_no" form:"account_no"`
	AccountName string `json:"account_name" form:"account_name"`
}
