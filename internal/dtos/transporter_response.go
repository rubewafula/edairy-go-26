package dtos

import "time"

type IndividualTransporterResponse struct {
	ID                uint64    `json:"id"`
	TransporterID     uint64    `json:"transporter_id"`
	RouteID           uint64    `json:"route_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	OtherNames        string    `json:"other_names"`
	Gender            string    `json:"gender"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	NationalIDNo      string    `json:"national_id_no"`
	KraPin            string    `json:"kra_pin"`
	MaritalStatus     string    `json:"marital_status"`
	NextOfKinFullName string    `json:"next_of_kin_full_name"`
	NextOfKinPhone    string    `json:"next_of_kin_phone"`
	PassportPhoto     string    `json:"passport_photo"`
	IDFrontPhoto      string    `json:"id_front_photo"`
	IDBackPhoto       string    `json:"id_back_photo"`
}

type CompanyTransporterResponse struct {
	ID                         uint64 `json:"id"`
	TransporterID              uint64 `json:"transporter_id"`
	RouteID                    uint64 `json:"route_id"`
	CompanyName                string `json:"company_name"`
	RegistrationNo             string `json:"registration_no"`
	KraPin                     string `json:"kra_pin"`
	ContactPersonName          string `json:"contact_person_name"`
	ContactPersonPhone         string `json:"contact_person_phone"`
	PostalAddress              string `json:"postal_address"`
	PostalCode                 string `json:"postal_code"`
	Town                       string `json:"town"`
	CertificateOfIncorporation string `json:"certificate_of_incorporation"`
}

type TransporterResponse struct {
	ID                uint64                         `json:"id"`
	TransporterNo     string                         `json:"transporter_no"`
	Category          string                         `json:"category"`
	TransporterTypeID uint64                         `json:"transporter_type_id"`
	PrimaryPhone      string                         `json:"primary_phone"`
	EmailAddress      string                         `json:"email_address"`
	RouteID           uint64                         `json:"route_id"`
	Status            string                         `json:"status"`
	Restricted        bool                           `json:"restricted"`
	CreatedAt         time.Time                      `json:"created_at"`
	UpdatedAt         time.Time                      `json:"updated_at"`
	Individual        *IndividualTransporterResponse `json:"individual,omitempty"`
	Company           *CompanyTransporterResponse    `json:"company,omitempty"`
}
