package services

import (
	"billing-service/internal/dto"
	"billing-service/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountExists     = errors.New("account already exists")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type BillingService struct {
	db *gorm.DB
}

func NewBillingService(db *gorm.DB) *BillingService {
	return &BillingService{db: db}
}

// CreateAccount creates a new billing account for a user
func (s *BillingService) CreateAccount(userID string) (*models.Account, error) {
	var existingAccount models.Account
	err := s.db.Where("user_id = ?", userID).First(&existingAccount).Error
	if err == nil {
		return nil, ErrAccountExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	account := &models.Account{
		ID:        uuid.New().String(),
		UserID:    userID,
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccount retrieves an account by user ID
func (s *BillingService) GetAccount(userID string) (*models.Account, error) {
	var account models.Account
	err := s.db.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// DepositMoney adds money to an account
func (s *BillingService) DepositMoney(userID string, amount float64) (*dto.DepositResponse, error) {
	var account models.Account
	err := s.db.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update balance
	newBalance := account.Balance + amount
	if err := tx.Model(&account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create payment record
	payment := &models.Payment{
		ID:        uuid.New().String(),
		UserID:    userID,
		Amount:    amount,
		Type:      "deposit",
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.DepositResponse{
		Success:    true,
		NewBalance: newBalance,
		Message:    "Money deposited successfully",
	}, nil
}

// WithdrawMoney withdraws money from an account
func (s *BillingService) WithdrawMoney(userID string, amount float64) (*dto.WithdrawResponse, error) {
	var account models.Account
	err := s.db.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	// Check if sufficient funds
	if account.Balance < amount {
		// Create failed payment record
		payment := &models.Payment{
			ID:        uuid.New().String(),
			UserID:    userID,
			Amount:    amount,
			Type:      "withdraw",
			Status:    "failed",
			CreatedAt: time.Now(),
		}
		s.db.Create(payment)

		return &dto.WithdrawResponse{
			Success:    false,
			NewBalance: account.Balance,
			Message:    "Insufficient funds",
		}, nil
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update balance
	newBalance := account.Balance - amount
	if err := tx.Model(&account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create payment record
	payment := &models.Payment{
		ID:        uuid.New().String(),
		UserID:    userID,
		Amount:    amount,
		Type:      "withdraw",
		Status:    "success",
		CreatedAt: time.Now(),
	}

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.WithdrawResponse{
		Success:    true,
		NewBalance: newBalance,
		Message:    "Money withdrawn successfully",
	}, nil
}

// GetPayments retrieves payment history for a user
func (s *BillingService) GetPayments(userID string, limit int) ([]models.Payment, error) {
	var payments []models.Payment
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
