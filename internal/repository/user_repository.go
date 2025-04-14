package repository

import (
	"github.com/Hooannn/EventPlatform/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetDetails(qUser entity.User) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload("Roles.Permissions").Where(qUser).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Get(qUser entity.User) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where(qUser).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
