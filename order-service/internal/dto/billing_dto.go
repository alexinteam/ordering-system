package dto

type CreateAccountRequest struct {
	UserID string `json:"user_id"`
}

type CreateAccountResponse struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type WithdrawRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type WithdrawResponse struct {
	Success    bool    `json:"success"`
	NewBalance float64 `json:"new_balance"`
	Message    string  `json:"message"`
}
