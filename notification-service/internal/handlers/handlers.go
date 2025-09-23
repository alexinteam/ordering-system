package handlers

import (
	"strconv"

	"notification-service/internal/dto"
	"notification-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	notificationService *services.NotificationService
}

func NewHandlers(notificationService *services.NotificationService) *Handlers {
	return &Handlers{
		notificationService: notificationService,
	}
}

func (h *Handlers) SendNotification(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := h.notificationService.SendNotification(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(201, response)
}

func (h *Handlers) GetNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	response, err := h.notificationService.GetNotifications(userID, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(200, response)
}

func (h *Handlers) GetAllNotifications(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	response, err := h.notificationService.GetAllNotifications(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(200, response)
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
