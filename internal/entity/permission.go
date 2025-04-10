package entity

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string  `json:"name" gorm:"unique;not null"`
	Description *string `json:"description"`
}

func (Permission) TableName() string {
	return "permissions"
}
