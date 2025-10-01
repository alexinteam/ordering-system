package services

import (
	"fmt"
	"log"
	"time"

	"delivery-service/internal/dto"
	"delivery-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryService struct {
	db *gorm.DB
}

func NewDeliveryService(db *gorm.DB) *DeliveryService {
	return &DeliveryService{
		db: db,
	}
}

func (s *DeliveryService) ReserveDelivery(req *dto.ReserveDeliveryRequest) (*dto.ReserveDeliveryResponse, error) {
	log.Printf("Reserving delivery for order %s, courier %s, slot %s", req.OrderID, req.CourierID, req.SlotID)

	var courier models.Courier
	if err := s.db.Where("id = ?", req.CourierID).First(&courier).Error; err != nil {
		return nil, fmt.Errorf("courier not found: %v", err)
	}

	if courier.Status != "available" {
		return nil, fmt.Errorf("courier is not available")
	}

	var slot models.DeliverySlot
	if err := s.db.Where("id = ? AND courier_id = ?", req.SlotID, req.CourierID).First(&slot).Error; err != nil {
		return nil, fmt.Errorf("delivery slot not found: %v", err)
	}

	if slot.Status != "available" {
		return nil, fmt.Errorf("delivery slot is not available")
	}

	reservation := &models.DeliveryReservation{
		ID:        uuid.New().String(),
		OrderID:   req.OrderID,
		CourierID: req.CourierID,
		SlotID:    req.SlotID,
		Address:   req.Address,
		Status:    "reserved",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := s.db.Create(reservation).Error; err != nil {
		log.Printf("Failed to create delivery reservation: %v", err)
		return nil, fmt.Errorf("failed to create delivery reservation: %v", err)
	}

	slot.Status = "reserved"
	slot.UpdatedAt = time.Now().Unix()
	s.db.Save(&slot)

	courier.Status = "busy"
	courier.UpdatedAt = time.Now().Unix()
	s.db.Save(&courier)

	log.Printf("Delivery reserved successfully. Reservation ID: %s", reservation.ID)
	return &dto.ReserveDeliveryResponse{
		ReservationID: reservation.ID,
		Status:        "reserved",
		Message:       "Delivery reserved successfully",
	}, nil
}

func (s *DeliveryService) CancelDeliveryReservation(req *dto.CancelDeliveryRequest) (*dto.CancelDeliveryResponse, error) {
	log.Printf("Cancelling delivery reservation %s", req.ReservationID)

	var reservation models.DeliveryReservation
	if err := s.db.Where("id = ?", req.ReservationID).First(&reservation).Error; err != nil {
		return nil, fmt.Errorf("delivery reservation not found: %v", err)
	}

	if reservation.Status != "reserved" {
		return nil, fmt.Errorf("only reserved deliveries can be cancelled")
	}

	var slot models.DeliverySlot
	if err := s.db.Where("id = ?", reservation.SlotID).First(&slot).Error; err == nil {
		slot.Status = "available"
		slot.UpdatedAt = time.Now().Unix()
		s.db.Save(&slot)
	}

	var courier models.Courier
	if err := s.db.Where("id = ?", reservation.CourierID).First(&courier).Error; err == nil {
		courier.Status = "available"
		courier.UpdatedAt = time.Now().Unix()
		s.db.Save(&courier)
	}

	reservation.Status = "cancelled"
	reservation.UpdatedAt = time.Now().Unix()
	s.db.Save(&reservation)

	log.Printf("Delivery reservation %s cancelled successfully", req.ReservationID)
	return &dto.CancelDeliveryResponse{
		Status:  "cancelled",
		Message: "Delivery reservation cancelled successfully",
	}, nil
}

func (s *DeliveryService) GetCouriers() ([]models.Courier, error) {
	var couriers []models.Courier
	if err := s.db.Find(&couriers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch couriers: %v", err)
	}
	return couriers, nil
}

func (s *DeliveryService) GetCourierSlots(courierID, date string) ([]models.DeliverySlot, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var slots []models.DeliverySlot
	query := s.db.Where("courier_id = ? AND date = ?", courierID, date)
	if err := query.Find(&slots).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch delivery slots: %v", err)
	}
	return slots, nil
}

func (s *DeliveryService) GetDeliveryReservation(reservationID string) (*models.DeliveryReservation, error) {
	var reservation models.DeliveryReservation
	if err := s.db.Where("id = ?", reservationID).First(&reservation).Error; err != nil {
		return nil, fmt.Errorf("delivery reservation not found: %v", err)
	}
	return &reservation, nil
}
