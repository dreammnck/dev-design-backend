package service

import (
	"backend/internal/payment"
	"backend/internal/payment/adapter"
	"backend/internal/payment/repository"
	"backend/internal/seats"
	"backend/internal/seats/service"
	"errors"
	"fmt"
	"time"
)

type PaymentService interface {
	Confirm(req payment.PaymentReq) (payment.PaymentRes, error)
	Process(req payment.PaymentReq) (payment.PaymentRes, error)
	Cancel(req payment.PaymentReq) (payment.PaymentRes, error)
}

type paymentService struct {
	repo       repository.PaymentRepository
	seatSvc    service.SeatService
	gatewayAda adapter.PaymentGatewayAdapter
}

func NewPaymentService(repo repository.PaymentRepository, seatSvc service.SeatService, gatewayAda adapter.PaymentGatewayAdapter) PaymentService {
	return &paymentService{repo: repo, seatSvc: seatSvc, gatewayAda: gatewayAda}
}

func (s *paymentService) Confirm(req payment.PaymentReq) (payment.PaymentRes, error) {
	if req.CustomerID == "" {
		return payment.PaymentRes{Status: "failed", Message: "Customer ID is required"}, errors.New("missing customer id")
	}

	// Verify seat ownership
	for _, seatID := range req.SeatIDs {
		seat, err := s.seatSvc.GetByID(seatID)
		if err != nil {
			return payment.PaymentRes{Status: "failed", Message: "Seat not found: " + seatID}, err
		}
		if seat.Status != seats.StatusReserved {
			return payment.PaymentRes{Status: "failed", Message: "Seat is not reserved: " + seatID}, errors.New("seat not reserved")
		}
		if seat.CustomerID == nil || *seat.CustomerID != req.CustomerID {
			return payment.PaymentRes{Status: "failed", Message: "Seat not reserved by this customer: " + seatID}, errors.New("seat ownership mismatch")
		}
	}

	_, err := s.repo.CreateBookingAndPayment(req, "DIRECT_CONFIRM", "success")
	if err != nil {
		return payment.PaymentRes{Status: "failed", Message: err.Error()}, err
	}

	// Update seat status to 'sold'
	for _, seatID := range req.SeatIDs {
		if err := s.seatSvc.UpdateStatus(seatID, seats.StatusSold); err != nil {
			fmt.Printf("Warning: Failed to update seat %s: %v\n", seatID, err)
		}
	}

	return payment.PaymentRes{Status: "success", Message: "Payment confirmed"}, nil
}

func (s *paymentService) Process(req payment.PaymentReq) (payment.PaymentRes, error) {
	// 1. Basic Validation
	if req.Amount <= 0 {
		return payment.PaymentRes{Status: "failed", Message: "Invalid amount"}, errors.New("invalid amount")
	}

	if req.CustomerID == "" {
		return payment.PaymentRes{Status: "failed", Message: "Customer ID is required"}, errors.New("missing customer id")
	}

	if req.PaymentMethod != "credit_card" {
		return payment.PaymentRes{Status: "failed", Message: fmt.Sprintf("Unsupported method: %s", req.PaymentMethod)}, errors.New("unsupported payment method")
	}

	// 1.1 Check if seats are reserved and belong to this customer
	for _, seatID := range req.SeatIDs {
		seat, err := s.seatSvc.GetByID(seatID)
		if err != nil {
			return payment.PaymentRes{Status: "failed", Message: "Seat not found: " + seatID}, err
		}
		if seat.Status != seats.StatusReserved {
			return payment.PaymentRes{Status: "failed", Message: "Seat is not reserved: " + seatID}, errors.New("seat not reserved")
		}
		if seat.CustomerID == nil || *seat.CustomerID != req.CustomerID {
			return payment.PaymentRes{Status: "failed", Message: "Seat not reserved by this customer: " + seatID}, errors.New("seat ownership mismatch")
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
	_, err = s.repo.CreateBookingAndPayment(req, gatewayRes.TransactionID, "success")
	if err != nil {
		return payment.PaymentRes{Status: "failed", Message: "Failed to save record: " + err.Error()}, err
	}

	// 4. Update seat status to 'sold'
	for _, seatID := range req.SeatIDs {
		if err := s.seatSvc.UpdateStatus(seatID, seats.StatusSold); err != nil {
			fmt.Printf("Warning: Failed to update seat %s: %v\n", seatID, err)
		}
	}

	return payment.PaymentRes{
		Status:    "success",
		PaymentID: gatewayRes.TransactionID,
		Message:   "Payment processed and seats sold successfully",
	}, nil
}

func (s *paymentService) Cancel(req payment.PaymentReq) (payment.PaymentRes, error) {
	if req.CustomerID == "" {
		return payment.PaymentRes{Status: "failed", Message: "Customer ID is required"}, errors.New("missing customer id")
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
		if seat.CustomerID == nil || *seat.CustomerID != req.CustomerID {
			return payment.PaymentRes{Status: "failed", Message: "Seat not reserved by this customer: " + seatID}, errors.New("seat ownership mismatch")
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
