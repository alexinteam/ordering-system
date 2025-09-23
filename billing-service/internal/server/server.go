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
	// Initialize database
	cfg := config.Load()
	db, err := database.Connect(cfg.Database)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Initialize services
	billingService := services.NewBillingService(db)

	// Initialize handlers
	handlers := handlers.NewHandlers(billingService)

	// Setup routes
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handlers.Health)

	// API routes
	api := router.Group("/api/v1")
	{
		// Account management
		api.POST("/accounts", handlers.CreateAccount)
		api.GET("/accounts/:user_id/balance", handlers.GetBalance)

		// Payment operations
		api.POST("/payments/deposit", handlers.DepositMoney)
		api.POST("/payments/withdraw", handlers.WithdrawMoney)
		api.GET("/payments/:user_id", handlers.GetPayments)
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
