package dtos

import "time"

// CreateMemberNextOfKinRequest defines the request body for creating a new member next of kin.
type CreateMemberNextOfKinRequest struct {
	MemberID               uint64  `json:"member_id" validate:"required"`
	FullName               string  `json:"full_name" validate:"required,max=255"`
	Relationship           *string `json:"relationship"`
	PhoneNumber            *string `json:"phone_number"`
	AlternativePhoneNumber *string `json:"alternative_phone_number"`
	EmailAddress           *string `json:"email_address" validate:"omitempty,email"`
	NationalIDNo           *string `json:"national_id_no"`
	PostalAddress          *string `json:"postal_address"`
	PhysicalAddress        *string `json:"physical_address"`
	Occupation             *string `json:"occupation"`
	IsPrimary              bool    `json:"is_primary"`
	Status                 bool    `json:"status"`
	Remarks                *string `json:"remarks"`
}

// UpdateMemberNextOfKinRequest defines the request body for updating an existing member next of kin.
type UpdateMemberNextOfKinRequest struct {
	MemberID               uint64  `json:"member_id"`
	FullName               string  `json:"full_name"`
	Relationship           *string `json:"relationship"`
	PhoneNumber            *string `json:"phone_number"`
	AlternativePhoneNumber *string `json:"alternative_phone_number"`
	EmailAddress           *string `json:"email_address" validate:"omitempty,email"`
	NationalIDNo           *string `json:"national_id_no"`
	PostalAddress          *string `json:"postal_address"`
	PhysicalAddress        *string `json:"physical_address"`
	Occupation             *string `json:"occupation"`
	IsPrimary              bool    `json:"is_primary"`
	Status                 bool    `json:"status"`
	Remarks                *string `json:"remarks"`
}

// MemberNextOfKinResponse defines the structure for a member next of kin response.
type MemberNextOfKinResponse struct {
	ID                     uint64    `json:"id"`
	MemberID               uint64    `json:"member_id"`
	MemberNo               string    `json:"member_no"`        // Joined from member_registrations
	MemberFullName         string    `json:"member_full_name"` // Joined from member_registrations
	FullName               string    `json:"full_name"`
	Relationship           *string   `json:"relationship"`
	PhoneNumber            *string   `json:"phone_number"`
	AlternativePhoneNumber *string   `json:"alternative_phone_number"`
	EmailAddress           *string   `json:"email_address"`
	NationalIDNo           *string   `json:"national_id_no"`
	PostalAddress          *string   `json:"postal_address"`
	PhysicalAddress        *string   `json:"physical_address"`
	Occupation             *string   `json:"occupation"`
	IsPrimary              bool      `json:"is_primary"`
	Status                 bool      `json:"status"`
	Remarks                *string   `json:"remarks"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	CreatedBy              uint64    `json:"created_by,omitempty"`
	UpdatedBy              uint64    `json:"updated_by,omitempty"`
}
