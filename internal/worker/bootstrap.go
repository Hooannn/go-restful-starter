package worker

import (
	"errors"
	"log"
	"strings"

	"github.com/Hooannn/go-restful-starter/configs"
	"github.com/Hooannn/go-restful-starter/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Bootstrap(cfg *configs.Config, db *gorm.DB) {
	log.Println("Bootstrapping...", cfg.AppName)

	var rUser entity.User

	result := db.Model(&entity.User{}).Take(&rUser)

	if result.RowsAffected > 0 {
		log.Println("✅ Root user exists, done bootstrapping")
		return
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("❌ Unexpected error, failed to bootstrap", result.Error)
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		permissions := []entity.Permission{}

		for _, permission := range strings.Split(cfg.RootRolePermissions, ",") {
			permissions = append(permissions, entity.Permission{Name: permission})
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(cfg.RootUserPassword), bcrypt.DefaultCost)
		if err := tx.Create(&entity.User{
			FirstName: "Root",
			LastName:  "User",
			Email:     cfg.RootUser,
			Password:  string(hash),
			Roles: []entity.Role{
				{
					Name:        cfg.RootRoleName,
					Description: &cfg.RootRoleDescription,
					Permissions: permissions,
				},
			},
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("❌ Unexpected error, failed to bootstrap", err)
	} else {
		log.Println("✅ Done bootstrapping")
	}
}
