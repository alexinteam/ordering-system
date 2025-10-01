package models

type Courier struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`
	Phone     string `json:"phone" gorm:"not null"`
	Status    string `json:"status" gorm:"not null"` // "available", "busy", "offline"
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type DeliverySlot struct {
	ID        string `json:"id" gorm:"primaryKey"`
	CourierID string `json:"courier_id" gorm:"not null"`
	Date      string `json:"date" gorm:"not null"`       // YYYY-MM-DD
	StartTime string `json:"start_time" gorm:"not null"` // HH:MM
	EndTime   string `json:"end_time" gorm:"not null"`   // HH:MM
	Status    string `json:"status" gorm:"not null"`     // "available", "reserved", "completed"
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type DeliveryReservation struct {
	ID        string `json:"id" gorm:"primaryKey"`
	OrderID   string `json:"order_id" gorm:"not null"`
	CourierID string `json:"courier_id" gorm:"not null"`
	SlotID    string `json:"slot_id" gorm:"not null"`
	Address   string `json:"address" gorm:"not null"`
	Status    string `json:"status" gorm:"not null"` // "reserved", "confirmed", "cancelled", "completed"
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
