package server

import (
	"log"

	"delivery-service/internal/config"
	"delivery-service/internal/database"
	"delivery-service/internal/handlers"
	"delivery-service/internal/services"

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

	// тестовые данные
	if err := database.CreateTestData(db); err != nil {
		log.Printf("Failed to create test data: %v", err)
	}

	deliveryService := services.NewDeliveryService(db)
	handlers := handlers.NewHandlers(deliveryService)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1")
	{
		api.GET("/couriers", handlers.GetCouriers)
		api.GET("/couriers/:courier_id/slots", handlers.GetCourierSlots)
		api.POST("/delivery/reserve", handlers.ReserveDelivery)
		api.POST("/delivery/cancel-reservation", handlers.CancelDeliveryReservation)
		api.GET("/reservations/:reservation_id", handlers.GetDeliveryReservation)
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
