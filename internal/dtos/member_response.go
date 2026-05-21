package dtos

import "time"

type MemberResponse struct {
	ID                uint64    `json:"id"`
	MemberNo          string    `json:"member_no"`
	MemberTypeID      uint64    `json:"member_type_id"`
	MemberTypeName    string    `json:"member_type_name"` // From member_types table
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	OtherNames        string    `json:"other_names"`
	RouteID           uint64    `json:"route_id"`
	RouteName         string    `json:"route_name"` // From routes table
	DateOfBirth       string    `json:"date_of_birth"`
	IDNumber          string    `json:"id_no"`
	Gender            string    `json:"gender"`
	BirthCity         string    `json:"birth_city"`
	PrimaryPhone      string    `json:"primary_phone"`
	SecondaryPhone    string    `json:"secondary_phone"`
	Email             string    `json:"email"`
	NumberOfCows      int       `json:"number_of_cows"`
	IdFrontPhoto      string    `json:"id_front_photo"`
	IdBackPhoto       string    `json:"id_back_photo"`
	PassportPhoto     string    `json:"passport_photo"`
	IdDateOfIssue     time.Time `json:"id_date_of_issue"`
	TaxNumber         string    `json:"tax_number"`
	MaritalStatus     string    `json:"marital_status"`
	Title             string    `json:"title"`
	NextOfKinFullName string    `json:"next_of_kin_full_name"`
	NextOfKinPhone    string    `json:"next_of_kin_phone"`
	Status            string    `json:"status"`
	DateRegistered    time.Time `json:"date_registered"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
