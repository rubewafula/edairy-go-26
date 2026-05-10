package dtos

type CreateTransporterRequest struct {
	TransporterNo     string `json:"transporter_no" validate:"required"`
	Category          string `json:"transporter_category" validate:"required,oneof=INDIVIDUAL COMPANY"`
	TransporterTypeID uint64 `json:"transporter_type_id"`
	PrimaryPhone      string `json:"primary_phone" validate:"required,max=15"`
	EmailAddress      string `json:"email_address" validate:"omitempty,email"`
	RouteID           uint64 `json:"route_id"`
	Status            string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED BLACKLISTED"`
	Restricted        bool   `json:"restricted"`

	// Individual Transporter Fields
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	OtherNames        string `json:"other_names"`
	Gender            string `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	DateOfBirth       string `json:"date_of_birth"`
	NationalIDNo      string `json:"national_id_no"`
	KraPin            string `json:"kra_pin"`
	MaritalStatus     string `json:"marital_status" validate:"omitempty,oneof=SINGLE MARRIED DIVORCED WIDOWED"`
	NextOfKinFullName string `json:"next_of_kin_full_name"`
	NextOfKinPhone    string `json:"next_of_kin_phone"`
	PassportPhoto     string `json:"passport_photo"`
	IDFrontPhoto      string `json:"id_front_photo"`
	IDBackPhoto       string `json:"id_back_photo"`

	// Company Transporter Fields
	CompanyName                string `json:"company_name"`
	RegistrationNo             string `json:"registration_no"`
	ContactPersonName          string `json:"contact_person_name"`
	ContactPersonPhone         string `json:"contact_person_phone"`
	PostalAddress              string `json:"postal_address"`
	PostalCode                 string `json:"postal_code"`
	Town                       string `json:"town"`
	CertificateOfIncorporation string `json:"certificate_of_incorporation"`
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
