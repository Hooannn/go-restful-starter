package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `json:"name" gorm:"unique;not null"`
	Description *string      `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

func (Role) TableName() string {
	return "roles"
}
