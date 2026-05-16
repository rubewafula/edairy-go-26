package dtos

import "time"

type UserResponse struct {
	ID              uint64               `json:"id" gorm:"column:id"`
	Name            string               `json:"name" gorm:"column:name"`
	Email           string               `json:"email" gorm:"column:email"`
	EmailVerifiedAt *string              `json:"email_verified_at" gorm:"column:email_verified_at"`
	IsVerified      bool                 `json:"is_verified" gorm:"column:is_verified"`
	CreatedAt       time.Time            `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time            `json:"updated_at" gorm:"column:updated_at"`
	Roles           []RoleResponse       `json:"roles" gorm:"-"`
	Permissions     []PermissionResponse `json:"permissions" gorm:"-"`
}
