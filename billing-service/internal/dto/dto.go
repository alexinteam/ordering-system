package dto

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// CreateAccountResponse represents the response after creating an account
type CreateAccountResponse struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

// WithdrawRequest represents the request to withdraw money from an account
type WithdrawRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// WithdrawResponse represents the response after withdrawing money
type WithdrawResponse struct {
	Success    bool    `json:"success"`
	NewBalance float64 `json:"new_balance"`
	Message    string  `json:"message"`
}

// DepositRequest represents the request to deposit money to an account
type DepositRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// DepositResponse represents the response after depositing money
type DepositResponse struct {
	Success    bool    `json:"success"`
	NewBalance float64 `json:"new_balance"`
	Message    string  `json:"message"`
}

// BalanceResponse represents the response with account balance
type BalanceResponse struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}
