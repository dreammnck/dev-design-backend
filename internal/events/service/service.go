package service

import (
	"backend/internal/events"
	"backend/internal/events/repository"
)

type EventService interface {
	GetBanner() ([]events.Event, error)
	GetAll() ([]events.Event, error)
	GetRecommend() ([]events.Event, error)
	GetComingSoon() ([]events.Event, error)
	GetByID(id string) (events.Event, error)
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) GetBanner() ([]events.Event, error) {
	return s.repo.GetBanner()
}

func (s *eventService) GetAll() ([]events.Event, error) {
	return s.repo.GetAll()
}

func (s *eventService) GetRecommend() ([]events.Event, error) {
	return s.repo.GetRecommend()
}

func (s *eventService) GetComingSoon() ([]events.Event, error) {
	return s.repo.GetComingSoon()
}

func (s *eventService) GetByID(id string) (events.Event, error) {
	return s.repo.GetByID(id)
}
