package repository

import (
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/google/uuid"
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

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Model(entity.User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindById(id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Model(entity.User{ID: uuid.MustParse(id)}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
