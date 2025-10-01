package services

import (
	"fmt"
	"log"
	"time"

	"warehouse-service/internal/dto"
	"warehouse-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WarehouseService struct {
	db *gorm.DB
}

func NewWarehouseService(db *gorm.DB) *WarehouseService {
	return &WarehouseService{
		db: db,
	}
}

func (s *WarehouseService) ReserveProduct(req *dto.ReserveProductRequest) (*dto.ReserveProductResponse, error) {
	log.Printf("Reserving product %s, quantity: %d for order %s", req.ProductID, req.Quantity, req.OrderID)

	var product models.Product
	if err := s.db.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Проверяем доступность товара
	availableStock := product.Stock - product.Reserved
	if availableStock < req.Quantity {
		log.Printf("Insufficient stock for product %s. Available: %d, Requested: %d",
			req.ProductID, availableStock, req.Quantity)
		return nil, fmt.Errorf("insufficient stock")
	}

	// Создаем резервацию
	reservation := &models.Reservation{
		ID:        uuid.New().String(),
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "reserved",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем резервацию
	if err := s.db.Create(reservation).Error; err != nil {
		log.Printf("Failed to create reservation: %v", err)
		return nil, fmt.Errorf("failed to create reservation: %v", err)
	}

	// Обновляем зарезервированное количество
	product.Reserved += req.Quantity
	product.UpdatedAt = time.Now()
	s.db.Save(&product)

	log.Printf("Product %s reserved successfully. Reservation ID: %s", req.ProductID, reservation.ID)
	return &dto.ReserveProductResponse{
		ReservationID: reservation.ID,
		Status:        "reserved",
		Message:       "Product reserved successfully",
	}, nil
}

func (s *WarehouseService) CancelReservation(req *dto.CancelReservationRequest) (*dto.CancelReservationResponse, error) {
	log.Printf("Cancelling reservation %s", req.ReservationID)

	var reservation models.Reservation
	if err := s.db.Where("id = ?", req.ReservationID).First(&reservation).Error; err != nil {
		return nil, fmt.Errorf("reservation not found: %v", err)
	}

	if reservation.Status != "reserved" {
		return nil, fmt.Errorf("only reserved items can be cancelled")
	}

	// Находим товар и уменьшаем зарезервированное количество
	var product models.Product
	if err := s.db.Where("id = ?", reservation.ProductID).First(&product).Error; err != nil {
		log.Printf("Product not found for reservation %s", req.ReservationID)
	} else {
		product.Reserved -= reservation.Quantity
		product.UpdatedAt = time.Now()
		s.db.Save(&product)
	}

	// Обновляем статус резервации
	reservation.Status = "cancelled"
	reservation.UpdatedAt = time.Now()
	s.db.Save(&reservation)

	log.Printf("Reservation %s cancelled successfully", req.ReservationID)
	return &dto.CancelReservationResponse{
		Status:  "cancelled",
		Message: "Reservation cancelled successfully",
	}, nil
}

func (s *WarehouseService) GetProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products: %v", err)
	}
	return products, nil
}

func (s *WarehouseService) GetProduct(productID string) (*models.Product, error) {
	var product models.Product
	if err := s.db.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}
	return &product, nil
}

func (s *WarehouseService) GetReservation(reservationID string) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := s.db.Where("id = ?", reservationID).First(&reservation).Error; err != nil {
		return nil, fmt.Errorf("reservation not found: %v", err)
	}
	return &reservation, nil
}
