package dtos

type CreateTransporterDriverRequest struct {
	TransporterID            uint64 `json:"transporter_id" validate:"required"`
	DriverNo                 string `json:"driver_no" validate:"required"`
	FirstName                string `json:"first_name" validate:"required"`
	LastName                 string `json:"last_name" validate:"required"`
	OtherNames               string `json:"other_names"`
	Gender                   string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	DateOfBirth              string `json:"date_of_birth"`
	NationalIDNo             string `json:"national_id_no" validate:"required"`
	KraPin                   string `json:"kra_pin"`
	PrimaryPhone             string `json:"primary_phone" validate:"required"`
	SecondaryPhone           string `json:"secondary_phone"`
	EmailAddress             string `json:"email_address" validate:"omitempty,email"`
	DrivingLicenseNo         string `json:"driving_license_no" validate:"required"`
	DrivingLicenseClass      string `json:"driving_license_class"`
	DrivingLicenseExpiry     string `json:"driving_license_expiry"`
	EmploymentDate           string `json:"employment_date"`
	Status                   string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED TERMINATED"`
	NextOfKinFullName        string `json:"next_of_kin_full_name"`
	NextOfKinPhone           string `json:"next_of_kin_phone"`
	PassportPhoto            string `json:"passport_photo"`
	DrivingLicenseFrontPhoto string `json:"driving_license_front_photo"`
	DrivingLicenseBackPhoto  string `json:"driving_license_back_photo"`
}

type UpdateTransporterDriverRequest struct {
	FirstName                string `json:"first_name" validate:"required"`
	LastName                 string `json:"last_name" validate:"required"`
	OtherNames               string `json:"other_names"`
	Gender                   string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	DateOfBirth              string `json:"date_of_birth"`
	NationalIDNo             string `json:"national_id_no" validate:"required"`
	KraPin                   string `json:"kra_pin"`
	PrimaryPhone             string `json:"primary_phone" validate:"required"`
	SecondaryPhone           string `json:"secondary_phone"`
	EmailAddress             string `json:"email_address" validate:"omitempty,email"`
	DrivingLicenseNo         string `json:"driving_license_no" validate:"required"`
	DrivingLicenseExpiry     string `json:"driving_license_expiry"`
	Status                   string `json:"status" validate:"required,oneof=ACTIVE INACTIVE SUSPENDED TERMINATED"`
	NextOfKinFullName        string `json:"next_of_kin_full_name"`
	NextOfKinPhone           string `json:"next_of_kin_phone"`
	PassportPhoto            string `json:"passport_photo"`
	DrivingLicenseFrontPhoto string `json:"driving_license_front_photo"`
	DrivingLicenseBackPhoto  string `json:"driving_license_back_photo"`
}
