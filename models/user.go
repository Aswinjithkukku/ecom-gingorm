package models

import "gorm.io/gorm"

type UserRole string

const (
	SuperAdmin UserRole = "SUPERADMIN"
	Admin      UserRole = "ADMIN"
	Client     UserRole = "CLIENT"
)

type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"not null" validate:"required,email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
