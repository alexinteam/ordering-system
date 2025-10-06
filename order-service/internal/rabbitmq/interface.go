package rabbitmq

type ProducerInterface interface {
	PublishOrderCreated(orderID, userID string, amount float64, email string) error
	PublishOrderProcessed(orderID, userID, status, email string) error
	Close() error
}
