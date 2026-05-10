package dtos

import "time"

type UserResponse struct {
	ID              uint64               `json:"ID" gorm:"column:id"`
	Name            string               `json:"Name" gorm:"column:name"`
	Email           string               `json:"Email" gorm:"column:email"`
	EmailVerifiedAt *string              `json:"EmailVerifiedAt" gorm:"column:email_verified_at"`
	IsVerified      bool                 `json:"IsVerified" gorm:"column:is_verified"`
	CreatedAt       time.Time            `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt       time.Time            `json:"UpdatedAt" gorm:"column:updated_at"`
	Roles           []RoleResponse       `json:"Roles" gorm:"-"`
	Permissions     []PermissionResponse `json:"Permissions" gorm:"-"`
}
