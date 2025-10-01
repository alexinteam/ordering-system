package handlers

import (
	"net/http"

	"delivery-service/internal/dto"
	"delivery-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	deliveryService *services.DeliveryService
}

func NewHandlers(deliveryService *services.DeliveryService) *Handlers {
	return &Handlers{
		deliveryService: deliveryService,
	}
}

func (h *Handlers) ReserveDelivery(c *gin.Context) {
	var req dto.ReserveDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.deliveryService.ReserveDelivery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) CancelDeliveryReservation(c *gin.Context) {
	var req dto.CancelDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.deliveryService.CancelDeliveryReservation(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetCouriers(c *gin.Context) {
	couriers, err := h.deliveryService.GetCouriers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, couriers)
}

func (h *Handlers) GetCourierSlots(c *gin.Context) {
	courierID := c.Param("courier_id")
	date := c.Query("date")

	slots, err := h.deliveryService.GetCourierSlots(courierID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slots)
}

func (h *Handlers) GetDeliveryReservation(c *gin.Context) {
	reservationID := c.Param("reservation_id")

	reservation, err := h.deliveryService.GetDeliveryReservation(reservationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "delivery"})
}
