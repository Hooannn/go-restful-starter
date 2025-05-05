package repository

import (
	"github.com/Hooannn/go-restful-starter/internal/entity"
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
	result := r.db.Preload("Roles.Permissions").Where(qUser).Limit(1).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) Get(qUser entity.User) (*entity.User, error) {
	var user entity.User
	result := r.db.Where(qUser).Limit(1).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetAll() (*[]entity.User, error) {
	var users []entity.User
	result := r.db.Find(&users)
	if err := result.Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *UserRepository) ExistsByUsername(username string) bool {
	var result byte
	r.db.Raw("SELECT 1 FROM \"users\" WHERE \"users\".\"email\" = ? LIMIT 1", username).Scan(&result)
	return result == 1
}

func (r *UserRepository) UpdatePasswordByUsername(username, password string) bool {
	result := r.db.Model(&entity.User{}).Where("email = ?", username).Update("password", password)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}
