package dtos

import "time"

type RoleResponse struct {
	ID          uint64               `json:"ID" gorm:"column:id"`
	Name        string               `json:"Name" gorm:"column:name"`
	GuardName   string               `json:"GuardName" gorm:"column:guard_name"`
	Permissions []PermissionResponse `json:"Permissions" gorm:"-"`
	CreatedAt   time.Time            `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt   time.Time            `json:"UpdatedAt" gorm:"column:updated_at"`
}
