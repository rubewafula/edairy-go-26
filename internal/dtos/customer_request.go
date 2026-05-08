package dtos

type CreateCustomerRequest struct {
	ClassID       uint64  `json:"class_id" validate:"required"`
	FullNames     string  `json:"full_names" validate:"required,max=255"`
	Phone         string  `json:"phone" validate:"required,max=15"`
	EmailAddress  string  `json:"email_address" validate:"omitempty,email"`
	CustomerNo    string  `json:"customer_no" validate:"required"`
	KraPin        string  `json:"kra_pin"`
	Status        string  `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	CreditLimit   float64 `json:"credit_limit"`
	PostalAddress string  `json:"postal_address"`
	PostalCode    string  `json:"postal_code"`
	PostalTown    string  `json:"postal_town"`
	SiteID        uint64  `json:"site_id"`
	Terms         string  `json:"terms"`
	Rate          float64 `json:"rate"`
}

type UpdateCustomerRequest struct {
	ClassID       uint64  `json:"class_id" validate:"required"`
	FullNames     string  `json:"full_names" validate:"required,max=255"`
	Phone         string  `json:"phone" validate:"required,max=15"`
	EmailAddress  string  `json:"email_address" validate:"omitempty,email"`
	KraPin        string  `json:"kra_pin"`
	Status        string  `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
	CreditLimit   float64 `json:"credit_limit"`
	PostalAddress string  `json:"postal_address"`
	PostalCode    string  `json:"postal_code"`
	PostalTown    string  `json:"postal_town"`
	SiteID        uint64  `json:"site_id"`
	Terms         string  `json:"terms"`
	Rate          float64 `json:"rate"`
}
