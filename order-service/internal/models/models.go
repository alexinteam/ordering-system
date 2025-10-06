package models

import "time"

type Order struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	UserID         string    `json:"user_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"` // "pending", "completed", "failed"
	IdempotencyKey string    `json:"idempotency_key" gorm:"uniqueIndex"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type IdempotencyRecord struct {
	Key       string    `json:"key" gorm:"primaryKey"`
	OrderID   string    `json:"order_id"`
	Status    string    `json:"status"`   // "processing", "completed", "failed"
	Response  string    `json:"response"` // JSON response stored as string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
