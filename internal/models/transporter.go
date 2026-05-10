package models

import (
	"time"
)

type Transporter struct {
	BaseModel
	TransporterNo     string `gorm:"unique;not null" json:"transporter_no"`
	Category          string `gorm:"type:enum('INDIVIDUAL','COMPANY');not null" json:"transporter_category"`
	TransporterTypeID uint64 `json:"transporter_type_id"`
	PrimaryPhone      string `json:"primary_phone"`
	EmailAddress      string `json:"email_address"`
	RouteID           uint64 `json:"route_id"`
	Status            string `gorm:"type:enum('ACTIVE','INACTIVE','SUSPENDED','BLACKLISTED');default:'ACTIVE'" json:"status"`
	Restricted        bool   `gorm:"not null;default:false" json:"restricted"`

	// Relationships
	Individual *IndividualTransporter `gorm:"foreignKey:TransporterID" json:"individual,omitempty"`
	Company    *CompanyTransporter    `gorm:"foreignKey:TransporterID" json:"company,omitempty"`
}

type IndividualTransporter struct {
	BaseModel
	TransporterID     uint64       `gorm:"unique;not null" json:"transporter_id"`
	Transporter       *Transporter `gorm:"foreignKey:TransporterID" json:"transporter,omitempty"`
	FirstName         string       `json:"first_name"`
	LastName          string       `json:"last_name"`
	OtherNames        string       `json:"other_names"`
	Gender            string       `gorm:"type:enum('MALE','FEMALE','OTHER')" json:"gender"`
	DateOfBirth       time.Time    `json:"date_of_birth"`
	NationalIDNo      string       `gorm:"unique" json:"national_id_no"`
	KraPin            string       `json:"kra_pin"`
	MaritalStatus     string       `gorm:"type:enum('SINGLE','MARRIED','DIVORCED','WIDOWED')" json:"marital_status"`
	NextOfKinFullName string       `json:"next_of_kin_full_name"`
	NextOfKinPhone    string       `json:"next_of_kin_phone"`
	PassportPhoto     string       `json:"passport_photo"`
	IDFrontPhoto      string       `json:"id_front_photo"`
	IDBackPhoto       string       `json:"id_back_photo"`
}

type CompanyTransporter struct {
	BaseModel
	TransporterID              uint64       `gorm:"unique;not null" json:"transporter_id"`
	Transporter                *Transporter `gorm:"foreignKey:TransporterID" json:"transporter,omitempty"`
	CompanyName                string       `gorm:"not null" json:"company_name"`
	RegistrationNo             string       `gorm:"unique" json:"registration_no"`
	KraPin                     string       `json:"kra_pin"`
	ContactPersonName          string       `json:"contact_person_name"`
	ContactPersonPhone         string       `json:"contact_person_phone"`
	PostalAddress              string       `json:"postal_address"`
	PostalCode                 string       `json:"postal_code"`
	Town                       string       `json:"town"`
	CertificateOfIncorporation string       `json:"certificate_of_incorporation"`
}

type TransporterVehicle struct {
	BaseModel
	TransporterID  uint64  `json:"transporter_id"`
	RouteID        uint64  `json:"route_id"`
	RegistrationNo string  `gorm:"unique" json:"registration_no"`
	VehicleType    string  `gorm:"type:enum('MOTORBIKE','PICKUP','VAN','TANKER','TRUCK')" json:"vehicle_type"`
	CapacityLitres float64 `gorm:"type:decimal(10,2)" json:"capacity_litres"`
	Active         bool    `gorm:"default:true" json:"active"`
}

type TransporterRouteAssignment struct {
	BaseModel
	TransporterID uint64    `json:"transporter_id"`
	RouteID       uint64    `json:"route_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Active        bool      `gorm:"default:true" json:"active"`
}

type TransporterDriver struct {
	BaseModel
	TransporterID            uint64    `gorm:"not null;index" json:"transporter_id"`
	DriverNo                 string    `gorm:"unique;not null" json:"driver_no"`
	FirstName                string    `gorm:"not null" json:"first_name"`
	LastName                 string    `gorm:"not null" json:"last_name"`
	OtherNames               string    `json:"other_names"`
	Gender                   string    `gorm:"type:enum('MALE','FEMALE','OTHER')" json:"gender"`
	DateOfBirth              time.Time `json:"date_of_birth"`
	NationalIDNo             string    `gorm:"unique;not null" json:"national_id_no"`
	KraPin                   string    `json:"kra_pin"`
	PrimaryPhone             string    `gorm:"not null" json:"primary_phone"`
	SecondaryPhone           string    `json:"secondary_phone"`
	EmailAddress             string    `json:"email_address"`
	DrivingLicenseNo         string    `gorm:"unique;not null" json:"driving_license_no"`
	DrivingLicenseClass      string    `json:"driving_license_class"`
	DrivingLicenseExpiry     time.Time `gorm:"index" json:"driving_license_expiry"`
	EmploymentDate           time.Time `json:"employment_date"`
	Status                   string    `gorm:"type:enum('ACTIVE','INACTIVE','SUSPENDED','TERMINATED');default:'ACTIVE';index" json:"status"`
	PassportPhoto            string    `json:"passport_photo"`
	DrivingLicenseFrontPhoto string    `json:"driving_license_front_photo"`
	DrivingLicenseBackPhoto  string    `json:"driving_license_back_photo"`
	NextOfKinFullName        string    `json:"next_of_kin_full_name"`
	NextOfKinPhone           string    `json:"next_of_kin_phone"`
}

type TransporterDriverAssignment struct {
	BaseModel
	TransporterDriverID  uint64     `gorm:"not null;index" json:"transporter_driver_id"`
	TransporterVehicleID uint64     `gorm:"not null;index" json:"transporter_vehicle_id"`
	AssignedFrom         time.Time  `gorm:"not null" json:"assigned_from"`
	AssignedTo           *time.Time `json:"assigned_to"`
	AssignmentType       string     `gorm:"type:enum('PRIMARY','TEMPORARY','RELIEF','EMERGENCY');default:'PRIMARY'" json:"assignment_type"`
	Active               bool       `gorm:"not null;default:true;index" json:"active"`
	Notes                string     `gorm:"type:text" json:"notes"`
}

type TransporterBankAccount struct {
	BaseModel
	TransporterID uint64 `json:"transporter_id"`
	BankID        uint64 `json:"bank_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	CreatedBy     uint64 `json:"created_by"`
	UpdatedBy     uint64 `json:"updated_by"`
}

type TransporterBenefit struct {
	BaseModel
	Name        string    `json:"name"`
	MinQuantity string    `json:"min_quantity"`
	Rate        string    `json:"rate"`
	RouteID     uint64    `json:"route_id"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedBy   uint64    `json:"created_by"`
	UpdatedBy   uint64    `json:"updated_by"`
}
