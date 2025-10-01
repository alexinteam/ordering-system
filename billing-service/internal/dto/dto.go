package dto

type ProcessPaymentRequest struct {
	OrderID       string  `json:"order_id" binding:"required"`
	UserID        string  `json:"user_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
}

type ProcessPaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type RefundPaymentRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
	Reason    string `json:"reason"`
}

type RefundPaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CreateAccountRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
	Message   string `json:"message"`
}

type WithdrawRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type WithdrawResponse struct {
	Success    bool    `json:"success"`
	NewBalance float64 `json:"new_balance"`
	Message    string  `json:"message"`
}
