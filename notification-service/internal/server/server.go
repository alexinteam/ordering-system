package server

import (
	"log"
	"notification-service/internal/config"
	"notification-service/internal/database"
	"notification-service/internal/handlers"
	"notification-service/internal/rabbitmq"
	"notification-service/internal/services"
	"os"

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

	notificationService := services.NewNotificationService(db)

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://admin:admin@rabbitmq.infrastructure.svc.cluster.local:5672/"
	}

	log.Printf("Creating RabbitMQ consumer...")
	consumer, err := rabbitmq.NewConsumer(rabbitmqURL, notificationService)
	if err == nil {
		log.Printf("RabbitMQ consumer created successfully")
		go consumer.Start()
	} else {
		log.Printf("Failed to create RabbitMQ consumer: %v", err)
	}

	handlers := handlers.NewHandlers(notificationService)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1")
	{
		api.POST("/notifications", handlers.SendNotification)
		api.GET("/notifications/:user_id", handlers.GetNotifications)
		api.GET("/notifications", handlers.GetAllNotifications)
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
