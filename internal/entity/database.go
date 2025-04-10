package entity

import (
	"github.com/Hooannn/EventPlatform/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	cfg := configs.LoadConfig(".env")
	db, err := gorm.Open(postgres.Open(cfg.DatabaseConnectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Role{}, &Permission{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
