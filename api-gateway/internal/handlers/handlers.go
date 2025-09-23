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
	notificationServiceURL string
}

func NewHandlers() *Handlers {
	return &Handlers{
		orderServiceURL:        getEnv("ORDER_SERVICE_URL", "http://order-service:8080"),
		billingServiceURL:      getEnv("BILLING_SERVICE_URL", "http://billing-service:8081"),
		notificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8082"),
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
