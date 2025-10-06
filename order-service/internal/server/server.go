package server

import (
	"log"
	"order-service/internal/config"
	"order-service/internal/database"
	"order-service/internal/handlers"
	"order-service/internal/rabbitmq"
	"order-service/internal/saga"
	"order-service/internal/services"
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

	billingServiceURL := os.Getenv("BILLING_SERVICE_URL")
	if billingServiceURL == "" {
		billingServiceURL = "http://billing-service:8082"
	}

	warehouseServiceURL := os.Getenv("WAREHOUSE_SERVICE_URL")
	if warehouseServiceURL == "" {
		warehouseServiceURL = "http://warehouse-service:8083"
	}

	deliveryServiceURL := os.Getenv("DELIVERY_SERVICE_URL")
	if deliveryServiceURL == "" {
		deliveryServiceURL = "http://delivery-service:8084"
	}

	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notification-service:8085"
	}

	log.Printf("Creating RabbitMQ producer...")
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://admin:admin@rabbitmq.infrastructure.svc.cluster.local:5672/"
	}

	rabbitmqProducer, err := rabbitmq.NewProducer(rabbitmqURL)
	if err != nil {
		log.Printf("Failed to create RabbitMQ producer: %v", err)
		rabbitmqProducer = nil
	} else {
		log.Printf("RabbitMQ producer created successfully")
	}

	log.Printf("Creating Saga manager...")
	sagaManager, err := saga.NewSagaManager(rabbitmqURL)
	if err != nil {
		log.Printf("Failed to create Saga manager: %v", err)
		panic("Failed to create Saga manager: " + err.Error())
	} else {
		log.Printf("Saga manager created successfully")
	}

	orderService := services.NewOrderService(db, billingServiceURL, warehouseServiceURL, deliveryServiceURL, notificationServiceURL, rabbitmqProducer, sagaManager)

	handlers := handlers.NewHandlers(orderService)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1")
	{
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders/:user_id", handlers.GetOrders)
	}

	return &Server{
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
