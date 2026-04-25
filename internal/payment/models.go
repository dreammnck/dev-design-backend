package payment

type PaymentReq struct {
	SeatIDs       []string               `json:"seatIds"`
	BookingID     string                 `json:"bookingId,omitempty"`
	UserID        string                 `json:"-"`             // Injected from JWT claims, not from request body
	PaymentMethod string                 `json:"paymentMethod"` // e.g., "credit_card", "qr_code"
	Amount        float64                `json:"amount"`
	Details       map[string]interface{} `json:"details"` // Method-specific info
}

type PaymentRes struct {
	Status    string `json:"status"`
	PaymentID string `json:"paymentId,omitempty"`
	Message   string `json:"message"`
}

type MyBookingRes struct {
	ID             string        `json:"id"`
	Concert        string        `json:"concert" gorm:"column:concert"`
	Date           string        `json:"date" gorm:"column:date"`
	Time           string        `json:"time" gorm:"column:time"`
	Venue          string        `json:"venue" gorm:"column:venue"`
	SectionSummary string        `json:"section" gorm:"-"`
	SeatSummary    string        `json:"seat" gorm:"-"`
	Quantity       int           `json:"quantity" gorm:"-"`
	TotalPriceText string        `json:"totalPriceText" gorm:"-"`
	TotalAmount    int           `json:"-" gorm:"column:total_amount"`
	Status         string        `json:"status"`
	ImageURL       string        `json:"imageUrl" gorm:"column:image_url"`
	Seats          []SeatDetails `json:"seats" gorm:"-"`
}

type SeatDetails struct {
	SeatID    string `json:"seatId"    gorm:"column:seat_id"`
	SeatLabel string `json:"seatLabel" gorm:"column:seat_label"`
	Section   string `json:"section"   gorm:"column:section"`
	QRPayload string `json:"qrPayload" gorm:"column:qr_payload"`
	QRUri     string `json:"qrUri"     gorm:"column:qr_uri"`
	Redeemed  bool   `json:"redeemed"  gorm:"column:redeemed"`
}
