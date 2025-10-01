package services

import (
	"fmt"
	"log"
	"time"

	"billing-service/internal/dto"
	"billing-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingService struct {
	db *gorm.DB
}

func NewBillingService(db *gorm.DB) *BillingService {
	return &BillingService{
		db: db,
	}
}

func (s *BillingService) ProcessPayment(req *dto.ProcessPaymentRequest) (*dto.ProcessPaymentResponse, error) {
	log.Printf("Processing payment for order %s, amount: %.2f", req.OrderID, req.Amount)

	// Создаем платеж
	payment := &models.Payment{
		ID:            uuid.New().String(),
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Status:        "pending",
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Сохраняем в базу
	if err := s.db.Create(payment).Error; err != nil {
		log.Printf("Failed to create payment: %v", err)
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	// Симулируем обработку платежа
	success := s.simulatePaymentProcessing(req.Amount, req.PaymentMethod)

	if success {
		payment.Status = "completed"
		payment.UpdatedAt = time.Now()
		s.db.Save(payment)

		log.Printf("Payment %s completed successfully", payment.ID)
		return &dto.ProcessPaymentResponse{
			PaymentID: payment.ID,
			Status:    "completed",
			Message:   "Payment processed successfully",
		}, nil
	} else {
		payment.Status = "failed"
		payment.UpdatedAt = time.Now()
		s.db.Save(payment)

		log.Printf("Payment %s failed", payment.ID)
		return &dto.ProcessPaymentResponse{
			PaymentID: payment.ID,
			Status:    "failed",
			Message:   "Payment processing failed",
		}, nil
	}
}

func (s *BillingService) RefundPayment(req *dto.RefundPaymentRequest) (*dto.RefundPaymentResponse, error) {
	log.Printf("Processing refund for payment %s", req.PaymentID)

	var payment models.Payment
	if err := s.db.Where("id = ?", req.PaymentID).First(&payment).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %v", err)
	}

	if payment.Status != "completed" {
		return nil, fmt.Errorf("only completed payments can be refunded")
	}

	// Обновляем статус
	payment.Status = "refunded"
	payment.UpdatedAt = time.Now()
	s.db.Save(&payment)

	log.Printf("Payment %s refunded successfully", payment.ID)
	return &dto.RefundPaymentResponse{
		Status:  "refunded",
		Message: "Payment refunded successfully",
	}, nil
}

func (s *BillingService) GetPayment(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	if err := s.db.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %v", err)
	}
	return &payment, nil
}

func (s *BillingService) CreateAccount(req *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	// Проверяем, существует ли уже аккаунт
	var existingAccount models.Account
	if err := s.db.Where("user_id = ?", req.UserID).First(&existingAccount).Error; err == nil {
		return &dto.CreateAccountResponse{
			AccountID: existingAccount.ID,
			Message:   "Account already exists",
		}, nil
	}

	// Создаем новый аккаунт
	account := &models.Account{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Balance:   1000.0, // Начальный баланс
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(account).Error; err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	log.Printf("Account created for user %s with ID %s", req.UserID, account.ID)
	return &dto.CreateAccountResponse{
		AccountID: account.ID,
		Message:   "Account created successfully",
	}, nil
}

func (s *BillingService) WithdrawMoney(req *dto.WithdrawRequest) (*dto.WithdrawResponse, error) {
	var account models.Account
	if err := s.db.Where("user_id = ?", req.UserID).First(&account).Error; err != nil {
		return nil, fmt.Errorf("account not found: %v", err)
	}

	if account.Balance < req.Amount {
		return &dto.WithdrawResponse{
			Success:    false,
			NewBalance: account.Balance,
			Message:    "Insufficient funds",
		}, nil
	}

	// Списываем деньги
	account.Balance -= req.Amount
	account.UpdatedAt = time.Now()
	s.db.Save(&account)

	log.Printf("Withdrawn %.2f from account %s. New balance: %.2f", req.Amount, account.ID, account.Balance)
	return &dto.WithdrawResponse{
		Success:    true,
		NewBalance: account.Balance,
		Message:    "Money withdrawn successfully",
	}, nil
}

func (s *BillingService) simulatePaymentProcessing(amount float64, method string) bool {
	// Платежи на сумму больше 1000 всегда проходят
	if amount > 1000 {
		return true
	}

	// Платежи через карту имеют 90% успешности
	if method == "card" {
		return time.Now().UnixNano()%10 < 9
	}

	// Остальные методы имеют 80% успешности
	return time.Now().UnixNano()%10 < 8
}
