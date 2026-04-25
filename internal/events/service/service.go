package service

import (
	"backend/internal/events"
	"backend/internal/events/repository"
	"errors"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrNotOwner      = errors.New("you do not own this event")
	ErrInvalidReview = errors.New("invalid review status: must be 'approved' or 'rejected'")
	ErrMissingReason = errors.New("reject_reason is required when rejecting an event")
)

type EventService interface {
	// Public
	GetBanner(userID string) ([]events.Event, error)
	GetAll(userID string, search string) ([]events.Event, error)
	GetRecommend(userID string) ([]events.Event, error)
	GetComingSoon(userID string) ([]events.Event, error)
	GetByID(id string, userID string) (events.Event, error)

	// Favorites
	ToggleFavorite(userID string, eventID string) (bool, error)
	GetFavoritedEvents(userID string) ([]events.Event, error)

	// Organization
	CreateEvent(orgID string, req events.EventCreateReq, imageUrl string) (events.Event, error)
	UpdateEvent(id string, orgID string, req events.EventUpdateReq) (events.Event, error)
	SubmitForReview(id string, orgID string) error
	GetMyEvents(orgID string) ([]events.Event, error)
	GetBookingSummary(eventID string, orgID string) ([]repository.BookingSummaryRow, error)
	AddSeats(eventID string, orgID string, req events.SeatBatchCreateReq) error
	RequestPayout(orgID string, req events.PayoutReq) (events.Payout, error)
	GetPayouts(orgID string) (map[string]interface{}, error)

	// Dashboard
	GetOrganizerOverallSummary(orgID string) (map[string]interface{}, error)
	GetOrganizerEventSales(eventID string, days int) ([]map[string]interface{}, error)
	GetOrganizerEventSeatsSummary(eventID string) (map[string]interface{}, error)
	GetOrganizerEventsCompare(orgID string) ([]map[string]interface{}, error)

	// Admin Dashboard
	GetAdminOverallSummary() (map[string]interface{}, error)
	GetAdminOrderStats(days int) ([]map[string]interface{}, error)
	GetAdminOrdersSummaryStatus() (map[string]interface{}, error)
	GetAdminTopSellingEvents(limit int) ([]map[string]interface{}, error)

	// Admin
	ReviewEvent(id string, req events.ReviewReq) error
	GetPendingEvents() ([]events.Event, error)
	GetAllAdminEvents() ([]map[string]interface{}, error)
	AdminEditEvent(id string, req events.AdminEditEventReq) error
	ProcessPayout(payoutID string, status events.PayoutStatus, note string) error
	GetAllPayouts() ([]map[string]interface{}, error)
	GetAllLocations() ([]events.Location, error)
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

// ── Public ────────────────────────────────────────────────────────────────────

func (s *eventService) GetBanner(userID string) ([]events.Event, error) {
	return s.repo.GetBanner(userID)
}

func (s *eventService) GetAll(userID string, search string) ([]events.Event, error) {
	return s.repo.GetAll(userID, search)
}

func (s *eventService) GetRecommend(userID string) ([]events.Event, error) {
	return s.repo.GetRecommend(userID)
}

func (s *eventService) GetComingSoon(userID string) ([]events.Event, error) {
	return s.repo.GetComingSoon(userID)
}

func (s *eventService) GetByID(id string, userID string) (events.Event, error) {
	return s.repo.GetByID(id, userID)
}

func (s *eventService) ToggleFavorite(userID string, eventID string) (bool, error) {
	return s.repo.ToggleFavorite(userID, eventID)
}

func (s *eventService) GetFavoritedEvents(userID string) ([]events.Event, error) {
	return s.repo.GetFavoritedEvents(userID)
}

// ── Organization ──────────────────────────────────────────────────────────────

func (s *eventService) CreateEvent(orgID string, req events.EventCreateReq, imageUrl string) (events.Event, error) {
	return s.repo.Create(orgID, req, imageUrl)
}

func (s *eventService) UpdateEvent(id string, orgID string, req events.EventUpdateReq) (events.Event, error) {
	return s.repo.Update(id, orgID, req)
}

func (s *eventService) SubmitForReview(id string, orgID string) error {
	return s.repo.SubmitForReview(id, orgID)
}

func (s *eventService) GetMyEvents(orgID string) ([]events.Event, error) {
	return s.repo.GetByOrganization(orgID)
}

func (s *eventService) GetBookingSummary(eventID string, orgID string) ([]repository.BookingSummaryRow, error) {
	return s.repo.GetBookingSummary(eventID, orgID)
}

func (s *eventService) AddSeats(eventID string, orgID string, req events.SeatBatchCreateReq) error {
	return s.repo.AddSeats(eventID, orgID, req.Seats)
}

func (s *eventService) RequestPayout(orgID string, req events.PayoutReq) (events.Payout, error) {
	return s.repo.CreatePayout(orgID, req)
}

func (s *eventService) GetPayouts(orgID string) (map[string]interface{}, error) {
	payouts, err := s.repo.GetPayouts(orgID)
	if err != nil {
		return nil, err
	}

	totalIncome, err := s.repo.GetTotalIncomeByOrgID(orgID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"payoutInfo":  payouts,
		"totalIncome": totalIncome,
	}, nil
}

// ── Admin ─────────────────────────────────────────────────────────────────────

func (s *eventService) ReviewEvent(id string, req events.ReviewReq) error {
	if req.Status != events.PublishStatusApproved && req.Status != events.PublishStatusRejected {
		return ErrInvalidReview
	}
	if req.Status == events.PublishStatusRejected && req.RejectReason == "" {
		return ErrMissingReason
	}
	return s.repo.Review(id, req.Status, req.RejectReason)
}

func (s *eventService) GetPendingEvents() ([]events.Event, error) {
	return s.repo.GetPendingEvents()
}

func (s *eventService) GetAllAdminEvents() ([]map[string]interface{}, error) {
	return s.repo.GetAllAdminEvents()
}

func (s *eventService) AdminEditEvent(id string, req events.AdminEditEventReq) error {
	return s.repo.AdminEditEvent(id, req)
}

func (s *eventService) ProcessPayout(payoutID string, status events.PayoutStatus, note string) error {
	return s.repo.ProcessPayout(payoutID, status, note)
}

func (s *eventService) GetAllPayouts() ([]map[string]interface{}, error) {
	return s.repo.GetAllPayouts()
}

func (s *eventService) GetAllLocations() ([]events.Location, error) {
	return s.repo.GetAllLocations()
}

func (s *eventService) GetOrganizerOverallSummary(orgID string) (map[string]interface{}, error) {
	return s.repo.OrganizerGetOverallSummary(orgID)
}

func (s *eventService) GetOrganizerEventSales(eventID string, days int) ([]map[string]interface{}, error) {
	return s.repo.OrganizerGetEventSales(eventID, days)
}

func (s *eventService) GetOrganizerEventSeatsSummary(eventID string) (map[string]interface{}, error) {
	return s.repo.OrganizerGetEventSeatsSummary(eventID)
}
func (s *eventService) GetAdminOverallSummary() (map[string]interface{}, error) {
	return s.repo.AdminGetOverallSummary()
}

func (s *eventService) GetAdminOrderStats(days int) ([]map[string]interface{}, error) {
	return s.repo.AdminGetOrderStats(days)
}

func (s *eventService) GetAdminOrdersSummaryStatus() (map[string]interface{}, error) {
	return s.repo.AdminGetOrdersSummaryStatus()
}

func (s *eventService) GetAdminTopSellingEvents(limit int) ([]map[string]interface{}, error) {
	return s.repo.AdminGetTopSellingEvents(limit)
}

func (s *eventService) GetOrganizerEventsCompare(orgID string) ([]map[string]interface{}, error) {
	return s.repo.OrganizerGetEventsCompare(orgID)
}
