package handlers

import (
	"order-service/internal/dto"
	"order-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	orderService *services.OrderService
}

func NewHandlers(orderService *services.OrderService) *Handlers {
	return &Handlers{
		orderService: orderService,
	}
}

func (h *Handlers) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := h.orderService.CreateOrder(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create order"})
		return
	}

	statusCode := 201
	if response.Status == "failed" {
		statusCode = 400
	}

	c.JSON(statusCode, response)
}

func (h *Handlers) GetOrders(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := h.orderService.GetOrders(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(200, gin.H{"orders": orders})
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
