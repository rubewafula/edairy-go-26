package models

import "time"

type Permission struct {
	BaseModel
	Name      string `gorm:"column:name"`
	GuardName string `gorm:"column:guard_name"`
}

type Role struct {
	BaseModel
	Name        string       `gorm:"column:name"`
	GuardName   string       `gorm:"column:guard_name"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;"`
}

type UserRole struct {
	BaseModel
	UserID uint64 `gorm:"primaryKey"`
	RoleID uint64 `gorm:"primaryKey"`
}

type UserPermission struct {
	BaseModel
	UserID       uint64 `gorm:"primaryKey"`
	PermissionID uint64 `gorm:"primaryKey"`
}

type RolePermission struct {
	BaseModel
	RoleID       uint64 `gorm:"primaryKey"`
	PermissionID uint64 `gorm:"primaryKey"`
}

type User struct {
	BaseModel
	Name              string     `gorm:"column:name"`
	Email             string     `gorm:"column:email;uniqueIndex"`
	EmailVerifiedAt   *string    `gorm:"column:email_verified_at"`
	Password          string     `gorm:"column:password"`
	RememberToken     string     `gorm:"column:remember_token"`
	IsVerified        bool       `gorm:"column:is_verified;default:0"`
	VerificationToken string     `gorm:"column:verification_token"`
	ResetToken        string     `gorm:"column:reset_token"`
	ResetTokenExpiry  *time.Time `gorm:"column:reset_token_expiry"`

	Roles       []Role       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Permissions []Permission `gorm:"many2many:user_permissions;constraint:OnDelete:CASCADE;"`
}
