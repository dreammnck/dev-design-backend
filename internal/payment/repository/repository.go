package repository

import (
	"backend/internal/payment"
	"fmt"
	"strings"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreateBookingAndPayment(req payment.PaymentReq, transactionID string, status string) (string, error)
	GetBookingsByUserID(userID string) ([]payment.MyBookingRes, error)
	GetBookingByID(bookingID string, userID string) ([]payment.MyBookingRes, error)
	GetAllBookings() ([]map[string]interface{}, error)
	GetAllPayments() ([]map[string]interface{}, error)
	UpdateBookingSeatQR(bookingID, seatID, qrPayload, qrUri string) error
	RedeemTicket(bookingID, seatID string) error
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
		type Booking struct {
			ID          string `gorm:"column:id;primaryKey"`
			EventID     string `gorm:"column:event_id"`
			UserID      string `gorm:"column:user_id"`
			TotalAmount int    `gorm:"column:total_amount"`
			Status      string `gorm:"column:status"`
		}

		b := Booking{
			ID:          uuid.New().String(),
			EventID:     seat.EventID,
			UserID:      req.UserID,
			TotalAmount: int(req.Amount),
			Status:      "confirmed",
		}
		bookingID = b.ID
		fmt.Printf("DEBUG: Creating booking with ID: [%s]\n", bookingID)

		if err := tx.Table("bookings").Create(&b).Error; err != nil {
			return fmt.Errorf("failed to create booking: %v", err)
		}

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

func (r *paymentRepository) GetAllBookings() ([]map[string]interface{}, error) {
	var bookings []map[string]interface{}
	err := r.db.Table("bookings").
		Select("bookings.id, events.title as event_title, users.username as customer_name, bookings.total_amount, bookings.status, bookings.created_at").
		Joins("left join events on events.id = bookings.event_id").
		Joins("left join users on users.id = bookings.user_id").
		Order("bookings.created_at desc").
		Find(&bookings).Error
	return bookings, err
}

func (r *paymentRepository) GetAllPayments() ([]map[string]interface{}, error) {
	var payments []map[string]interface{}
	err := r.db.Table("payments").
		Select("payments.id, bookings.id as booking_id, payments.transaction_id, payments.amount, payments.payment_method, payments.status, payments.created_at").
		Joins("left join bookings on bookings.id = payments.booking_id").
		Order("payments.created_at desc").
		Find(&payments).Error
	return payments, err
}
func (r *paymentRepository) GetBookingsByUserID(userID string) ([]payment.MyBookingRes, error) {
	var finalBookings []payment.MyBookingRes

	err := r.db.Table("bookings").
		Select("bookings.id, events.title as concert, events.event_date as date, events.event_time as time, locations.name as venue, bookings.total_amount, bookings.status, events.image_url").
		Joins("join events on events.id = bookings.event_id").
		Joins("join locations on locations.id = events.location_id").
		Where("bookings.user_id = ?", userID).
		Order("bookings.created_at desc").
		Find(&finalBookings).Error

	if err != nil {
		return nil, err
	}

	for i := range finalBookings {
		b := &finalBookings[i]
		b.TotalPriceText = fmt.Sprintf("฿%s", formatPrice(b.TotalAmount))

		// Determine status based on event date
		now := time.Now()
		eventDate, err := time.Parse("2006-01-02", b.Date[:10]) // date is often ISO format or YYYY-MM-DD
		if err == nil {
			if b.Status == "confirmed" {
				if eventDate.After(now) {
					b.Status = "upcoming"
				} else {
					b.Status = "completed"
				}
			}
		}

		// Fetch seats for this booking
		var seats []payment.SeatDetails
		r.db.Table("booking_seats").
			Select("booking_seats.seat_id, seats.seat_number as seat_label, seats.seat_type as section, booking_seats.qr_payload, booking_seats.qr_uri, booking_seats.redeemed").
			Joins("join seats on seats.id = booking_seats.seat_id").
			Where("booking_seats.booking_id = ?", b.ID).
			Scan(&seats)

		b.Seats = seats
		b.Quantity = len(seats)

		// Summarize sections and seats for the top level
		if len(seats) > 0 {
			sections := []string{}
			labels := []string{}
			seenSections := make(map[string]bool)
			for _, s := range seats {
				if !seenSections[s.Section] {
					sections = append(sections, s.Section)
					seenSections[s.Section] = true
				}
				labels = append(labels, s.SeatLabel)
			}
			finalBookings[i].SectionSummary = strings.Join(sections, ", ")
			finalBookings[i].SeatSummary = strings.Join(labels, ", ")
		}
	}

	return finalBookings, nil
}

func (r *paymentRepository) GetBookingByID(bookingID string, userID string) ([]payment.MyBookingRes, error) {
	var finalBookings []payment.MyBookingRes

	err := r.db.Table("bookings").
		Select("bookings.id, events.title as concert, events.event_date as date, events.event_time as time, locations.name as venue, bookings.total_amount, bookings.status, events.image_url").
		Joins("join events on events.id = bookings.event_id").
		Joins("join locations on locations.id = events.location_id").
		Where("bookings.id = ? AND bookings.user_id = ?", bookingID, userID).
		Find(&finalBookings).Error

	if err != nil {
		return nil, err
	}

	for i := range finalBookings {
		b := &finalBookings[i]
		b.TotalPriceText = fmt.Sprintf("฿%s", formatPrice(b.TotalAmount))

		// Determine status based on event date
		now := time.Now()
		eventDate, err := time.Parse("2006-01-02", b.Date[:10])
		if err == nil {
			if b.Status == "confirmed" {
				if eventDate.After(now) {
					b.Status = "upcoming"
				} else {
					b.Status = "completed"
				}
			}
		}

		// Fetch seats
		var seats []payment.SeatDetails
		r.db.Table("booking_seats").
			Select("booking_seats.seat_id, seats.seat_number as seat_label, seats.seat_type as section, booking_seats.qr_payload, booking_seats.qr_uri, booking_seats.redeemed").
			Joins("join seats on seats.id = booking_seats.seat_id").
			Where("booking_seats.booking_id = ?", b.ID).
			Scan(&seats)

		b.Seats = seats
		b.Quantity = len(seats)

		if len(seats) > 0 {
			sections := []string{}
			labels := []string{}
			seenSections := make(map[string]bool)
			for _, s := range seats {
				if !seenSections[s.Section] {
					sections = append(sections, s.Section)
					seenSections[s.Section] = true
				}
				labels = append(labels, s.SeatLabel)
			}
			b.SectionSummary = strings.Join(sections, ", ")
			b.SeatSummary = strings.Join(labels, ", ")
		}
	}

	return finalBookings, nil
}

func (r *paymentRepository) UpdateBookingSeatQR(bookingID, seatID, qrPayload, qrUri string) error {
	return r.db.Table("booking_seats").
		Where("booking_id = ? AND seat_id = ?", bookingID, seatID).
		Updates(map[string]interface{}{
			"qr_payload": qrPayload,
			"qr_uri":     qrUri,
		}).Error
}

func (r *paymentRepository) RedeemTicket(bookingID, seatID string) error {
	result := r.db.Table("booking_seats").
		Where("booking_id = ? AND seat_id = ? AND redeemed = ?", bookingID, seatID, false).
		Update("redeemed", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ticket not found or already redeemed")
	}
	return nil
}

func formatPrice(price int) string {
	s := fmt.Sprintf("%d", price)
	if len(s) <= 3 {
		return s
	}
	var res []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		res = append([]string{s[start:i]}, res...)
	}
	return strings.Join(res, ",")
}
