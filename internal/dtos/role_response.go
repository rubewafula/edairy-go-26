package dtos

import "time"

type RoleResponse struct {
	ID          uint64               `json:"id" gorm:"column:id"`
	Name        string               `json:"name" gorm:"column:name"`
	GuardName   string               `json:"guard_name" gorm:"column:guard_name"`
	Permissions []PermissionResponse `json:"permissions" gorm:"-"`
	CreatedAt   time.Time            `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time            `json:"updated_at" gorm:"column:updated_at"`
}
