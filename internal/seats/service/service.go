package service

import (
	"backend/internal/seats"
	"backend/internal/seats/repository"
	"errors"
	"time"
)

type SeatService interface {
	GetByEventID(id string) ([]seats.Seat, error)
	GetByID(id string) (seats.Seat, error)
	ReserveSeat(id string, customerID string) error
	UpdateStatus(id string, status string) error
	ClearExpiredReservations(timeout time.Duration) error
}

type seatService struct {
	repo repository.SeatRepository
}

func NewSeatService(repo repository.SeatRepository) SeatService {
	return &seatService{repo: repo}
}

func (s *seatService) GetByEventID(id string) ([]seats.Seat, error) {
	return s.repo.GetByEventID(id)
}

func (s *seatService) GetByID(id string) (seats.Seat, error) {
	return s.repo.GetByID(id)
}

func (s *seatService) ReserveSeat(id string, customerID string) error {
	if customerID == "" {
		return errors.New("customer id is required")
	}

	seat, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if seat.Status != seats.StatusAvailable {
		// Check if it was reserved more than 5 minutes ago
		if seat.Status == seats.StatusReserved && seat.ReservedAt != nil {
			if time.Since(*seat.ReservedAt) > 5*time.Minute {
				// Can re-reserve
			} else {
				return errors.New("seat is not available")
			}
		} else {
			return errors.New("seat is not available")
		}
	}

	now := time.Now()
	return s.repo.UpdateStatus(id, seats.StatusReserved, &now, &customerID)
}

func (s *seatService) UpdateStatus(id string, status string) error {
	if status == seats.StatusAvailable {
		return s.repo.UpdateStatus(id, status, nil, nil)
	}
	now := time.Now()
	return s.repo.UpdateStatus(id, status, &now, nil)
}

func (s *seatService) ClearExpiredReservations(timeout time.Duration) error {
	return s.repo.ClearExpiredReservations(timeout)
}

