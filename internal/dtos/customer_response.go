package dtos

import "time"

type CustomerResponse struct {
	ID               uint64    `json:"id"`
	CustomerTypeID   uint64    `json:"customer_type_id"`
	CustomerTypeName string    `json:"customer_type_name"`
	FullNames        string    `json:"full_names"`
	Phone            string    `json:"phone"`
	EmailAddress     string    `json:"email_address"`
	CustomerNo       string    `json:"customer_no"`
	KraPin           string    `json:"kra_pin"`
	Status           string    `json:"status"`
	CreditLimit      float64   `json:"credit_limit"`
	PostalAddress    string    `json:"postal_address"`
	PostalCode       string    `json:"postal_code"`
	PostalTown       string    `json:"postal_town"`
	SiteID           uint64    `json:"site_id"`
	Terms            string    `json:"terms"`
	Rate             float64   `json:"rate"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
