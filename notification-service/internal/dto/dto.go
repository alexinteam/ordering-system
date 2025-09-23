package dto

import "notification-service/internal/models"

type SendNotificationRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
}

type SendNotificationResponse struct {
	NotificationID string `json:"notification_id"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}

type GetNotificationsResponse struct {
	Notifications []models.Notification `json:"notifications"`
	Total         int64                 `json:"total"`
}
