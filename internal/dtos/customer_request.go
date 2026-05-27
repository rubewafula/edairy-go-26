package dtos

type CreateCustomerRequest struct {
	CustomerTypeID uint64  `json:"customer_type_id" validate:"required"`
	FullNames      string  `json:"full_names" validate:"required,max=255"`
	Phone          string  `json:"phone" validate:"required,max=15"`
	EmailAddress   string  `json:"email_address" validate:"omitempty,email"`
	CustomerNo     string  `json:"customer_no" validate:"required"`
	KraPin         string  `json:"kra_pin"`
	Status         string  `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	CreditLimit    float64 `json:"credit_limit"`
	PostalAddress  string  `json:"postal_address"`
	PostalCode     string  `json:"postal_code"`
	PostalTown     string  `json:"postal_town"`
	Terms          uint64  `json:"terms"`
	Rate           float64 `json:"rate"`
}

type UpdateCustomerRequest struct {
	CustomerTypeID uint64  `json:"customer_type_id" validate:"required"`
	FullNames      string  `json:"full_names" validate:"required,max=255"`
	Phone          string  `json:"phone" validate:"required,max=15"`
	EmailAddress   string  `json:"email_address" validate:"omitempty,email"`
	KraPin         string  `json:"kra_pin"`
	Status         string  `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	CreditLimit    float64 `json:"credit_limit"`
	PostalAddress  string  `json:"postal_address"`
	PostalCode     string  `json:"postal_code"`
	PostalTown     string  `json:"postal_town"`
	Terms          uint64  `json:"terms"`
	Rate           float64 `json:"rate"`
}
