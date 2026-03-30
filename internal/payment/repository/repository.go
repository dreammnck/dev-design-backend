package repository

import (
	"backend/internal/payment"
	"fmt"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreateBookingAndPayment(req payment.PaymentReq, transactionID string, status string) (string, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// CreateBookingAndPayment executes a transaction to create booking, link seats, and record payment
func (r *paymentRepository) CreateBookingAndPayment(req payment.PaymentReq, transactionID string, status string) (string, error) {
	var bookingID string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get EventID from the first seat (assuming all seats belong to the same event)
		var seat struct {
			EventID string
		}
		if err := tx.Table("seats").Select("event_id").Where("id = ?", req.SeatIDs[0]).Scan(&seat).Error; err != nil {
			return fmt.Errorf("failed to find seat event: %v", err)
		}

		// 2. Create Booking
		newBooking := struct {
			EventID     string `gorm:"column:event_id"`
			CustomerID  string `gorm:"column:customer_id"`
			TotalAmount int    `gorm:"column:total_amount"`
			Status      string `gorm:"column:status"`
		}{
			EventID:     seat.EventID,
			CustomerID:  req.CustomerID,
			TotalAmount: int(req.Amount),
			Status:      "confirmed", // If payment is success
		}

		if err := tx.Table("bookings").Create(&newBooking).Error; err != nil {
			return fmt.Errorf("failed to create booking: %v", err)
		}

		// Get the generated Booking ID if using gen_random_uuid
		// Since we didn't define a struct with GORM tags for bookings here, we can query it back or use a more defined struct
		var createdBooking struct {
			ID string
		}
		tx.Table("bookings").Order("created_at desc").Limit(1).Scan(&createdBooking)
		bookingID = createdBooking.ID

		// 3. Link Seats to Booking
		for _, seatID := range req.SeatIDs {
			bookingSeat := map[string]interface{}{
				"booking_id": bookingID,
				"seat_id":    seatID,
			}
			if err := tx.Table("booking_seats").Create(&bookingSeat).Error; err != nil {
				return fmt.Errorf("failed to link seat %s: %v", seatID, err)
			}
		}

		// 4. Create Payment record
		paymentRecord := map[string]interface{}{
			"booking_id":     bookingID,
			"transaction_id": transactionID,
			"amount":         int(req.Amount),
			"payment_method": req.PaymentMethod,
			"status":         status,
		}
		if err := tx.Table("payments").Create(&paymentRecord).Error; err != nil {
			return fmt.Errorf("failed to create payment record: %v", err)
		}

		return nil
	})

	return bookingID, err
}
