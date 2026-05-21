package models

import "time"

type Permission struct {
	BaseModel
	Name      string `gorm:"column:name" json:"name"`
	GuardName string `gorm:"column:guard_name" json:"guard_name"`
}

type Role struct {
	BaseModel
	Name        string       `gorm:"column:name" json:"name"`
	GuardName   string       `gorm:"column:guard_name" json:"guard_name"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;" json:"permissions"`
}

type UserRole struct {
	BaseModel
	UserID uint64 `gorm:"primaryKey" json:"user_id"`
	RoleID uint64 `gorm:"primaryKey" json:"role_id"`
}

type UserPermission struct {
	BaseModel
	UserID       uint64 `gorm:"primaryKey" json:"user_id"`
	PermissionID uint64 `gorm:"primaryKey" json:"permission_id"`
}

type RolePermission struct {
	BaseModel
	RoleID       uint64 `gorm:"primaryKey" json:"role_id"`
	PermissionID uint64 `gorm:"primaryKey" json:"permission_id"`
}

type User struct {
	BaseModel
	Name              string     `gorm:"column:name" json:"name"`
	Email             string     `gorm:"column:email;uniqueIndex" json:"email"`
	EmailVerifiedAt   *string    `gorm:"column:email_verified_at" json:"email_verified_at"`
	Password          string     `gorm:"column:password" json:"password"`
	RememberToken     string     `gorm:"column:remember_token" json:"remember_token"`
	IsVerified        bool       `gorm:"column:is_verified;default:0" json:"is_verified"`
	VerificationToken string     `gorm:"column:verification_token" json:"verification_token"`
	ResetToken        string     `gorm:"column:reset_token" json:"reset_token"`
	ResetTokenExpiry  *time.Time `gorm:"column:reset_token_expiry" json:"reset_token_expiry"`

	Roles       []Role       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;" json:"roles"`
	Permissions []Permission `gorm:"many2many:user_permissions;constraint:OnDelete:CASCADE;" json:"permissions"`
}
