package models

import "time"

// Account represents a user's billing account
type Account struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"uniqueIndex"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Payment represents a payment transaction
type Payment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`   // "deposit" or "withdraw"
	Status    string    `json:"status"` // "success" or "failed"
	CreatedAt time.Time `json:"created_at"`
}
