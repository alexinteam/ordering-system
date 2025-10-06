package dto

type SendNotificationRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Email   string `json:"email"`
}

type SendNotificationResponse struct {
	NotificationID string `json:"notification_id"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}
