package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"order-service/internal/dto"
	"order-service/internal/models"
	"order-service/internal/rabbitmq"
	"order-service/internal/saga"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db                     *gorm.DB
	billingServiceURL      string
	warehouseServiceURL    string
	deliveryServiceURL     string
	notificationServiceURL string
	rabbitmqProducer       rabbitmq.ProducerInterface
	sagaManager            *saga.SagaManager
}

func NewOrderService(db *gorm.DB, billingServiceURL, warehouseServiceURL, deliveryServiceURL, notificationServiceURL string, rabbitmqProducer rabbitmq.ProducerInterface, sagaManager *saga.SagaManager) *OrderService {
	return &OrderService{
		db:                     db,
		billingServiceURL:      billingServiceURL,
		warehouseServiceURL:    warehouseServiceURL,
		deliveryServiceURL:     deliveryServiceURL,
		notificationServiceURL: notificationServiceURL,
		rabbitmqProducer:       rabbitmqProducer,
		sagaManager:            sagaManager,
	}
}

func (s *OrderService) CreateOrder(req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	// Проверяем идемпотентность
	existingRecord, err := s.checkIdempotency(req.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	if existingRecord != nil {
		// Возвращаем сохраненный ответ
		var response dto.CreateOrderResponse
		if err := json.Unmarshal([]byte(existingRecord.Response), &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal saved response: %v", err)
		}
		log.Printf("Returning cached response for idempotency key: %s", req.IdempotencyKey)
		return &response, nil
	}

	// Создаем запись идемпотентности в статусе "processing"
	idempotencyRecord := &models.IdempotencyRecord{
		Key:       req.IdempotencyKey,
		Status:    "processing",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(idempotencyRecord).Error; err != nil {
		return nil, fmt.Errorf("failed to create idempotency record: %v", err)
	}

	order := &models.Order{
		ID:             uuid.New().String(),
		UserID:         req.UserID,
		Amount:         req.Amount,
		Status:         "pending",
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.Create(order).Error; err != nil {
		// Удаляем запись идемпотентности при ошибке
		s.db.Delete(idempotencyRecord)
		return nil, err
	}

	// Обновляем запись идемпотентности с ID заказа
	idempotencyRecord.OrderID = order.ID
	s.db.Save(idempotencyRecord)

	log.Printf("Starting saga for order %s", order.ID)

	sagaData := map[string]interface{}{
		"order_id":   order.ID,
		"user_id":    req.UserID,
		"amount":     req.Amount,
		"email":      req.Email,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
		"address":    req.Address,
	}

	steps := []*saga.SagaStep{
		{
			Name: "process_payment",
			Execute: func(ctx context.Context, data interface{}) error {
				return s.processPayment(ctx, data)
			},
			Compensate: func(ctx context.Context, data interface{}) error {
				return s.refundPayment(ctx, data)
			},
		},
		{
			Name: "reserve_product",
			Execute: func(ctx context.Context, data interface{}) error {
				return s.reserveProduct(ctx, data)
			},
			Compensate: func(ctx context.Context, data interface{}) error {
				return s.cancelProductReservation(ctx, data)
			},
		},
		{
			Name: "reserve_delivery",
			Execute: func(ctx context.Context, data interface{}) error {
				return s.reserveDelivery(ctx, data)
			},
			Compensate: func(ctx context.Context, data interface{}) error {
				return s.cancelDeliveryReservation(ctx, data)
			},
		},
	}

	sagaInstance := s.sagaManager.CreateSaga(order.ID, steps, sagaData)

	err = s.sagaManager.ExecuteSaga(context.Background(), sagaInstance)
	if err != nil {
		log.Printf("Saga failed for order %s: %v", order.ID, err)
		s.updateOrderStatus(order.ID, "failed")

		s.sendNotification(req.UserID, req.Email, "error", "Заказ не удалось оформить. Попробуйте еще раз.")

		response := &dto.CreateOrderResponse{
			OrderID: order.ID,
			Status:  "failed",
			Message: "Order creation failed",
		}

		// Сохраняем ответ в записи идемпотентности
		s.saveIdempotencyResponse(req.IdempotencyKey, response, "failed")

		return response, nil
	}

	s.updateOrderStatus(order.ID, "completed")

	s.sendNotification(req.UserID, req.Email, "success", "Ваш заказ успешно оформлен!")

	if s.rabbitmqProducer != nil {
		s.rabbitmqProducer.PublishOrderProcessed(order.ID, req.UserID, "completed", req.Email)
	}

	log.Printf("Order %s completed successfully", order.ID)

	response := &dto.CreateOrderResponse{
		OrderID: order.ID,
		Status:  "completed",
		Message: "Order created successfully",
	}

	// Сохраняем ответ в записи идемпотентности
	s.saveIdempotencyResponse(req.IdempotencyKey, response, "completed")

	return response, nil
}

func (s *OrderService) processPayment(ctx context.Context, data interface{}) error {
	sagaData := data.(map[string]interface{})
	orderID := sagaData["order_id"].(string)
	userID := sagaData["user_id"].(string)
	amount := sagaData["amount"].(float64)

	log.Printf("Processing payment for order %s, amount: %.2f", orderID, amount)

	if err := s.createAccountIfNotExists(userID); err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}

	paymentReq := map[string]interface{}{
		"order_id":       orderID,
		"user_id":        userID,
		"amount":         amount,
		"payment_method": "card",
	}

	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return fmt.Errorf("failed to marshal payment request: %v", err)
	}

	resp, err := http.Post(s.billingServiceURL+"/api/v1/payments/process", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to process payment: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("payment failed with status %d: %s", resp.StatusCode, string(body))
	}

	var paymentResp ProcessPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return fmt.Errorf("failed to decode payment response: %v", err)
	}

	if paymentResp.Status != "completed" {
		return fmt.Errorf("payment not completed: %s", paymentResp.Message)
	}

	sagaData["payment_id"] = paymentResp.PaymentID

	log.Printf("Payment processed successfully for order %s", orderID)
	return nil
}

