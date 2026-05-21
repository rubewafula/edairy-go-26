package dtos

import "time"

type TransporterDriverResponse struct {
	ID                       uint64     `json:"id"`
	TransporterID            uint64     `json:"transporter_id"`
	DriverNo                 string     `json:"driver_no"`
	FirstName                string     `json:"first_name"`
	LastName                 string     `json:"last_name"`
	OtherNames               string     `json:"other_names"`
	Gender                   string     `json:"gender"`
	DateOfBirth              *time.Time `json:"date_of_birth"`
	NationalIDNo             string     `json:"national_id_no"`
	KraPin                   string     `json:"kra_pin"`
	PrimaryPhone             string     `json:"primary_phone"`
	SecondaryPhone           string     `json:"secondary_phone"`
	EmailAddress             string     `json:"email_address"`
	DrivingLicenseNo         string     `json:"driving_license_no"`
	DrivingLicenseClass      string     `json:"driving_license_class"`
	DrivingLicenseExpiry     *time.Time `json:"driving_license_expiry"`
	EmploymentDate           *time.Time `json:"employment_date"`
	Status                   string     `json:"status"`
	PassportPhoto            string     `json:"passport_photo"`
	DrivingLicenseFrontPhoto string     `json:"driving_license_front_photo"`
	DrivingLicenseBackPhoto  string     `json:"driving_license_back_photo"`
	NextOfKinFullName        string     `json:"next_of_kin_full_name"`
	NextOfKinPhone           string     `json:"next_of_kin_phone"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}
