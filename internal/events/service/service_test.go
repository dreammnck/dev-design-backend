package service

import (
	"backend/internal/events"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository
type MockEventRepo struct {
	mock.Mock
}

func (m *MockEventRepo) GetBanner() ([]events.Event, error) {
	args := m.Called()
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetAll() ([]events.Event, error) {
	args := m.Called()
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetRecommend() ([]events.Event, error) {
	args := m.Called()
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetComingSoon() ([]events.Event, error) {
	args := m.Called()
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetByID(id string) (events.Event, error) {
	args := m.Called(id)
	return args.Get(0).(events.Event), args.Error(1)
}

func TestGetBanner_Success(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)
	expected := []events.Event{{ID: "1", Title: "Event 1"}}

	repo.On("GetBanner").Return(expected, nil)

	res, err := svc.GetBanner()

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetByID_Success(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)
	expected := events.Event{ID: "1", Title: "Event 1"}

	repo.On("GetByID", "1").Return(expected, nil)

	res, err := svc.GetByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetByID_Error(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)

	repo.On("GetByID", "invalid").Return(events.Event{}, errors.New("not found"))

	_, err := svc.GetByID("invalid")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}
