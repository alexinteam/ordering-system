package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	orderServiceURL        string
	billingServiceURL      string
	warehouseServiceURL    string
	deliveryServiceURL     string
	notificationServiceURL string
}

func NewHandlers() *Handlers {
	return &Handlers{
		orderServiceURL:        getEnv("ORDER_SERVICE_URL", "http://order-service:8081"),
		billingServiceURL:      getEnv("BILLING_SERVICE_URL", "http://billing-service:8082"),
		warehouseServiceURL:    getEnv("WAREHOUSE_SERVICE_URL", "http://warehouse-service:8083"),
		deliveryServiceURL:     getEnv("DELIVERY_SERVICE_URL", "http://delivery-service:8084"),
		notificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8085"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (h *Handlers) proxyRequest(c *gin.Context, serviceURL, path string) {
	targetURL := serviceURL + path

	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create request"})
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to make request"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read response"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// Order Service Routes
func (h *Handlers) CreateOrder(c *gin.Context) {
	h.proxyRequest(c, h.orderServiceURL, "/api/v1/orders")
}

func (h *Handlers) GetOrders(c *gin.Context) {
	userID := c.Param("user_id")
	h.proxyRequest(c, h.orderServiceURL, "/api/v1/orders/"+userID)
}

// Billing Service Routes
func (h *Handlers) CreateAccount(c *gin.Context) {
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/accounts")
}

func (h *Handlers) DepositMoney(c *gin.Context) {
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/payments/deposit")
}

func (h *Handlers) WithdrawMoney(c *gin.Context) {
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/payments/withdraw")
}

func (h *Handlers) GetBalance(c *gin.Context) {
	userID := c.Param("user_id")
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/accounts/"+userID+"/balance")
}

func (h *Handlers) GetPayments(c *gin.Context) {
	userID := c.Param("user_id")
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/payments/"+userID)
}

// Notification Service Routes
func (h *Handlers) SendNotification(c *gin.Context) {
	h.proxyRequest(c, h.notificationServiceURL, "/api/v1/notifications")
}

func (h *Handlers) GetNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	h.proxyRequest(c, h.notificationServiceURL, "/api/v1/notifications/"+userID)
}

func (h *Handlers) GetAllNotifications(c *gin.Context) {
	h.proxyRequest(c, h.notificationServiceURL, "/api/v1/notifications")
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "service": "api-gateway"})
}

// Service health checks
func (h *Handlers) CheckOrderService(c *gin.Context) {
	h.proxyRequest(c, h.orderServiceURL, "/health")
}

func (h *Handlers) CheckBillingService(c *gin.Context) {
	h.proxyRequest(c, h.billingServiceURL, "/health")
}

func (h *Handlers) CheckNotificationService(c *gin.Context) {
	h.proxyRequest(c, h.notificationServiceURL, "/health")
}

func (h *Handlers) CheckWarehouseService(c *gin.Context) {
	h.proxyRequest(c, h.warehouseServiceURL, "/health")
}

func (h *Handlers) CheckDeliveryService(c *gin.Context) {
	h.proxyRequest(c, h.deliveryServiceURL, "/health")
}

// Warehouse Service Routes
func (h *Handlers) GetProducts(c *gin.Context) {
	h.proxyRequest(c, h.warehouseServiceURL, "/api/v1/products")
}

func (h *Handlers) GetProduct(c *gin.Context) {
	productID := c.Param("product_id")
	h.proxyRequest(c, h.warehouseServiceURL, "/api/v1/products/"+productID)
}

func (h *Handlers) ReserveProduct(c *gin.Context) {
	h.proxyRequest(c, h.warehouseServiceURL, "/api/v1/products/reserve")
}

func (h *Handlers) CancelReservation(c *gin.Context) {
	h.proxyRequest(c, h.warehouseServiceURL, "/api/v1/products/cancel-reservation")
}

// Delivery Service Routes
func (h *Handlers) GetCouriers(c *gin.Context) {
	h.proxyRequest(c, h.deliveryServiceURL, "/api/v1/couriers")
}

func (h *Handlers) GetCourierSlots(c *gin.Context) {
	courierID := c.Param("courier_id")
	h.proxyRequest(c, h.deliveryServiceURL, "/api/v1/couriers/"+courierID+"/slots")
}

func (h *Handlers) ReserveDelivery(c *gin.Context) {
	h.proxyRequest(c, h.deliveryServiceURL, "/api/v1/delivery/reserve")
}

func (h *Handlers) CancelDeliveryReservation(c *gin.Context) {
	h.proxyRequest(c, h.deliveryServiceURL, "/api/v1/delivery/cancel-reservation")
}

// Billing Service Additional Routes
func (h *Handlers) ProcessPayment(c *gin.Context) {
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/payments/process")
}

func (h *Handlers) RefundPayment(c *gin.Context) {
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/payments/refund")
}

func (h *Handlers) GetPayment(c *gin.Context) {
	paymentID := c.Param("payment_id")
	h.proxyRequest(c, h.billingServiceURL, "/api/v1/payments/"+paymentID)
}
