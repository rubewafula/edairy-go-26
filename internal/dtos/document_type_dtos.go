package dtos

import "time"

type CreateDocumentTypeRequest struct {
	DocumentType string `json:"document_type" validate:"required"`
	Description  string `json:"description"`
}

type UpdateDocumentTypeRequest struct {
	DocumentType string `json:"document_type"`
	Description  string `json:"description"`
}

type DocumentTypeResponse struct {
	ID           uint64    `json:"id"`
	DocumentType string    `json:"document_type"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uint64    `json:"created_by,omitempty"`
	UpdatedBy    uint64    `json:"updated_by,omitempty"`
}
