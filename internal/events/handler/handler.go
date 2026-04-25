package handler

import (
	"backend/internal/events"
	"backend/internal/events/service"
)

type EventHandler struct {
	svc service.EventService
}

func NewEventHandler(svc service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

type FlattenedEvent struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Location     string `json:"location"`
	LocationType string `json:"locationType"`
	Time         string `json:"time"`
	Price    int    `json:"price"`
	Image    string `json:"image"`
	IsFav    bool   `json:"isFav"`
	Detail   string `json:"detail,omitempty"`
}

func flattenEvents(evts []events.Event) []FlattenedEvent {
	data := make([]FlattenedEvent, 0, len(evts))
	for _, e := range evts {
		data = append(data, flattenEvent(e))
	}
	return data
}

func flattenEvent(e events.Event) FlattenedEvent {
	timeStr := ""
	if e.Time != nil {
		timeStr = *e.Time
	}
	return FlattenedEvent{
		ID:       e.ID,
		Title:    e.Title,
		Date:     e.Date,
		Location:     e.Location.Name,
		LocationType: e.Location.Type,
		Time:         timeStr,
		Price:    e.Price,
		Image:    e.Image,
		IsFav:    e.IsFav,
		Detail:   e.Detail,
	}
}
