package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/config"
	"order-service/internal/dto"
)

type BillingClient struct {
	baseURL string
}

func NewBillingClient(cfg config.ServicesConfig) *BillingClient {
	return &BillingClient{
		baseURL: cfg.BillingServiceURL,
	}
}

func (c *BillingClient) CreateAccount(userID string) error {
	req := dto.CreateAccountRequest{UserID: userID}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(c.baseURL+"/accounts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 409 { // 409 means account already exists
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *BillingClient) WithdrawMoney(userID string, amount float64) (*dto.WithdrawResponse, error) {
	req := dto.WithdrawRequest{UserID: userID, Amount: amount}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(c.baseURL+"/payments/withdraw", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw money: %w", err)
	}
	defer resp.Body.Close()

	var withdrawResp dto.WithdrawResponse
	if err := json.NewDecoder(resp.Body).Decode(&withdrawResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &withdrawResp, nil
}
