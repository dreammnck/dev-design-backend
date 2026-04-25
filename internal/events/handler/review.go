package handler

import (
	"backend/internal/auth"
	"backend/internal/events"
	"backend/internal/events/repository"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /events/pending
func (h *EventHandler) GetPendingEvents(c *gin.Context) {
	evts, err := h.svc.GetPendingEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": evts})
}

// GET /organizer/event
func (h *EventHandler) OrganizerListEvents(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	evts, err := h.svc.GetMyEvents(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": mapToAdminEvents(evts)})
}

// GET /admin/event
func (h *EventHandler) AdminListEvents(c *gin.Context) {
	data, err := h.svc.GetAllAdminEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// POST /events/:id/review
func (h *EventHandler) ReviewEvent(c *gin.Context) {
	id := c.Param("id")

	var req events.ReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := h.svc.ReviewEvent(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	updatedEvt, err := h.svc.GetByID(id, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch updated event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Event review submitted: " + string(req.Status),
		"data":    updatedEvt,
	})
}

// PATCH /admin/editEvent/:id
func (h *EventHandler) AdminEditEvent(c *gin.Context) {
	id := c.Param("id")
	claims, _ := middleware.GetClaims(c)

	// If not admin, check ownership
	if claims.Role != auth.RoleAdmin {
		existingEvt, err := h.svc.GetByID(id, "")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Event not found"})
			return
		}

		// Ensure organization_id matches for organizers
		if claims.Role == auth.RoleOrganization {
			if existingEvt.OrganizationID == nil || *existingEvt.OrganizationID != claims.UserID {
				c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "You do not have permission to edit this event"})
				return
			}
		}
	}

	var req events.AdminEditEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := h.svc.AdminEditEvent(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	updatedEvt, err := h.svc.GetByID(id, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch updated event"})
		return
	}

	// Return the formatted updated event for Admin within an array to match the UI expectation
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   []AdminEventRes{mapToAdminEvent(updatedEvt)},
	})
}

// ── Mapping ───────────────────────────────────────────────────────────────

type AdminEventRes struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	ImageURL      string               `json:"imageUrl"`
	EventDate     string               `json:"eventDate"`
	EventTime     *string              `json:"eventTime"`
	Price         int                  `json:"price"`
	PublishStatus events.PublishStatus `json:"publishStatus"`
	RejectReason  *string              `json:"rejectReason"`
	PublishedAt   *string              `json:"publishedAt"`
	OrganizerName string               `json:"organizerName"`
	Location      AdminEventLocation   `json:"location"`
}

type AdminEventLocation struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	City    *string `json:"city"`
	Country *string `json:"country"`
}

func mapToAdminEvents(evts []events.Event) []AdminEventRes {
	res := make([]AdminEventRes, 0, len(evts))
	for _, e := range evts {
		res = append(res, mapToAdminEvent(e))
	}
	return res
}

func mapToAdminEvent(e events.Event) AdminEventRes {
	var publishedAt *string
	if e.PublishedAt != nil {
		p := e.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		publishedAt = &p
	}

	city := &e.Location.City
	if *city == "" {
		city = nil
	}
	country := &e.Location.Country
	if *country == "" {
		country = nil
	}

	return AdminEventRes{
		ID:            e.ID,
		Title:         e.Title,
		ImageURL:      e.Image,
		EventDate:     e.Date,
		EventTime:     e.Time,
		Price:         e.Price,
		PublishStatus: e.PublishStatus,
		RejectReason:  e.RejectReason,
		PublishedAt:   publishedAt,
		Location: AdminEventLocation{
			ID:      e.Location.ID,
			Name:    e.Location.Name,
			City:    city,
			Country: country,
		},
	}
}

// GET /organizer/events/booking/:eventId
func (h *EventHandler) GetOrganizerEventBookings(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	eventID := c.Param("eventId")
	rows, err := h.svc.GetBookingSummary(eventID, claims.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mapToOrganizerBookings(rows),
	})
}

type OrganizerBookingRes struct {
	BookingID     string   `json:"bookingId"`
	Buyer         string   `json:"buyer"`
	BuyerEmail    string   `json:"buyerEmail"`
	BuyerPhone    string   `json:"buyerPhone"`
	Notes         string   `json:"notes"`
	Seats         []string `json:"seats"`
	TicketType    *string  `json:"ticketType"`
	Qty           int      `json:"qty"`
	Total         int      `json:"total"`
	PaymentMethod string   `json:"paymentMethod"`
	StatusKey     string   `json:"statusKey"`
	StatusLabel   string   `json:"statusLabel"`
	CreatedAt     string   `json:"createdAt"`
}

// GET /organizer/events/summary/:eventId
func (h *EventHandler) GetOrganizerEventSummary(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	eventID := c.Param("eventId")
	rows, err := h.svc.GetBookingSummary(eventID, claims.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	totalBooking := 0
	totalTicket := len(rows)
	totalPrice := 0

	// Track unique bookings for count and price
	seenBookings := make(map[string]bool)
	for _, r := range rows {
		if !seenBookings[r.BookingID] {
			totalBooking++
			totalPrice += r.Amount
			seenBookings[r.BookingID] = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"totalBooking": totalBooking,
			"totalTicket":  totalTicket,
			"totalPrice":   totalPrice,
		},
	})
}

func mapToOrganizerBookings(rows []repository.BookingSummaryRow) []OrganizerBookingRes {
	// Consolidation map: bookingID -> object
	m := make(map[string]*OrganizerBookingRes)
	var order []string

	for _, r := range rows {
		if _, ok := m[r.BookingID]; !ok {
			label := "Confirmed"
			if r.Status == "paid" {
				label = "Paid"
			} else if r.Status == "pending" {
				label = "Pending"
			}

			tType := r.TicketType
			var tTypePtr *string
			if tType != "" {
				tTypePtr = &tType
			}

			m[r.BookingID] = &OrganizerBookingRes{
				BookingID:     r.BookingID,
				Buyer:         r.Buyer,
				BuyerEmail:    r.BuyerEmail,
				BuyerPhone:    r.BuyerPhone,
				Notes:         "-", // Default note if not in DB
				Seats:         []string{},
				TicketType:    tTypePtr,
				Qty:           0,
				Total:         r.Amount,
				PaymentMethod: "card", // Mock if not in DB
				StatusKey:     r.Status,
				StatusLabel:   label,
				CreatedAt:     r.BookedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			order = append(order, r.BookingID)
		}
		m[r.BookingID].Seats = append(m[r.BookingID].Seats, r.SeatNumber)
		m[r.BookingID].Qty++
	}

	res := make([]OrganizerBookingRes, 0, len(order))
	for _, id := range order {
		res = append(res, *m[id])
	}
	return res
}
