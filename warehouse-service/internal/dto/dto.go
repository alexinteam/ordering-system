package dto

type ReserveProductRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type ReserveProductResponse struct {
	ReservationID string `json:"reservation_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type CancelReservationRequest struct {
	ReservationID string `json:"reservation_id" binding:"required"`
	Reason        string `json:"reason"`
}

type CancelReservationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
