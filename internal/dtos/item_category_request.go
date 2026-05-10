package dtos

import "time"

type CreateItemCategoryRequest struct {
	Name             string `json:"Name" validate:"required,max=255"`
	Description      string `json:"Description" validate:"required,max=255"`
	ParentCategoryID uint64 `json:"ParentCategoryID"`
}

type UpdateItemCategoryRequest struct {
	Name             string `json:"Name" validate:"required,max=255"`
	Description      string `json:"Description" validate:"required,max=255"`
	ParentCategoryID uint64 `json:"ParentCategoryID"`
}

type ItemCategoryResponse struct {
	ID                 uint64    `json:"ID"`
	Name               string    `json:"Name"`
	Description        string    `json:"Description"`
	ParentCategoryID   uint64    `json:"ParentCategoryID"`
	ParentCategoryName string    `json:"ParentCategoryName"`
	CreatedAt          time.Time `json:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt"`
}
