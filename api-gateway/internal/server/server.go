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
			billing.GET("/accounts/:user_id/balance", handlers.GetBalance)
			billing.GET("/payments/:user_id", handlers.GetPayments)
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
