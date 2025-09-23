package services

import (
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) SendNotification(req *dto.SendNotificationRequest) (*dto.SendNotificationResponse, error) {
	notification := &models.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Message:   req.Message,
		Type:      req.Type,
		Email:     req.Email,
		Status:    "sent", // В упрощенной версии всегда считаем отправленным
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return &dto.SendNotificationResponse{
		NotificationID: notification.ID,
		Status:         "sent",
		Message:        "Notification sent successfully",
	}, nil
}

func (s *NotificationService) GetNotifications(userID string, limit, offset int) (*dto.GetNotificationsResponse, error) {
	var notifications []models.Notification
	var total int64

	if err := s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	return &dto.GetNotificationsResponse{
		Notifications: notifications,
		Total:         total,
	}, nil
}

func (s *NotificationService) GetAllNotifications(limit, offset int) (*dto.GetNotificationsResponse, error) {
	var notifications []models.Notification
	var total int64

	if err := s.db.Model(&models.Notification{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := s.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	return &dto.GetNotificationsResponse{
		Notifications: notifications,
		Total:         total,
	}, nil
}
