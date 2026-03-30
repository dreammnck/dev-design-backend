package repository

import (
	"backend/internal/seats"
	"time"

	"gorm.io/gorm"
)

type SeatRepository interface {
	GetByEventID(id string) ([]seats.Seat, error)
	GetByID(id string) (seats.Seat, error)
	UpdateStatus(id string, status string, reservedAt *time.Time, customerID *string) error
	ClearExpiredReservations(timeout time.Duration) error
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) GetByEventID(id string) ([]seats.Seat, error) {
	var s []seats.Seat
	err := r.db.Where("event_id = ?", id).Find(&s).Error
	return s, err
}

func (r *seatRepository) GetByID(id string) (seats.Seat, error) {
	var s seats.Seat
	err := r.db.Where("id = ?", id).First(&s).Error
	return s, err
}

func (r *seatRepository) UpdateStatus(id string, status string, reservedAt *time.Time, customerID *string) error {
	return r.db.Model(&seats.Seat{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"reserved_at": reservedAt,
		"customer_id": customerID,
	}).Error
}

func (r *seatRepository) ClearExpiredReservations(timeout time.Duration) error {
	threshold := time.Now().Add(-timeout)
	return r.db.Model(&seats.Seat{}).
		Where("status = ?", seats.StatusReserved).
		Where("reserved_at < ?", threshold).
		Updates(map[string]interface{}{
			"status":      seats.StatusAvailable,
			"reserved_at": nil,
			"customer_id": nil,
		}).Error
}
