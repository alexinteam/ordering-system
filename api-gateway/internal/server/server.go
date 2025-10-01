package server

import (
	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	handlers := handlers.NewHandlers()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", handlers.Health)

	router.GET("/health/orders", handlers.CheckOrderService)
	router.GET("/health/billing", handlers.CheckBillingService)
	router.GET("/health/warehouse", handlers.CheckWarehouseService)
	router.GET("/health/delivery", handlers.CheckDeliveryService)
	router.GET("/health/notifications", handlers.CheckNotificationService)

	api := router.Group("/api/v1")
	{
		orders := api.Group("/orders")
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("/:user_id", handlers.GetOrders)
		}

		billing := api.Group("/billing")
		{
			billing.POST("/accounts", handlers.CreateAccount)
			billing.POST("/payments/deposit", handlers.DepositMoney)
			billing.POST("/payments/withdraw", handlers.WithdrawMoney)
			billing.POST("/payments/process", handlers.ProcessPayment)
			billing.POST("/payments/refund", handlers.RefundPayment)
			billing.GET("/accounts/:user_id/balance", handlers.GetBalance)
			billing.GET("/payments/user/:user_id", handlers.GetPayments)
			billing.GET("/payments/payment/:payment_id", handlers.GetPayment)
		}

		warehouse := api.Group("/warehouse")
		{
			warehouse.GET("/products", handlers.GetProducts)
			warehouse.GET("/products/:product_id", handlers.GetProduct)
			warehouse.POST("/products/reserve", handlers.ReserveProduct)
			warehouse.POST("/products/cancel-reservation", handlers.CancelReservation)
		}

		delivery := api.Group("/delivery")
		{
			delivery.GET("/couriers", handlers.GetCouriers)
			delivery.GET("/couriers/:courier_id/slots", handlers.GetCourierSlots)
			delivery.POST("/reserve", handlers.ReserveDelivery)
			delivery.POST("/cancel-reservation", handlers.CancelDeliveryReservation)
		}

		notifications := api.Group("/notifications")
		{
			notifications.POST("", handlers.SendNotification)
			notifications.GET("/:user_id", handlers.GetNotifications)
			notifications.GET("", handlers.GetAllNotifications)
		}
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
