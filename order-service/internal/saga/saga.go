package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type SagaStep struct {
	Name        string
	Execute     func(ctx context.Context, data interface{}) error
	Compensate  func(ctx context.Context, data interface{}) error
	Completed   bool
	Compensated bool
}

type Saga struct {
	ID          string
	Steps       []*SagaStep
	Data        map[string]interface{}
	Status      string // "running", "completed", "failed", "compensated"
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type SagaManager struct {
	rabbitmqURL string
	conn        *amqp.Connection
	channel     *amqp.Channel
}

func NewSagaManager(rabbitmqURL string) (*SagaManager, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	err = channel.ExchangeDeclare(
		"saga_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	return &SagaManager{
		rabbitmqURL: rabbitmqURL,
		conn:        conn,
		channel:     channel,
	}, nil
}

func (sm *SagaManager) CreateSaga(id string, steps []*SagaStep, data map[string]interface{}) *Saga {
	return &Saga{
		ID:        id,
		Steps:     steps,
		Data:      data,
		Status:    "running",
		CreatedAt: time.Now(),
	}
}

func (sm *SagaManager) ExecuteSaga(ctx context.Context, saga *Saga) error {
	log.Printf("Starting saga %s with %d steps", saga.ID, len(saga.Steps))

	for i, step := range saga.Steps {
		log.Printf("Executing step %d: %s", i+1, step.Name)

		err := step.Execute(ctx, saga.Data)
		if err != nil {
			log.Printf("Step %s failed: %v", step.Name, err)
			saga.Status = "failed"

			// Выполняем компенсирующие действия
			sm.compensateSaga(ctx, saga, i)
			return fmt.Errorf("saga failed at step %s: %v", step.Name, err)
		}

		step.Completed = true
		log.Printf("Step %s completed successfully", step.Name)

		// Публикуем событие о завершении шага
		sm.publishSagaEvent(saga.ID, "step_completed", map[string]interface{}{
			"step_name":  step.Name,
			"step_index": i,
		})
	}

	saga.Status = "completed"
	now := time.Now()
	saga.CompletedAt = &now

	log.Printf("Saga %s completed successfully", saga.ID)

	sm.publishSagaEvent(saga.ID, "saga_completed", nil)

	return nil
}

func (sm *SagaManager) compensateSaga(ctx context.Context, saga *Saga, failedStepIndex int) {
	log.Printf("Starting compensation for saga %s", saga.ID)

	for i := failedStepIndex - 1; i >= 0; i-- {
		step := saga.Steps[i]
		if step.Completed && !step.Compensated {
			log.Printf("Compensating step: %s", step.Name)

			err := step.Compensate(ctx, saga.Data)
			if err != nil {
				log.Printf("Compensation failed for step %s: %v", step.Name, err)
			} else {
				step.Compensated = true
				log.Printf("Step %s compensated successfully", step.Name)
			}

			sm.publishSagaEvent(saga.ID, "step_compensated", map[string]interface{}{
				"step_name":  step.Name,
				"step_index": i,
			})
		}
	}

	saga.Status = "compensated"
	log.Printf("Saga %s compensation completed", saga.ID)

	// Публикуем событие о завершении компенсации
	sm.publishSagaEvent(saga.ID, "saga_compensated", nil)
}

func (sm *SagaManager) publishSagaEvent(sagaID, eventType string, data map[string]interface{}) {
	event := map[string]interface{}{
		"saga_id":    sagaID,
		"event_type": eventType,
		"timestamp":  time.Now(),
		"data":       data,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal saga event: %v", err)
		return
	}

	err = sm.channel.Publish(
		"saga_events",
		fmt.Sprintf("saga.%s.%s", sagaID, eventType),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish saga event: %v", err)
	}
}

func (sm *SagaManager) Close() error {
	if sm.channel != nil {
		sm.channel.Close()
	}
	if sm.conn != nil {
		return sm.conn.Close()
	}
	return nil
}
