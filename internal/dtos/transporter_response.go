package dtos

import "time"

type IndividualTransporterResponse struct {
	ID                uint64    `json:"ID"`
	TransporterID     uint64    `json:"TransporterID"`
	RouteID           uint64    `json:"RouteID"`
	FirstName         string    `json:"FirstName"`
	LastName          string    `json:"LastName"`
	OtherNames        string    `json:"OtherNames"`
	Gender            string    `json:"Gender"`
	DateOfBirth       time.Time `json:"DateOfBirth"`
	NationalIDNo      string    `json:"NationalIDNo"`
	KraPin            string    `json:"KraPin"`
	MaritalStatus     string    `json:"MaritalStatus"`
	NextOfKinFullName string    `json:"NextOfKinFullName"`
	NextOfKinPhone    string    `json:"NextOfKinPhone"`
	PassportPhoto     string    `json:"PassportPhoto"`
	IDFrontPhoto      string    `json:"IDFrontPhoto"`
	IDBackPhoto       string    `json:"IDBackPhoto"`
}

type CompanyTransporterResponse struct {
	ID                         uint64 `json:"ID"`
	TransporterID              uint64 `json:"TransporterID"`
	RouteID                    uint64 `json:"RouteID"`
	CompanyName                string `json:"CompanyName"`
	RegistrationNo             string `json:"RegistrationNo"`
	KraPin                     string `json:"KraPin"`
	ContactPersonName          string `json:"ContactPersonName"`
	ContactPersonPhone         string `json:"ContactPersonPhone"`
	PostalAddress              string `json:"PostalAddress"`
	PostalCode                 string `json:"PostalCode"`
	Town                       string `json:"Town"`
	CertificateOfIncorporation string `json:"CertificateOfIncorporation"`
}

type TransporterResponse struct {
	ID                uint64                         `json:"ID"`
	TransporterNo     string                         `json:"TransporterNo"`
	Category          string                         `json:"Category"`
	TransporterTypeID uint64                         `json:"TransporterTypeID"`
	PrimaryPhone      string                         `json:"PrimaryPhone"`
	EmailAddress      string                         `json:"EmailAddress"`
	RouteID           uint64                         `json:"RouteID"`
	Status            string                         `json:"Status"`
	Restricted        bool                           `json:"Restricted"`
	CreatedAt         time.Time                      `json:"CreatedAt"`
	UpdatedAt         time.Time                      `json:"UpdatedAt"`
	Individual        *IndividualTransporterResponse `json:"Individual,omitempty"`
	Company           *CompanyTransporterResponse    `json:"Company,omitempty"`
}