func (s *OrderService) refundPayment(ctx context.Context, data interface{}) error {
	sagaData := data.(map[string]interface{})
	paymentID, exists := sagaData["payment_id"]
	if !exists {
		log.Printf("No payment ID found for refund")
		return nil
	}

	log.Printf("Refunding payment %s", paymentID)

	refundReq := map[string]string{
		"payment_id": paymentID.(string),
		"reason":     "Order cancelled",
	}

	jsonData, err := json.Marshal(refundReq)
	if err != nil {
		return fmt.Errorf("failed to marshal refund request: %v", err)
	}

	resp, err := http.Post(s.billingServiceURL+"/api/v1/payments/refund", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to refund payment: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refund failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Payment %s refunded successfully", paymentID)
	return nil
}

func (s *OrderService) reserveProduct(ctx context.Context, data interface{}) error {
	sagaData := data.(map[string]interface{})
	orderID := sagaData["order_id"].(string)
	productID := sagaData["product_id"].(string)
	quantity := sagaData["quantity"].(int)

	log.Printf("Reserving product %s, quantity: %d for order %s", productID, quantity, orderID)

	reserveReq := map[string]interface{}{
		"order_id":   orderID,
		"product_id": productID,
		"quantity":   quantity,
	}

	jsonData, err := json.Marshal(reserveReq)
	if err != nil {
		return fmt.Errorf("failed to marshal reserve request: %v", err)
	}

	resp, err := http.Post(s.warehouseServiceURL+"/api/v1/products/reserve", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to reserve product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("product reservation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var reserveResp ReserveProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&reserveResp); err != nil {
		return fmt.Errorf("failed to decode reserve response: %v", err)
	}

	if reserveResp.Status != "reserved" {
		return fmt.Errorf("product not reserved: %s", reserveResp.Message)
	}

	sagaData["reservation_id"] = reserveResp.ReservationID

	log.Printf("Product reserved successfully for order %s", orderID)
	return nil
}

func (s *OrderService) cancelProductReservation(ctx context.Context, data interface{}) error {
	sagaData := data.(map[string]interface{})
	reservationID, exists := sagaData["reservation_id"]
	if !exists {
		log.Printf("No reservation ID found for cancellation")
		return nil
	}

	log.Printf("Cancelling product reservation %s", reservationID)

	cancelReq := map[string]string{
		"reservation_id": reservationID.(string),
		"reason":         "Order cancelled",
	}

	jsonData, err := json.Marshal(cancelReq)
	if err != nil {
		return fmt.Errorf("failed to marshal cancel request: %v", err)
	}

	resp, err := http.Post(s.warehouseServiceURL+"/api/v1/products/cancel-reservation", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to cancel reservation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("reservation cancellation failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Product reservation %s cancelled successfully", reservationID)
	return nil
}

func (s *OrderService) reserveDelivery(ctx context.Context, data interface{}) error {
	sagaData := data.(map[string]interface{})
	orderID := sagaData["order_id"].(string)
	address := sagaData["address"].(string)

	log.Printf("Reserving delivery for order %s to address: %s", orderID, address)

	resp, err := http.Get(s.deliveryServiceURL + "/api/v1/couriers")
	if err != nil {
		return fmt.Errorf("failed to get couriers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get couriers with status %d", resp.StatusCode)
	}

	var couriers []Courier
	if err := json.NewDecoder(resp.Body).Decode(&couriers); err != nil {
		return fmt.Errorf("failed to decode couriers: %v", err)
	}

	if len(couriers) == 0 {
		return fmt.Errorf("no available couriers")
	}

	courier := couriers[0]

	resp, err = http.Get(fmt.Sprintf("%s/api/v1/couriers/%s/slots", s.deliveryServiceURL, courier.ID))
	if err != nil {
		return fmt.Errorf("failed to get courier slots: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get courier slots with status %d", resp.StatusCode)
	}

	var slots []DeliverySlot
	if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
		return fmt.Errorf("failed to decode slots: %v", err)
	}

	if len(slots) == 0 {
		return fmt.Errorf("no available delivery slots")
	}

	slot := slots[0]

	reserveReq := map[string]interface{}{
		"order_id":   orderID,
		"courier_id": courier.ID,
		"slot_id":    slot.ID,
		"address":    address,
	}

	jsonData, err := json.Marshal(reserveReq)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery reserve request: %v", err)
	}

	resp, err = http.Post(s.deliveryServiceURL+"/api/v1/delivery/reserve", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to reserve delivery: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delivery reservation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var deliveryResp ReserveDeliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&deliveryResp); err != nil {
		return fmt.Errorf("failed to decode delivery response: %v", err)
	}

	if deliveryResp.Status != "reserved" {
		return fmt.Errorf("delivery not reserved: %s", deliveryResp.Message)
	}

	sagaData["delivery_reservation_id"] = deliveryResp.ReservationID

	log.Printf("Delivery reserved successfully for order %s", orderID)
	return nil
}

func (s *OrderService) cancelDeliveryReservation(ctx context.Context, data interface{}) error {
	sagaData := data.(map[string]interface{})
	deliveryReservationID, exists := sagaData["delivery_reservation_id"]
	if !exists {
		log.Printf("No delivery reservation ID found for cancellation")
		return nil
	}

	log.Printf("Cancelling delivery reservation %s", deliveryReservationID)

	cancelReq := map[string]string{
		"reservation_id": deliveryReservationID.(string),
		"reason":         "Order cancelled",
	}

	jsonData, err := json.Marshal(cancelReq)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery cancel request: %v", err)
	}

	resp, err := http.Post(s.deliveryServiceURL+"/api/v1/delivery/cancel-reservation", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to cancel delivery reservation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delivery reservation cancellation failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Delivery reservation %s cancelled successfully", deliveryReservationID)
	return nil
}

func (s *OrderService) createAccountIfNotExists(userID string) error {
	createAccountReq := map[string]string{
		"user_id": userID,
	}

	jsonData, err := json.Marshal(createAccountReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.billingServiceURL+"/api/v1/accounts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 409 {
		return nil
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create account, status: %d", resp.StatusCode)
	}

	return nil
}

func (s *OrderService) withdrawMoney(userID string, amount float64) (*WithdrawResponse, error) {
	withdrawReq := map[string]interface{}{
		"user_id": userID,
		"amount":  amount,
	}

	jsonData, err := json.Marshal(withdrawReq)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.billingServiceURL+"/api/v1/payments/withdraw", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var withdrawResp WithdrawResponse
	if err := json.Unmarshal(body, &withdrawResp); err != nil {
		return nil, err
	}

	return &withdrawResp, nil
}

func (s *OrderService) sendNotification(userID, email, notificationType, message string) {
	notificationReq := map[string]string{
		"user_id": userID,
		"email":   email,
		"type":    notificationType,
		"message": message,
	}

	jsonData, err := json.Marshal(notificationReq)
	if err != nil {
		return // Log error in production
	}

	http.Post(s.notificationServiceURL+"/api/v1/notifications", "application/json", bytes.NewBuffer(jsonData))
}

func (s *OrderService) updateOrderStatus(orderID, status string) {
	s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status)
}

func (s *OrderService) GetOrders(userID string) ([]models.Order, error) {
	var orders []models.Order
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

type WithdrawResponse struct {
	Success    bool    `json:"success"`
	NewBalance float64 `json:"new_balance"`
	Message    string  `json:"message"`
}

type ProcessPaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type ReserveProductResponse struct {
	ReservationID string `json:"reservation_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type ReserveDeliveryResponse struct {
	ReservationID string `json:"reservation_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type Courier struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type DeliverySlot struct {
	ID        string `json:"id"`
	CourierID string `json:"courier_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// Методы для работы с идемпотентностью
func (s *OrderService) checkIdempotency(key string) (*models.IdempotencyRecord, error) {
	var record models.IdempotencyRecord
	err := s.db.Where("key = ?", key).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Запись не найдена - это нормально
		}
		return nil, err
	}
	return &record, nil
}

func (s *OrderService) saveIdempotencyResponse(key string, response *dto.CreateOrderResponse, status string) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	return s.db.Model(&models.IdempotencyRecord{}).
		Where("key = ?", key).
		Updates(map[string]interface{}{
			"status":     status,
			"response":   string(responseJSON),
			"updated_at": time.Now(),
		}).Error
}
