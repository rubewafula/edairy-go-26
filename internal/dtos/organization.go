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
	ID          uint64    `json:"id"`
	AddressType string    `json:"address_type"`
	City        string    `json:"city"`
	Code        string    `json:"code"`
	Country     string    `json:"country"`
	Line1       string    `json:"line1"`
	Line2       string    `json:"line2"`
	Line3       string    `json:"line3"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrganizationDocument DTOs
type CreateOrganizationDocumentRequest struct {
	AstraID               uint64 `json:"astra_id"`
	DocumentTypeID        uint64 `json:"document_type_id" validate:"required"`        // This is the ID of the document type
	DocumentName          string `json:"document_name" validate:"required"`           // Original file name
	DocumentContentBase64 string `json:"document_content_base64" validate:"required"` // Base64 encoded file content
	Submitted             bool   `json:"submitted"`
}

type UpdateOrganizationDocumentRequest struct {
	AstraID               uint64 `json:"astra_id"`
	DocumentTypeID        uint64 `json:"document_type_id"` // Corrected typo
	DocumentName          string `json:"document_name"`
	DocumentContentBase64 string `json:"document_content_base64"`
	Submitted             bool   `json:"submitted"`
}

type OrganizationDocumentResponse struct {
	ID               uint64    `json:"id"`
	AstraID          uint64    `json:"astra_id"`
	DocumentTypeID   uint64    `json:"document_type_id"`   // Changed to ID
	DocumentTypeName string    `json:"document_type_name"` // Added for display
	DocumentURL      string    `json:"document_url"`       // The URL to the stored document
	Submitted        bool      `json:"submitted"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
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
	ID        uint64    `json:"id"`
	Issue     string    `json:"issue"`
	Comment   string    `json:"comment"`
	Iteration int       `json:"iteration"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	ID             uint64    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	IDNo           string    `json:"id_no"`
	Position       string    `json:"position"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	LinkStatus     string    `json:"link_status"`
	LivenessPassed bool      `json:"liveness_passed"`
	Submitted      bool      `json:"submitted"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
	ID             uint64    `json:"id"`
	WalletTypeID   uint64    `json:"wallet_type_id"`
	WalletTypeName string    `json:"wallet_type_name"`
	WalletID       string    `json:"wallet_id"`
	WalletName     string    `json:"wallet_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
