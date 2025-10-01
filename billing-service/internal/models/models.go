package models

import (
	"time"
)

type Payment struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	OrderID       string    `json:"order_id" gorm:"not null"`
	UserID        string    `json:"user_id" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null"` // "pending", "completed", "failed", "refunded"
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Account struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;unique"`
	Balance   float64   `json:"balance" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
