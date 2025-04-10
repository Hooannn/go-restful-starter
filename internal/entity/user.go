package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Password     string         `json:"-"`
	Phone        *string        `json:"phone"`
	Email        string         `json:"email"`
	Birthday     *time.Time     `json:"birthday"`
	MemberNumber *string        `json:"member_number"`
	ActivatedAt  *time.Time     `json:"activated_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Roles        []Role         `json:"roles" gorm:"many2many:user_roles;"`
}

func (User) TableName() string {
	return "users"
}
