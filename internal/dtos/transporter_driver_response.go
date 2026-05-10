package dtos

import "time"

type TransporterDriverResponse struct {
	ID                       uint64    `json:"ID"`
	TransporterID            uint64    `json:"TransporterID"`
	DriverNo                 string    `json:"DriverNo"`
	FirstName                string    `json:"FirstName"`
	LastName                 string    `json:"LastName"`
	OtherNames               string    `json:"OtherNames"`
	Gender                   string    `json:"Gender"`
	DateOfBirth              time.Time `json:"DateOfBirth"`
	NationalIDNo             string    `json:"NationalIDNo"`
	KraPin                   string    `json:"KraPin"`
	PrimaryPhone             string    `json:"PrimaryPhone"`
	SecondaryPhone           string    `json:"SecondaryPhone"`
	EmailAddress             string    `json:"EmailAddress"`
	DrivingLicenseNo         string    `json:"DrivingLicenseNo"`
	DrivingLicenseClass      string    `json:"DrivingLicenseClass"`
	DrivingLicenseExpiry     time.Time `json:"DrivingLicenseExpiry"`
	EmploymentDate           time.Time `json:"EmploymentDate"`
	Status                   string    `json:"Status"`
	PassportPhoto            string    `json:"PassportPhoto"`
	DrivingLicenseFrontPhoto string    `json:"DrivingLicenseFrontPhoto"`
	DrivingLicenseBackPhoto  string    `json:"DrivingLicenseBackPhoto"`
	NextOfKinFullName        string    `json:"NextOfKinFullName"`
	NextOfKinPhone           string    `json:"NextOfKinPhone"`
	CreatedAt                time.Time `json:"CreatedAt"`
	UpdatedAt                time.Time `json:"UpdatedAt"`
}
