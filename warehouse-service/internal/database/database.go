package database

import (
	"fmt"
	"log"

	"warehouse-service/internal/config"
	"warehouse-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Warehouse service connected to database successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Product{}, &models.Reservation{}); err != nil {
		return err
	}
	return nil
}

func CreateTestProducts(db *gorm.DB) error {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		return nil
	}

	products := []models.Product{
		{
			ID:          "prod_1",
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       1500.00,
			Stock:       10,
			Reserved:    0,
		},
		{
			ID:          "prod_2",
			Name:        "Smartphone",
			Description: "Latest smartphone model",
			Price:       800.00,
			Stock:       25,
			Reserved:    0,
		},
		{
			ID:          "prod_3",
			Name:        "Headphones",
			Description: "Wireless noise-cancelling headphones",
			Price:       200.00,
			Stock:       50,
			Reserved:    0,
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			return fmt.Errorf("failed to create test product %s: %v", product.ID, err)
		}
	}

	log.Println("Test products created successfully")
	return nil
}
