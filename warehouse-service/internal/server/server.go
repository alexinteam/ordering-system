package server

import (
	"log"

	"warehouse-service/internal/config"
	"warehouse-service/internal/database"
	"warehouse-service/internal/handlers"
	"warehouse-service/internal/services"

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

	// Создаем тестовые товары
	if err := database.CreateTestProducts(db); err != nil {
		log.Printf("Failed to create test products: %v", err)
	}

	warehouseService := services.NewWarehouseService(db)
	handlers := handlers.NewHandlers(warehouseService)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1")
	{
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:product_id", handlers.GetProduct)
		api.POST("/products/reserve", handlers.ReserveProduct)
		api.POST("/products/cancel-reservation", handlers.CancelReservation)
		api.GET("/reservations/:reservation_id", handlers.GetReservation)
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
