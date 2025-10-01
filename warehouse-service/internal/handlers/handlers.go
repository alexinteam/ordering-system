package handlers

import (
	"net/http"

	"warehouse-service/internal/dto"
	"warehouse-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	warehouseService *services.WarehouseService
}

func NewHandlers(warehouseService *services.WarehouseService) *Handlers {
	return &Handlers{
		warehouseService: warehouseService,
	}
}

func (h *Handlers) ReserveProduct(c *gin.Context) {
	var req dto.ReserveProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.warehouseService.ReserveProduct(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) CancelReservation(c *gin.Context) {
	var req dto.CancelReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.warehouseService.CancelReservation(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetProducts(c *gin.Context) {
	products, err := h.warehouseService.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handlers) GetProduct(c *gin.Context) {
	productID := c.Param("product_id")

	product, err := h.warehouseService.GetProduct(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handlers) GetReservation(c *gin.Context) {
	reservationID := c.Param("reservation_id")

	reservation, err := h.warehouseService.GetReservation(reservationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "warehouse"})
}
