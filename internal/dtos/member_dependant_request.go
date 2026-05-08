package dtos

type CreateMemberDependantRequest struct {
	MemberID     uint64 `json:"member_id" validate:"required"`
	Name         string `json:"name" validate:"required,max=255"`
	NationalID   string `json:"national_id"`
	Relationship string `json:"relationship" validate:"required"`
	MobileNo     string `json:"mobile_no"`
	Gender       string `json:"gender" validate:"required,oneof=MALE FEMALE"`
	DateOfBirth  string `json:"date_of_birth" validate:"required,datetime"`
	BirthCertNo  string `json:"birth_cert_no"`
	Email        string `json:"email" validate:"omitempty,email"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
	Town         string `json:"town"`
}

type UpdateMemberDependantRequest struct {
	MemberID     uint64 `json:"member_id" validate:"required"`
	Name         string `json:"name" validate:"required,max=255"`
	NationalID   string `json:"national_id"`
	Relationship string `json:"relationship" validate:"required"`
	MobileNo     string `json:"mobile_no"`
	Gender       string `json:"gender" validate:"required,oneof=MALE FEMALE"`
	DateOfBirth  string `json:"date_of_birth" validate:"required,datetime"`
	BirthCertNo  string `json:"birth_cert_no"`
	Email        string `json:"email" validate:"omitempty,email"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
	Town         string `json:"town"`
}
