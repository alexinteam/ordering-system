package server

import (
	"billing-service/internal/config"
	"billing-service/internal/database"
	"billing-service/internal/handlers"
	"billing-service/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	cfg := config.Load()
	db, err := database.Connect(cfg.Database)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err := database.Migrate(db); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	billingService := services.NewBillingService(db)
	handlers := handlers.NewHandlers(billingService)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1")
	{
		api.POST("/payments/process", handlers.ProcessPayment)
		api.POST("/payments/refund", handlers.RefundPayment)
		api.GET("/payments/:payment_id", handlers.GetPayment)
		api.POST("/accounts", handlers.CreateAccount)
		api.POST("/payments/withdraw", handlers.WithdrawMoney)
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
