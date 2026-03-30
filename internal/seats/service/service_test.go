package service

import (
	"backend/internal/seats"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeatRepository
type MockSeatRepo struct {
	mock.Mock
}

func (m *MockSeatRepo) GetByEventID(id string) ([]seats.Seat, error) {
	args := m.Called(id)
	return args.Get(0).([]seats.Seat), args.Error(1)
}

func (m *MockSeatRepo) GetByID(id string) (seats.Seat, error) {
	args := m.Called(id)
	return args.Get(0).(seats.Seat), args.Error(1)
}

func (m *MockSeatRepo) UpdateStatus(id string, status string, reservedAt *time.Time, customerID *string) error {
	args := m.Called(id, status, reservedAt, customerID)
	return args.Error(0)
}

func (m *MockSeatRepo) ClearExpiredReservations(timeout time.Duration) error {
	args := m.Called(timeout)
	return args.Error(0)
}

func TestReserveSeat_Success_Available(t *testing.T) {
	repo := new(MockSeatRepo)
	svc := NewSeatService(repo)
	id := "seat-1"

	repo.On("GetByID", id).Return(seats.Seat{SeatID: id, Status: seats.StatusAvailable}, nil)
	repo.On("UpdateStatus", id, seats.StatusReserved, mock.Anything, mock.Anything).Return(nil)

	err := svc.ReserveSeat(id, "customer-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestReserveSeat_Success_Expired(t *testing.T) {
	repo := new(MockSeatRepo)
	svc := NewSeatService(repo)
	id := "seat-1"
	expiredTime := time.Now().Add(-10 * time.Minute)

	repo.On("GetByID", id).Return(seats.Seat{
		SeatID:     id,
		Status:     seats.StatusReserved,
		ReservedAt: &expiredTime,
	}, nil)
	repo.On("UpdateStatus", id, seats.StatusReserved, mock.Anything, mock.Anything).Return(nil)

	err := svc.ReserveSeat(id, "customer-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestReserveSeat_Fail_AlreadyReserved(t *testing.T) {
	repo := new(MockSeatRepo)
	svc := NewSeatService(repo)
	id := "seat-1"
	recentTime := time.Now().Add(-1 * time.Minute)

	repo.On("GetByID", id).Return(seats.Seat{
		SeatID:     id,
		Status:     seats.StatusReserved,
		ReservedAt: &recentTime,
	}, nil)

	err := svc.ReserveSeat(id, "customer-1")

	assert.Error(t, err)
	assert.Equal(t, "seat is not available", err.Error())
}

func TestReserveSeat_Fail_NotFound(t *testing.T) {
	repo := new(MockSeatRepo)
	svc := NewSeatService(repo)
	id := "non-existent"

	repo.On("GetByID", id).Return(seats.Seat{}, errors.New("not found"))

	err := svc.ReserveSeat(id, "customer-1")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}
