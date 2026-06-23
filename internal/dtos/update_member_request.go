package dtos

type UpdateMemberRequest struct {
	MemberTypeID uint64 `json:"member_type_id" validate:"required"`
	FirstName    string `json:"first_name" validate:"required,max=128"`
	LastName     string `json:"last_name" validate:"required,max=128"`
	OtherNames   string `json:"other_names" validate:"max=128"`
	RouteID      uint64 `json:"route_id"`
	DOB          string `json:"dob"`

	IDNo      string `json:"id_no" validate:"max=25"`
	MemberNo  string `json:"member_no"`
	BirthCity string `json:"birth_city"`
	Gender    string `json:"gender" validate:"required,oneof=MALE FEMALE"`

	PrimaryPhone   string `json:"primary_phone" validate:"max=15"`
	SecondaryPhone string `json:"secondary_phone"`
	Email          string `json:"email" validate:"omitempty,email"`
	NumberOfCows   int    `json:"number_of_cows" validate:"min=0"`

	IDFrontPhoto  string                         `json:"id_front_photo" validate:"max=255"`
	IDBackPhoto   string                         `json:"id_back_photo" validate:"max=255"`
	PassportPhoto string                         `json:"passport_photo" validate:"max=255"`
	IDDateOfIssue string                         `json:"id_date_of_issue" validate:"omitempty"`
	NextOfKins    []CreateMemberNextOfKinRequest `json:"next_of_kins" form:"next_of_kins"`

	TaxNumber     string `json:"tax_number"`
	MaritalStatus string `json:"marital_status"`
	Title         string `json:"title"`

	NextOfKinFullName string `json:"next_of_kin_full_name"`
	NextOfKinPhone    string `json:"next_of_kin_phone"`

	BankID      uint64 `json:"bank_id" form:"bank_id"`
	BankBranch  string `json:"bank_branch" form:"bank_branch"`
	AccountNo   string `json:"account_no" form:"account_no"`
	AccountName string `json:"account_name" form:"account_name"`
}
