package repository

import (
	"backend/internal/events"

	"gorm.io/gorm"
)

type EventRepository interface {
	GetBanner() ([]events.Event, error)
	GetAll() ([]events.Event, error)
	GetRecommend() ([]events.Event, error)
	GetComingSoon() ([]events.Event, error)
	GetByID(id string) (events.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetBanner() ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").Where("is_banner = ?", true).Find(&evts).Error
	return evts, err
}

func (r *eventRepository) GetAll() ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").Find(&evts).Error
	return evts, err
}

func (r *eventRepository) GetRecommend() ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").Where("is_recommend = ?", true).Find(&evts).Error
	return evts, err
}

func (r *eventRepository) GetComingSoon() ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").Where("is_coming_soon = ?", true).Find(&evts).Error
	return evts, err
}

func (r *eventRepository) GetByID(id string) (events.Event, error) {
	var evt events.Event
	err := r.db.Preload("Location").Where("id = ?", id).First(&evt).Error
	return evt, err
}
