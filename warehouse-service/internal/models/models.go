package models

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	Reserved    int       `json:"reserved" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Reservation struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	OrderID   string    `json:"order_id" gorm:"not null"`
	ProductID string    `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"` // "reserved", "confirmed", "cancelled"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
