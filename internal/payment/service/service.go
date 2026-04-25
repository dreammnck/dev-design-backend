package service

import (
	auth "backend/internal/auth/service"
	"backend/internal/payment"
	"backend/internal/payment/adapter"
	"backend/internal/payment/repository"
	"backend/internal/seats"
	"backend/internal/seats/service"
	"backend/pkg/notification"
	"backend/pkg/storage"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type PaymentService interface {
	Confirm(req payment.PaymentReq) (payment.PaymentRes, error)
	Process(req payment.PaymentReq) (payment.PaymentRes, error)
	Cancel(req payment.PaymentReq) (payment.PaymentRes, error)
	GetAllBookings() ([]map[string]interface{}, error)
	GetMyBookings(userID string) ([]payment.MyBookingRes, error)
	GetAllPayments() ([]map[string]interface{}, error)
	GetTicketByID(bookingID string, userID string) ([]payment.MyBookingRes, error)
	RedeemTicket(qrPayload string) error
}

type paymentService struct {
	repo            repository.PaymentRepository
	seatSvc         service.SeatService
	gatewayAda      adapter.PaymentGatewayAdapter
	authSvc         auth.AuthService
	notificationSvc notification.NotificationService
}

func NewPaymentService(
	repo repository.PaymentRepository,
	seatSvc service.SeatService,
	gatewayAda adapter.PaymentGatewayAdapter,
	authSvc auth.AuthService,
	notificationSvc notification.NotificationService,
) PaymentService {
	return &paymentService{
		repo:            repo,
		seatSvc:         seatSvc,
		gatewayAda:      gatewayAda,
		authSvc:         authSvc,
		notificationSvc: notificationSvc,
	}
}

func (s *paymentService) Confirm(req payment.PaymentReq) (payment.PaymentRes, error) {
	if req.UserID == "" {
		return payment.PaymentRes{Status: "failed", Message: "User ID is required"}, errors.New("missing user id")
	}

	// Verify seat is reserved (Admin can confirm any reserved seat)
	for _, seatID := range req.SeatIDs {
		seat, err := s.seatSvc.GetByID(seatID)
		if err != nil {
			return payment.PaymentRes{Status: "failed", Message: "Seat not found: " + seatID}, err
		}
		if seat.Status != seats.StatusReserved {
			return payment.PaymentRes{Status: "failed", Message: "Seat is not reserved: " + seatID}, errors.New("seat not reserved")
		}
	}

	bookingID, err := s.repo.CreateBookingAndPayment(req, "DIRECT_CONFIRM", "success")
	if err != nil {
		return payment.PaymentRes{Status: "failed", Message: err.Error()}, err
	}

	// Update seat status to 'sold' and Generate QR
	for _, seatID := range req.SeatIDs {
		if err := s.seatSvc.UpdateStatus(seatID, seats.StatusSold); err != nil {
			fmt.Printf("Warning: Failed to update seat %s: %v\n", seatID, err)
		}
		// Generate and Save QR
		go s.processSeatQR(bookingID, seatID)
	}

	// Trigger Notification
	go s.triggerTicketNotification(req.UserID, req.SeatIDs)

	return payment.PaymentRes{Status: "success", Message: "Payment confirmed"}, nil
}

func (s *paymentService) Process(req payment.PaymentReq) (payment.PaymentRes, error) {
	// 1. Basic Validation
	if req.Amount <= 0 {
		return payment.PaymentRes{Status: "failed", Message: "Invalid amount"}, errors.New("invalid amount")
	}

	if req.UserID == "" {
		return payment.PaymentRes{Status: "failed", Message: "User ID is required"}, errors.New("missing user id")
	}

	if req.PaymentMethod != "credit_card" {
		return payment.PaymentRes{Status: "failed", Message: fmt.Sprintf("Unsupported method: %s", req.PaymentMethod)}, errors.New("unsupported payment method")
	}

	// 1.1 Check if seats are reserved and belong to this user
	for _, seatID := range req.SeatIDs {
		seat, err := s.seatSvc.GetByID(seatID)
		if err != nil {
			return payment.PaymentRes{Status: "failed", Message: "Seat not found: " + seatID}, err
		}
		if seat.Status != seats.StatusReserved {
			return payment.PaymentRes{Status: "failed", Message: "Seat is not reserved: " + seatID}, errors.New("seat not reserved")
		}
		if seat.CustomerID == nil || *seat.CustomerID != req.UserID {
			return payment.PaymentRes{Status: "failed", Message: "Seat not reserved by this user: " + seatID}, errors.New("seat ownership mismatch")
		}
	}

	// Safe extraction
	expiryMonth, _ := req.Details["expiry_month"].(float64)
	expiryYear, _ := req.Details["expiry_year"].(float64)
	cardNumber, _ := req.Details["card_number"].(string)
	cvv, _ := req.Details["cvv"].(string)

	gatewayReq := adapter.GatewayPaymentReq{
		Amount:      req.Amount,
		Currency:    "THB",
		CardNumber:  cardNumber,
		ExpiryMonth: int(expiryMonth),
		ExpiryYear:  int(expiryYear),
		CVV:         cvv,
		MerchantID:  "MERC-001",
		OrderID:     req.BookingID,
	}

	if gatewayReq.OrderID == "" {
		gatewayReq.OrderID = fmt.Sprintf("ORD-%s-001", time.Now().Format("20060102"))
	}

	// 2. Call External Gateway via Adapter
	// เราเปิดใช้งานตัวนี้อีกครั้งเพื่อให้ระบบสมบูรณ์
	gatewayRes, err := s.gatewayAda.ProcessPayment(gatewayReq)
	if err != nil {
		return payment.PaymentRes{Status: "failed", Message: fmt.Sprintf("Gateway communication error: %v", err)}, err
	}

	if gatewayRes.Status != "success" {
		return payment.PaymentRes{
			Status:  "failed",
			Message: gatewayRes.Message,
		}, fmt.Errorf("gateway error: %s", gatewayRes.Message)
	}

	// 3. Update Local DB: Create Booking and Payment record
	bookingID, err := s.repo.CreateBookingAndPayment(req, gatewayRes.TransactionID, "success")
	if err != nil {
		return payment.PaymentRes{Status: "failed", Message: "Failed to save record: " + err.Error()}, err
	}

	// 4. Update seat status to 'sold'
	for _, seatID := range req.SeatIDs {
		if err := s.seatSvc.UpdateStatus(seatID, seats.StatusSold); err != nil {
			fmt.Printf("Warning: Failed to update seat %s: %v\n", seatID, err)
		}
		// Generate and Save QR
		go s.processSeatQR(bookingID, seatID)
	}

	// Trigger Notification in background (One email with all tickets)
	go s.triggerTicketNotification(req.UserID, req.SeatIDs)

	return payment.PaymentRes{
		Status:    "success",
		PaymentID: gatewayRes.TransactionID,
		Message:   "Payment processed and seats sold successfully",
	}, nil
}

func (s *paymentService) triggerTicketNotification(userID string, seatIDs []string) error {
	user, err := s.authSvc.GetUser(userID)
	if err != nil {
		fmt.Printf("Error: Failed to fetch user %s for notification: %v\n", userID, err)
		return err
	}

	var tickets []notification.TicketData
	for _, seatID := range seatIDs {
		seat, err := s.seatSvc.GetByID(seatID)
		if err != nil {
			fmt.Printf("Error: Failed to fetch seat %s for notification: %v\n", seatID, err)
			continue
		}
		tickets = append(tickets, notification.TicketData{
			EventID:    seat.EventID,
			SeatID:     seatID,
			SeatNumber: seat.SeatNumber,
		})
	}

	if len(tickets) > 0 {
		if err := s.notificationSvc.SendTicketsEmail(user.Email, tickets); err != nil {
			fmt.Printf("Error: Failed to send batch tickets email to %s: %v\n", user.Email, err)
			return err
		}
	}
	return nil
}

func (s *paymentService) Cancel(req payment.PaymentReq) (payment.PaymentRes, error) {
	if req.UserID == "" {
		return payment.PaymentRes{Status: "failed", Message: "User ID is required"}, errors.New("missing user id")
	}

	// Verify seat ownership before cancelling
	for _, seatID := range req.SeatIDs {
		seat, err := s.seatSvc.GetByID(seatID)
		if err != nil {
			return payment.PaymentRes{Status: "failed", Message: "Seat not found: " + seatID}, err
		}
		if seat.Status != seats.StatusReserved {
			return payment.PaymentRes{Status: "failed", Message: "Seat is not reserved: " + seatID}, errors.New("seat not reserved")
		}
		if seat.CustomerID == nil || *seat.CustomerID != req.UserID {
			return payment.PaymentRes{Status: "failed", Message: "Seat not reserved by this user: " + seatID}, errors.New("seat ownership mismatch")
		}
	}

	// Revert seat status to available
	for _, seatID := range req.SeatIDs {
		_ = s.seatSvc.UpdateStatus(seatID, seats.StatusAvailable)
	}

	return payment.PaymentRes{
		Status:  "cancelled",
		Message: "Payment cancelled and seats released",
	}, nil
}

func (s *paymentService) GetAllBookings() ([]map[string]interface{}, error) {
	return s.repo.GetAllBookings()
}

func (s *paymentService) GetMyBookings(userID string) ([]payment.MyBookingRes, error) {
	return s.repo.GetBookingsByUserID(userID)
}

func (s *paymentService) GetAllPayments() ([]map[string]interface{}, error) {
	return s.repo.GetAllPayments()
}

func (s *paymentService) GetTicketByID(bookingID string, userID string) ([]payment.MyBookingRes, error) {
	return s.repo.GetBookingByID(bookingID, userID)
}

func (s *paymentService) RedeemTicket(qrPayload string) error {
	var p struct {
		BookingID string `json:"bookingId"`
		SeatID    string `json:"seatId"`
	}

	if err := json.Unmarshal([]byte(qrPayload), &p); err != nil {
		return fmt.Errorf("invalid QR payload format")
	}

	if p.BookingID == "" || p.SeatID == "" {
		return fmt.Errorf("missing bookingId or seatId in QR")
	}

	return s.repo.RedeemTicket(p.BookingID, p.SeatID)
}

func (s *paymentService) processSeatQR(bookingID, seatID string) {
	seat, err := s.seatSvc.GetByID(seatID)
	if err != nil {
		fmt.Printf("Error: Failed to fetch seat %s for QR: %v\n", seatID, err)
		return
	}

	// Create Payload
	payloadObj := map[string]string{
		"bookingId": bookingID,
		"seatLabel": seat.SeatNumber,
		"seatId":    seatID,
	}
	payloadJSON, _ := json.Marshal(payloadObj)
	qrPayload := string(payloadJSON)

	// Generate QR Image
	qrBytes, err := notification.GenerateQRCode(qrPayload)
	if err != nil {
		fmt.Printf("Error: Failed to generate QR for seat %s: %v\n", seatID, err)
		return
	}

	// Upload to GCS
	filename := fmt.Sprintf("tickets/%s/%s.png", bookingID, seatID)
	qrUri, err := storage.UploadBytes(filename, qrBytes, "image/png")
	if err != nil {
		fmt.Printf("Error: Failed to upload QR for seat %s: %v\n", seatID, err)
		return
	}

	// Update DB
	if err := s.repo.UpdateBookingSeatQR(bookingID, seatID, qrPayload, qrUri); err != nil {
		fmt.Printf("Error: Failed to update DB for seat %s: %v\n", seatID, err)
	}
}
