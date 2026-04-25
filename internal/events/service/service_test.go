package service

import (
	"backend/internal/events"
	"backend/internal/events/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository — implements all methods of repository.EventRepository
type MockEventRepo struct {
	mock.Mock
}

func (m *MockEventRepo) GetBanner(userID string) ([]events.Event, error) {
	args := m.Called(userID)
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetAll(userID string, search string) ([]events.Event, error) {
	args := m.Called(userID, search)
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetRecommend(userID string) ([]events.Event, error) {
	args := m.Called(userID)
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetComingSoon(userID string) ([]events.Event, error) {
	args := m.Called(userID)
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetByID(id string, userID string) (events.Event, error) {
	args := m.Called(id, userID)
	return args.Get(0).(events.Event), args.Error(1)
}

func (m *MockEventRepo) ToggleFavorite(userID string, eventID string) (bool, error) {
	args := m.Called(userID, eventID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEventRepo) GetFavoritedEvents(userID string) ([]events.Event, error) {
	args := m.Called(userID)
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) Create(orgID string, req events.EventCreateReq, imageUrl string) (events.Event, error) {
	args := m.Called(orgID, req, imageUrl)
	return args.Get(0).(events.Event), args.Error(1)
}

func (m *MockEventRepo) Update(id string, orgID string, req events.EventUpdateReq) (events.Event, error) {
	args := m.Called(id, orgID, req)
	return args.Get(0).(events.Event), args.Error(1)
}

func (m *MockEventRepo) SubmitForReview(id string, orgID string) error {
	args := m.Called(id, orgID)
	return args.Error(0)
}

func (m *MockEventRepo) GetByOrganization(orgID string) ([]events.Event, error) {
	args := m.Called(orgID)
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) GetBookingSummary(eventID string, orgID string) ([]repository.BookingSummaryRow, error) {
	args := m.Called(eventID, orgID)
	return args.Get(0).([]repository.BookingSummaryRow), args.Error(1)
}

func (m *MockEventRepo) AddSeats(eventID string, orgID string, seats []events.SeatInput) error {
	args := m.Called(eventID, orgID, seats)
	return args.Error(0)
}

func (m *MockEventRepo) Review(id string, status events.PublishStatus, rejectReason string) error {
	args := m.Called(id, status, rejectReason)
	return args.Error(0)
}

func (m *MockEventRepo) GetPendingEvents() ([]events.Event, error) {
	args := m.Called()
	return args.Get(0).([]events.Event), args.Error(1)
}

func (m *MockEventRepo) AdminEditEvent(id string, req events.AdminEditEventReq) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockEventRepo) CreatePayout(orgID string, req events.PayoutReq) (events.Payout, error) {
	args := m.Called(orgID, req)
	return args.Get(0).(events.Payout), args.Error(1)
}

func (m *MockEventRepo) GetPayouts(orgID string) ([]events.Payout, error) {
	args := m.Called(orgID)
	return args.Get(0).([]events.Payout), args.Error(1)
}

func (m *MockEventRepo) ProcessPayout(payoutID string, status events.PayoutStatus, note string) error {
	args := m.Called(payoutID, status, note)
	return args.Error(0)
}

func (m *MockEventRepo) OrganizerGetOverallSummary(orgID string) (map[string]interface{}, error) {
	args := m.Called(orgID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) OrganizerGetEventSales(eventID string, days int) ([]map[string]interface{}, error) {
	args := m.Called(eventID, days)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) OrganizerGetEventSeatsSummary(eventID string) (map[string]interface{}, error) {
	args := m.Called(eventID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) OrganizerGetEventsCompare(orgID string) ([]map[string]interface{}, error) {
	args := m.Called(orgID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) GetAllAdminEvents() ([]map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) GetAllPayouts() ([]map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) GetTotalIncomeByOrgID(orgID string) (int, error) {
	args := m.Called(orgID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepo) GetAllLocations() ([]events.Location, error) {
	args := m.Called()
	return args.Get(0).([]events.Location), args.Error(1)
}

func (m *MockEventRepo) AdminGetOverallSummary() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) AdminGetOrderStats(days int) ([]map[string]interface{}, error) {
	args := m.Called(days)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) AdminGetOrdersSummaryStatus() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockEventRepo) AdminGetTopSellingEvents(limit int) ([]map[string]interface{}, error) {
	args := m.Called(limit)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestGetBanner_Success(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)
	expected := []events.Event{{ID: "1", Title: "Event 1"}}

	repo.On("GetBanner", "").Return(expected, nil)

	res, err := svc.GetBanner("")

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetByID_Success(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)
	expected := events.Event{ID: "1", Title: "Event 1"}

	repo.On("GetByID", "1", "").Return(expected, nil)

	res, err := svc.GetByID("1", "")

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetByID_Error(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)

	repo.On("GetByID", "invalid", "").Return(events.Event{}, errors.New("not found"))

	_, err := svc.GetByID("invalid", "")

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}

func TestReviewEvent_RejectWithoutReason(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)

	err := svc.ReviewEvent("event-1", events.ReviewReq{
		Status:       events.PublishStatusRejected,
		RejectReason: "",
	})

	assert.Error(t, err)
	assert.Equal(t, ErrMissingReason, err)
}

func TestReviewEvent_Approve(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)

	repo.On("Review", "event-1", events.PublishStatusApproved, "").Return(nil)

	err := svc.ReviewEvent("event-1", events.ReviewReq{
		Status: events.PublishStatusApproved,
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestReviewEvent_InvalidStatus(t *testing.T) {
	repo := new(MockEventRepo)
	svc := NewEventService(repo)

	err := svc.ReviewEvent("event-1", events.ReviewReq{
		Status: "unknown",
	})

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidReview, err)
}
