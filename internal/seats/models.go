package seats

import "time"

const (
	StatusAvailable = "available"
	StatusReserved  = "reserved"
	StatusSold      = "sold"
)

type Seat struct {
	SeatID     string     `json:"seatId" gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EventID    string     `json:"eventId" gorm:"type:uuid;column:event_id"`
	SeatNumber string     `json:"seatNumber" gorm:"column:seat_number"`
	Status     string     `json:"status" gorm:"column:status"`
	Price      int        `json:"price" gorm:"column:price"`
	SeatType   string     `json:"seatType" gorm:"column:seat_type"`

	CustomerID *string    `json:"userId" gorm:"type:uuid;column:user_id"`
	ReservedAt *time.Time `json:"reservedAt" gorm:"column:reserved_at"`
}

type ReservationMessage struct {
	SeatID      string    `json:"seatId"`
	CustomerID  string    `json:"userId"`
	SubmittedAt time.Time `json:"submittedAt"`
}
