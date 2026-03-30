package seats

import "time"

const (
	StatusAvailable = "available"
	StatusReserved  = "reserved"
	StatusSold       = "sold"
)

type Seat struct {
	SeatID     string     `json:"seatId" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EventID    string     `json:"eventId" gorm:"type:uuid;column:event_id"`
	SeatNumber string     `json:"seatNumber" gorm:"column:seat_number"`
	Status     string     `json:"status" gorm:"column:status"`
	CustomerID *string    `json:"customerId" gorm:"type:uuid;column:customer_id"`
	ReservedAt *time.Time `json:"reservedAt" gorm:"column:reserved_at"`
}

type ReservationMessage struct {
	SeatID      string    `json:"seatId"`
	CustomerID  string    `json:"customerId"`
	SubmittedAt time.Time `json:"submittedAt"`
}
