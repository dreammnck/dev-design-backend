package payment

type PaymentReq struct {
	SeatIDs       []string               `json:"seatIds"`
	BookingID     string                 `json:"bookingId,omitempty"`
	CustomerID    string                 `json:"customerId"`
	PaymentMethod string                 `json:"paymentMethod"` // e.g., "credit_card", "qr_code"
	Amount        float64                `json:"amount"`
	Details       map[string]interface{} `json:"details"` // Method-specific info
}

type PaymentRes struct {
	Status    string `json:"status"`
	PaymentID string `json:"paymentId,omitempty"`
	Message   string `json:"message"`
}
