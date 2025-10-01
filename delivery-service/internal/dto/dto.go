package dto

type ReserveDeliveryRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	CourierID string `json:"courier_id" binding:"required"`
	SlotID    string `json:"slot_id" binding:"required"`
	Address   string `json:"address" binding:"required"`
}

type ReserveDeliveryResponse struct {
	ReservationID string `json:"reservation_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type CancelDeliveryRequest struct {
	ReservationID string `json:"reservation_id" binding:"required"`
	Reason        string `json:"reason"`
}

type CancelDeliveryResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
