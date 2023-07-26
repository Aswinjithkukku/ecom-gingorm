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
	Id       int      `json:"id" gorm:"primary_key"`
	Name     string   `json:"name"`
	Email    string   `json:"email" gorm:"not null" validate:"required,email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role" gorm:"type:enum('SUPERADMIN', 'ADMIN', 'CLIENT')"`
}
