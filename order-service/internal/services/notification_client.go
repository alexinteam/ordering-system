package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/config"
	"order-service/internal/dto"
)

type NotificationClient struct {
	baseURL string
}

func NewNotificationClient(cfg config.ServicesConfig) *NotificationClient {
	return &NotificationClient{
		baseURL: cfg.NotificationServiceURL,
	}
}

func (c *NotificationClient) SendNotification(userID, message, notificationType, email string) error {
	req := dto.SendNotificationRequest{
		UserID:  userID,
		Message: message,
		Type:    notificationType,
		Email:   email,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(c.baseURL+"/notifications", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
