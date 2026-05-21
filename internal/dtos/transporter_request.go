package dtos

import (
	"mime/multipart"
)

type CreateTransporterRequest struct {
	TransporterNo string `json:"transporter_no" form:"transporter_no" validate:"required"`

	Category string `json:"transporter_category" form:"transporter_category" validate:"required,oneof=INDIVIDUAL COMPANY"`

	PrimaryPhone string `json:"primary_phone" form:"primary_phone" validate:"required,max=15"`

	EmailAddress string `json:"email_address" form:"email_address" validate:"omitempty,email"`

	RouteID uint64 `json:"route_id" form:"route_id"`

	Status string `json:"status" form:"status" validate:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED BLACKLISTED"`

	Restricted bool `json:"restricted" form:"restricted"`

	// -----------------------------------
	// Individual Transporter Fields
	// -----------------------------------

	FirstName string `json:"first_name" form:"first_name"`

	LastName string `json:"last_name" form:"last_name"`

	OtherNames string `json:"other_names" form:"other_names"`

	Gender string `json:"gender" form:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`

	DateOfBirth string `json:"date_of_birth" form:"date_of_birth"`

	NationalIDNo string `json:"national_id_no" form:"national_id_no"`

	KraPin string `json:"kra_pin" form:"kra_pin"`

	MaritalStatus string `json:"marital_status" form:"marital_status" validate:"omitempty,oneof=SINGLE MARRIED DIVORCED WIDOWED"`

	NextOfKinFullName string `json:"next_of_kin_full_name" form:"next_of_kin_full_name"`

	NextOfKinPhone string `json:"next_of_kin_phone" form:"next_of_kin_phone"`

	PassportPhoto *multipart.FileHeader `json:"passport_photo" form:"passport_photo" binding:"-"`

	IDFrontPhoto *multipart.FileHeader `json:"id_front_photo" form:"id_front_photo" binding:"-"`

	IDBackPhoto *multipart.FileHeader `json:"id_back_photo" form:"id_back_photo" binding:"-"`

	// -----------------------------------
	// Company Transporter Fields
	// -----------------------------------

	CompanyName string `json:"company_name" form:"company_name"`

	RegistrationNo string `json:"registration_no" form:"registration_no"`

	ContactPersonName string `json:"contact_person_name" form:"contact_person_name"`

	ContactPersonPhone string `json:"contact_person_phone" form:"contact_person_phone"`

	PostalAddress string `json:"postal_address" form:"postal_address"`

	PostalCode string `json:"postal_code" form:"postal_code"`

	Town string `json:"town" form:"town"`

	CertificateOfIncorporation *multipart.FileHeader `json:"certificate_of_incorporation" form:"certificate_of_incorporation" binding:"-"`
}

type UpdateTransporterRequest struct {
	TransporterTypeID uint64 `json:"transporter_type_id"`
	PrimaryPhone      string `json:"primary_phone" validate:"required,max=15"`
	EmailAddress      string `json:"email_address" validate:"omitempty,email"`
	RouteID           uint64 `json:"route_id"`
	Status            string `json:"status" validate:"required,oneof=ACTIVE INACTIVE SUSPENDED BLACKLISTED"`
	Restricted        bool   `json:"restricted"`

	// Fields for subtypes would typically be updated via separate specialized endpoints
	// or handled here depending on your frontend preference.
}

type UpdateIndividualTransporterRequest struct {
	RouteID           uint64 `json:"route_id"`
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	OtherNames        string `json:"other_names"`
	Gender            string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	DateOfBirth       string `json:"date_of_birth" validate:"required"`
	NationalIDNo      string `json:"national_id_no" validate:"required"`
	KraPin            string `json:"kra_pin"`
	MaritalStatus     string `json:"marital_status" validate:"omitempty,oneof=SINGLE MARRIED DIVORCED WIDOWED"`
	NextOfKinFullName string `json:"next_of_kin_full_name"`
	NextOfKinPhone    string `json:"next_of_kin_phone"`
	PassportPhoto     string `json:"passport_photo"`
	IDFrontPhoto      string `json:"id_front_photo"`
	IDBackPhoto       string `json:"id_back_photo"`
}

type UpdateCompanyTransporterRequest struct {
	RouteID                    uint64 `json:"route_id"`
	CompanyName                string `json:"company_name" validate:"required"`
	RegistrationNo             string `json:"registration_no" validate:"required"`
	KraPin                     string `json:"kra_pin"`
	ContactPersonName          string `json:"contact_person_name"`
	ContactPersonPhone         string `json:"contact_person_phone"`
	PostalAddress              string `json:"postal_address"`
	PostalCode                 string `json:"postal_code"`
	Town                       string `json:"town"`
	CertificateOfIncorporation string `json:"certificate_of_incorporation"`
}
