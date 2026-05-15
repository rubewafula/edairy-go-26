package dtos

import "time"

// OrganizationAddress DTOs
type CreateOrganizationAddressRequest struct {
	AddressType string `json:"address_type" validate:"required"`
	City        string `json:"city"`
	Code        string `json:"code"`
	Country     string `json:"country"`
	Line1       string `json:"line1" validate:"required"`
	Line2       string `json:"line2"`
	Line3       string `json:"line3"`
	State       string `json:"state"`
}

type UpdateOrganizationAddressRequest struct {
	AddressType string `json:"address_type"`
	City        string `json:"city"`
	Code        string `json:"code"`
	Country     string `json:"country"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
	Line3       string `json:"line3"`
	State       string `json:"state"`
}

type OrganizationAddressResponse struct {
	ID          uint64    `json:"ID"`
	AddressType string    `json:"AddressType"`
	City        string    `json:"City"`
	Code        string    `json:"Code"`
	Country     string    `json:"Country"`
	Line1       string    `json:"Line1"`
	Line2       string    `json:"Line2"`
	Line3       string    `json:"Line3"`
	State       string    `json:"State"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

// OrganizationBank DTOs
type CreateOrganizationBankRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateOrganizationBankRequest struct {
	Name string `json:"name" validate:"required"`
}

type OrganizationBankResponse struct {
	ID        uint64    `json:"ID"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// OrganizationDocument DTOs
type CreateOrganizationDocumentRequest struct {
	AstraID      uint64 `json:"astra_id"`
	DocumentType string `json:"document_type" validate:"required"`
	Document     string `json:"document" validate:"required"` // Base64 encoded or file path
	Submitted    bool   `json:"submitted"`
}

type UpdateOrganizationDocumentRequest struct {
	AstraID      uint64 `json:"astra_id"`
	DocumentType string `json:"document_type"`
	Document     string `json:"document"`
	Submitted    bool   `json:"submitted"`
}

type OrganizationDocumentResponse struct {
	ID           uint64    `json:"ID"`
	AstraID      uint64    `json:"AstraID"`
	DocumentType string    `json:"DocumentType"`
	Document     string    `json:"Document"`
	Submitted    bool      `json:"Submitted"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}

// OrganizationKybComment DTOs
type CreateOrganizationKybCommentRequest struct {
	Issue     string `json:"issue" validate:"required"`
	Comment   string `json:"comment" validate:"required"`
	Iteration int    `json:"iteration"`
}

type UpdateOrganizationKybCommentRequest struct {
	Issue     string `json:"issue"`
	Comment   string `json:"comment"`
	Iteration int    `json:"iteration"`
}

type OrganizationKybCommentResponse struct {
	ID        uint64    `json:"ID"`
	Issue     string    `json:"Issue"`
	Comment   string    `json:"Comment"`
	Iteration int       `json:"Iteration"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// OrganizationLeadership DTOs
type CreateOrganizationLeadershipRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	IDNo      string `json:"id_no" validate:"required"`
	Position  string `json:"position" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"email"`
}

type UpdateOrganizationLeadershipRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IDNo      string `json:"id_no"`
	Position  string `json:"position"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type OrganizationLeadershipResponse struct {
	ID             uint64    `json:"ID"`
	FirstName      string    `json:"FirstName"`
	LastName       string    `json:"LastName"`
	IDNo           string    `json:"IDNo"`
	Position       string    `json:"Position"`
	Phone          string    `json:"Phone"`
	Email          string    `json:"Email"`
	Status         string    `json:"Status"`
	LinkStatus     string    `json:"LinkStatus"`
	LivenessPassed bool      `json:"LivenessPassed"`
	Submitted      bool      `json:"Submitted"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}

// OrganizationMember DTOs
type CreateOrganizationMemberRequest struct {
	CustomerID     uint64 `json:"customer_id" validate:"required"`
	CustomerType   string `json:"customer_type" validate:"required"`
	ManuallyRatify bool   `json:"manually_ratify"`
	NextLevel      string `json:"next_level"`
	Status         string `json:"status"`
	AstraID        string `json:"astra_id"`
	CreditLimit    uint64 `json:"credit_limit"`
	LinkStatus     string `json:"link_status"`
	LivenessPassed bool   `json:"liveness_passed"`
	AstraRemarks   string `json:"astra_remarks"`
	UUID           string `json:"uuid"`
	AuthCreated    bool   `json:"auth_created"`
	Locale         string `json:"locale"`
}

type UpdateOrganizationMemberRequest struct {
	ManuallyRatify bool   `json:"manually_ratify"`
	NextLevel      string `json:"next_level"`
	Status         string `json:"status"`
	AstraID        string `json:"astra_id"`
	CreditLimit    uint64 `json:"credit_limit"`
	LinkStatus     string `json:"link_status"`
	LivenessPassed bool   `json:"liveness_passed"`
	AstraRemarks   string `json:"astra_remarks"`
	AuthCreated    bool   `json:"auth_created"`
	Locale         string `json:"locale"`
}

type OrganizationMemberResponse struct {
	ID             uint64    `json:"ID"`
	CustomerID     uint64    `json:"CustomerID"`
	CustomerName   string    `json:"CustomerName"`
	CustomerType   string    `json:"CustomerType"`
	ManuallyRatify bool      `json:"ManuallyRatify"`
	NextLevel      string    `json:"NextLevel"`
	Status         string    `json:"Status"`
	AstraID        string    `json:"AstraID"`
	CreditLimit    uint64    `json:"CreditLimit"`
	LinkStatus     string    `json:"LinkStatus"`
	LivenessPassed bool      `json:"LivenessPassed"`
	AstraRemarks   string    `json:"AstraRemarks"`
	UUID           string    `json:"UUID"`
	AuthCreated    bool      `json:"AuthCreated"`
	Locale         string    `json:"Locale"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}

// OrganizationWallet DTOs
type CreateOrganizationWalletRequest struct {
	WalletTypeID uint64 `json:"wallet_type_id" validate:"required"`
	WalletID     string `json:"wallet_id" validate:"required"`
	WalletName   string `json:"wallet_name" validate:"required"`
}

type UpdateOrganizationWalletRequest struct {
	WalletTypeID uint64 `json:"wallet_type_id"`
	WalletID     string `json:"wallet_id"`
	WalletName   string `json:"wallet_name"`
}

type OrganizationWalletResponse struct {
	ID             uint64    `json:"ID"`
	WalletTypeID   uint64    `json:"WalletTypeID"`
	WalletTypeName string    `json:"WalletTypeName"`
	WalletID       string    `json:"WalletID"`
	WalletName     string    `json:"WalletName"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
}
