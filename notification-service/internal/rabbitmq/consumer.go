package rabbitmq

import (
	"encoding/json"
	"log"
	"notification-service/internal/dto"
	"notification-service/internal/services"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	service *services.NotificationService
}

func NewConsumer(url string, service *services.NotificationService) (*Consumer, error) {
	log.Printf("Connecting to RabbitMQ at %s", url)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		conn.Close()
		return nil, err
	}

	// Declare exchange
	err = channel.ExchangeDeclare("order.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare exchange: %v", err)
		conn.Close()
		return nil, err
	}

	// Declare queue
	queue, err := channel.QueueDeclare("notification.queue", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		conn.Close()
		return nil, err
	}

	// Bind queue to exchange
	err = channel.QueueBind(queue.Name, "order.processed", "order.events", false, nil)
	if err != nil {
		log.Printf("Failed to bind queue: %v", err)
		conn.Close()
		return nil, err
	}

	log.Printf("RabbitMQ consumer created successfully")
	return &Consumer{
		conn:    conn,
		channel: channel,
		service: service,
	}, nil
}

func (c *Consumer) Start() error {
	log.Printf("Starting RabbitMQ consumer...")

	msgs, err := c.channel.Consume(
		"notification.queue", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Printf("Failed to register consumer: %v", err)
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)

			var event OrderProcessedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Process the event
			c.processEvent(&event)
		}
	}()

	log.Printf("RabbitMQ consumer started successfully")
	return nil
}

func (c *Consumer) processEvent(event *OrderProcessedEvent) {
	log.Printf("Processing order.processed event for order %s", event.OrderID)

	var message string
	var notificationType string

	if event.Status == "completed" {
		message = "Ваш заказ успешно оформлен!"
		notificationType = "success"
	} else {
		message = "Недостаточно средств на счете"
		notificationType = "error"
	}

	// Create notification
	req := &dto.SendNotificationRequest{
		UserID:  event.UserID,
		Message: message,
		Type:    notificationType,
		Email:   event.Email,
	}

	_, err := c.service.SendNotification(req)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
	} else {
		log.Printf("Notification sent successfully for order %s", event.OrderID)
	}
}

func (c *Consumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type OrderProcessedEvent struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Status  string `json:"status"`
	Email   string `json:"email"`
}
