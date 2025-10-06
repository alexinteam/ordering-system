package dto

type CreateOrderRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Email          string  `json:"email" binding:"required,email"`
	ProductID      string  `json:"product_id" binding:"required"`
	Quantity       int     `json:"quantity" binding:"required,gt=0"`
	Address        string  `json:"address" binding:"required"`
	IdempotencyKey string  `json:"idempotency_key" binding:"required"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
