package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(url string) (*Producer, error) {
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

	// Declare exchanges
	err = channel.ExchangeDeclare("order.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare exchange: %v", err)
		conn.Close()
		return nil, err
	}

	log.Printf("RabbitMQ producer created successfully")
	return &Producer{
		conn:    conn,
		channel: channel,
	}, nil
}

func (p *Producer) PublishOrderCreated(orderID, userID string, amount float64, email string) error {
	event := OrderCreatedEvent{
		OrderID: orderID,
		UserID:  userID,
		Amount:  amount,
		Email:   email,
	}

	return p.publishEvent("order.created", event)
}

func (p *Producer) PublishOrderProcessed(orderID, userID, status, email string) error {
	event := OrderProcessedEvent{
		OrderID: orderID,
		UserID:  userID,
		Status:  status,
		Email:   email,
	}

	return p.publishEvent("order.processed", event)
}

func (p *Producer) publishEvent(routingKey string, event interface{}) error {
	log.Printf("Publishing event with routing key %s", routingKey)

	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	err = p.channel.Publish(
		"order.events", // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Successfully published event with routing key %s", routingKey)
	return nil
}

func (p *Producer) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

type OrderCreatedEvent struct {
	OrderID string  `json:"order_id"`
	UserID  string  `json:"user_id"`
	Amount  float64 `json:"amount"`
	Email   string  `json:"email"`
}

type OrderProcessedEvent struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Status  string `json:"status"`
	Email   string `json:"email"`
}
