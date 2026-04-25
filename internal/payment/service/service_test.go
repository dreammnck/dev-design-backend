package service

import (
	authModel "backend/internal/auth"
	"backend/internal/payment"
	"backend/internal/payment/adapter"
	"backend/internal/seats"
	"backend/pkg/notification"
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

func (m *MockPaymentRepo) GetAllBookings() ([]map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockPaymentRepo) GetAllPayments() ([]map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockPaymentRepo) UpdateBookingSeatQR(bookingID, seatID, qrPayload, qrUri string) error {
	args := m.Called(bookingID, seatID, qrPayload, qrUri)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetBookingsByUserID(userID string) ([]payment.MyBookingRes, error) {
	args := m.Called(userID)
	return args.Get(0).([]payment.MyBookingRes), args.Error(1)
}
func (m *MockPaymentRepo) GetBookingByID(bookingID string, userID string) ([]payment.MyBookingRes, error) {
	args := m.Called(bookingID, userID)
	return args.Get(0).([]payment.MyBookingRes), args.Error(1)
}

func (m *MockPaymentRepo) RedeemTicket(bookingID, seatID string) error {
	args := m.Called(bookingID, seatID)
	return args.Error(0)
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

type MockAuthSvc struct {
	mock.Mock
}

func (m *MockAuthSvc) GetUser(id string) (*authModel.UserInfo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.UserInfo), args.Error(1)
}

func (m *MockAuthSvc) GetAllUsers() ([]authModel.UserInfo, error) {
	args := m.Called()
	return args.Get(0).([]authModel.UserInfo), args.Error(1)
}

func (m *MockAuthSvc) Login(req authModel.LoginReq) (*authModel.LoginRes, error) {
	args := m.Called(req)
	return args.Get(0).(*authModel.LoginRes), args.Error(1)
}

func (m *MockAuthSvc) AdminUpdateUser(id string, req authModel.UpdateUserReq) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockAuthSvc) Register(req authModel.RegisterReq) (*authModel.UserInfo, error) {
	args := m.Called(req)
	return args.Get(0).(*authModel.UserInfo), args.Error(1)
}

func (m *MockAuthSvc) UpdateRole(id string, req authModel.UpdateRoleReq) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockAuthSvc) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockNotificationSvc struct {
	mock.Mock
}

func (m *MockNotificationSvc) SendTicketEmail(toEmail string, eventID string, seatID string) error {
	args := m.Called(toEmail, eventID, seatID)
	return args.Error(0)
}

func (m *MockNotificationSvc) SendTicketsEmail(toEmail string, tickets []notification.TicketData) error {
	args := m.Called(toEmail, tickets)
	return args.Error(0)
}

func TestConfirm_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	seatSvc := new(MockSeatSvc)
	authSvc := new(MockAuthSvc)
	notifSvc := new(MockNotificationSvc)
	svc := NewPaymentService(repo, seatSvc, nil, authSvc, notifSvc)
	custID := "cust-1"

	req := payment.PaymentReq{
		UserID:  custID,
		SeatIDs: []string{"seat-1"},
		Amount:  100,
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
	svc := NewPaymentService(nil, seatSvc, nil, nil, nil)
	custID := "cust-1"

	req := payment.PaymentReq{
		UserID:  custID,
		SeatIDs: []string{"seat-1"},
	}

	seatSvc.On("GetByID", "seat-1").Return(seats.Seat{SeatID: "seat-1", Status: seats.StatusReserved, CustomerID: &custID}, nil)
	seatSvc.On("UpdateStatus", "seat-1", seats.StatusAvailable).Return(nil)

	res, err := svc.Cancel(req)

	assert.NoError(t, err)
	assert.Equal(t, "cancelled", res.Status)
	seatSvc.AssertExpectations(t)
}

func TestCancel_MissingUserID(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil)
	req := payment.PaymentReq{
		SeatIDs: []string{"seat-1"},
	}

	res, err := svc.Cancel(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "User ID is required")
}

func TestCancel_OwnershipMismatch(t *testing.T) {
	seatSvc := new(MockSeatSvc)
	svc := NewPaymentService(nil, seatSvc, nil, nil, nil)
	ownerID := "cust-owner"

	req := payment.PaymentReq{
		UserID:  "cust-other",
		SeatIDs: []string{"seat-1"},
	}

	seatSvc.On("GetByID", "seat-1").Return(seats.Seat{SeatID: "seat-1", Status: seats.StatusReserved, CustomerID: &ownerID}, nil)

	res, err := svc.Cancel(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "Seat not reserved by this user")
}

func TestProcessPayment_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	seatSvc := new(MockSeatSvc)
	adapterMock := new(MockGatewayAdapter)
	authSvc := new(MockAuthSvc)
	notifSvc := new(MockNotificationSvc)
	svc := NewPaymentService(repo, seatSvc, adapterMock, authSvc, notifSvc)

	req := payment.PaymentReq{
		SeatIDs:       []string{"seat-1"},
		UserID:        "cust-1",
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
	svc := NewPaymentService(nil, nil, nil, nil, nil)
	req := payment.PaymentReq{Amount: 0}

	res, err := svc.Process(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "Invalid amount")
}

func TestProcessPayment_MissingUserID(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil)
	req := payment.PaymentReq{Amount: 100}

	res, err := svc.Process(req)

	assert.Error(t, err)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Message, "User ID is required")
}

func TestGetMyBookings_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	svc := NewPaymentService(repo, nil, nil, nil, nil)
	userID := "user-123"

	expected := []payment.MyBookingRes{
		{ID: "bk-001", Concert: "Event title", Status: "confirmed"},
	}

	repo.On("GetBookingsByUserID", userID).Return(expected, nil)

	res, err := svc.GetMyBookings(userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repo.AssertExpectations(t)
}
