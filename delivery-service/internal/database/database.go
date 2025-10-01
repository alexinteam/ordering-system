package database

import (
	"fmt"
	"log"
	"time"

	"delivery-service/internal/config"
	"delivery-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Delivery service connected to database successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Courier{}, &models.DeliverySlot{}, &models.DeliveryReservation{}); err != nil {
		return err
	}
	return nil
}

func CreateTestData(db *gorm.DB) error {
	// Проверяем, есть ли уже данные
	var courierCount int64
	db.Model(&models.Courier{}).Count(&courierCount)
	if courierCount > 0 {
		return nil
	}

	// Создаем курьеров
	couriers := []models.Courier{
		{
			ID:        "courier_1",
			Name:      "Иван Петров",
			Phone:     "+7-900-123-45-67",
			Status:    "available",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		{
			ID:        "courier_2",
			Name:      "Мария Сидорова",
			Phone:     "+7-900-234-56-78",
			Status:    "available",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		{
			ID:        "courier_3",
			Name:      "Алексей Козлов",
			Phone:     "+7-900-345-67-89",
			Status:    "available",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}

	for _, courier := range couriers {
		if err := db.Create(&courier).Error; err != nil {
			return fmt.Errorf("failed to create courier %s: %v", courier.ID, err)
		}
	}

	// Создаем слоты доставки на сегодня и завтра
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	timeSlots := []struct {
		start string
		end   string
	}{
		{"09:00", "12:00"},
		{"12:00", "15:00"},
		{"15:00", "18:00"},
		{"18:00", "21:00"},
	}

	for _, courier := range couriers {
		for _, date := range []string{today, tomorrow} {
			for _, timeSlot := range timeSlots {
				slot := models.DeliverySlot{
					ID:        fmt.Sprintf("slot_%s_%s_%s", courier.ID, date, timeSlot.start),
					CourierID: courier.ID,
					Date:      date,
					StartTime: timeSlot.start,
					EndTime:   timeSlot.end,
					Status:    "available",
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				}
				if err := db.Create(&slot).Error; err != nil {
					return fmt.Errorf("failed to create slot %s: %v", slot.ID, err)
				}
			}
		}
	}

	log.Println("Test data created successfully")
	return nil
}
