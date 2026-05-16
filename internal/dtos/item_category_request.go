package dtos

import "time"

type CreateItemCategoryRequest struct {
	Name             string `json:"name" validate:"required,max=255"`
	Description      string `json:"description" validate:"required,max=255"`
	ParentCategoryID uint64 `json:"parent_category_id"`
}

type UpdateItemCategoryRequest struct {
	Name             string `json:"name" validate:"required,max=255"`
	Description      string `json:"description" validate:"required,max=255"`
	ParentCategoryID uint64 `json:"parent_category_id"`
}

type ItemCategoryResponse struct {
	ID                 uint64    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ParentCategoryID   uint64    `json:"parent_category_id"`
	ParentCategoryName string    `json:"parent_category_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
