package database

import (
	"log"

	"billing-service/internal/config"
	"billing-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Billing service connected to database successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Payment{}, &models.Account{}); err != nil {
		return err
	}
	return nil
}
