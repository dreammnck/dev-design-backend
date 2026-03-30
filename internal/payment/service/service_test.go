package service

import (
	"backend/internal/payment"
	"backend/internal/payment/adapter"
	"backend/internal/seats"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) CreateBookingAndPayment(req payment.PaymentReq, transactionID string, status string) (string, error) {
	args := m.Called(req, transactionID, status)
	return args.String(0), args.Error(1)
}

type MockSeatSvc struct {
	mock.Mock
}

func (m *MockSeatSvc) GetByEventID(id string) ([]seats.Seat, error) {
	args := m.Called(id)
	return args.Get(0).([]seats.Seat), args.Error(1)
}

func (m *MockSeatSvc) GetByID(id string) (seats.Seat, error) {
	args := m.Called(id)
	return args.Get(0).(seats.Seat), args.Error(1)
}

func (m *MockSeatSvc) ReserveSeat(id string, customerID string) error {
	args := m.Called(id, customerID)
	return args.Error(0)
}

func (m *MockSeatSvc) UpdateStatus(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockSeatSvc) ClearExpiredReservations(timeout time.Duration) error {
	args := m.Called(timeout)
	return args.Error(0)
}

type MockGatewayAdapter struct {
	mock.Mock
}

func (m *MockGatewayAdapter) ProcessPayment(req adapter.GatewayPaymentReq) (*adapter.GatewayPaymentRes, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*adapter.GatewayPaymentRes), args.Error(1)
}

func TestConfirm_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	seatSvc := new(MockSeatSvc)
	svc := NewPaymentService(repo, seatSvc, nil)
	custID := "cust-1"

	req := payment.PaymentReq{
		CustomerID: custID,
		SeatIDs:    []string{"seat-1"},
		Amount:     100,
	}

	seatSvc.On("GetByID", "seat-1").Return(seats.Seat{SeatID: "seat-1", Status: seats.StatusReserved, CustomerID: &custID}, nil)
	repo.On("CreateBookingAndPayment", req, "DIRECT_CONFIRM", "success").Return("book-1", nil)
	seatSvc.On("UpdateStatus", "seat-1", seats.StatusSold).Return(nil)

	res, err := svc.Confirm(req)

	assert.NoError(t, err)
	assert.Equal(t, "success", res.Status)
	repo.AssertExpectations(t)
	seatSvc.AssertExpectations(t)
}

func TestCancel_Success(t *testing.T) {
	seatSvc := new(MockSeatSvc)
	svc := NewPaymentService(nil, seatSvc, nil)
	custID := "cust-1"

	req := payment.PaymentReq{
		CustomerID: custID,
		SeatIDs:    []string{"seat-1"},
	}

	seatSvc.On("GetByID", "seat-1").Return(seats.Seat{SeatID: "seat-1", Status: seats.StatusReserved, CustomerID: &custID}, nil)
	seatSvc.On("UpdateStatus", "seat-1", seats.StatusAvailable).Return(nil)

	res, err := svc.Cancel(req)

	assert.NoError(t, err)
	assert.Equal(t, "cancelled", res.Status)
	seatSvc.AssertExpectations(t)
}

func TestCancel_MissingCustomerID(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil)
	req := payment.PaymentReq{
		SeatIDs: []string{"seat-1"},
	}

	res, err := svc.Cancel(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "Customer ID is required")
}

func TestCancel_OwnershipMismatch(t *testing.T) {
	seatSvc := new(MockSeatSvc)
	svc := NewPaymentService(nil, seatSvc, nil)
	ownerID := "cust-owner"

	req := payment.PaymentReq{
		CustomerID: "cust-other",
		SeatIDs:    []string{"seat-1"},
	}

	seatSvc.On("GetByID", "seat-1").Return(seats.Seat{SeatID: "seat-1", Status: seats.StatusReserved, CustomerID: &ownerID}, nil)

	res, err := svc.Cancel(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "Seat not reserved by this customer")
}

func TestProcessPayment_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	seatSvc := new(MockSeatSvc)
	adapterMock := new(MockGatewayAdapter)
	svc := NewPaymentService(repo, seatSvc, adapterMock)

	req := payment.PaymentReq{
		SeatIDs:       []string{"seat-1"},
		CustomerID:    "cust-1",
		Amount:        100,
		PaymentMethod: "credit_card",
		Details: map[string]interface{}{
			"card_number":  "1234",
			"expiry_month": 12.0,
			"expiry_year":  25.0,
			"cvv":          "123",
		},
	}

	custID := "cust-1"
	seatSvc.On("GetByID", "seat-1").Return(seats.Seat{SeatID: "seat-1", Status: seats.StatusReserved, CustomerID: &custID}, nil)
	adapterMock.On("ProcessPayment", mock.Anything).Return(&adapter.GatewayPaymentRes{
		Status:        "success",
		TransactionID: "tx-123",
	}, nil)
	repo.On("CreateBookingAndPayment", req, "tx-123", "success").Return("book-123", nil)
	seatSvc.On("UpdateStatus", "seat-1", seats.StatusSold).Return(nil)

	res, err := svc.Process(req)

	assert.NoError(t, err)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "tx-123", res.PaymentID)
	seatSvc.AssertExpectations(t)
	adapterMock.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestProcessPayment_InvalidAmount(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil)
	req := payment.PaymentReq{Amount: 0}

	res, err := svc.Process(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "Invalid amount")
}

func TestProcessPayment_MissingCustomerID(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil)
	req := payment.PaymentReq{Amount: 100}

	res, err := svc.Process(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "Customer ID is required")
}
