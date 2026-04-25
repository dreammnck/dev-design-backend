package repository

import (
	"backend/internal/events"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrEventNotFound = errors.New("event not found")
var ErrNotOwner = errors.New("you do not own this event")

type EventRepository interface {
	// Public (approved only)
	GetBanner(userID string) ([]events.Event, error)
	GetAll(userID string, search string) ([]events.Event, error)
	GetRecommend(userID string) ([]events.Event, error)
	GetComingSoon(userID string) ([]events.Event, error)
	GetByID(id string, userID string) (events.Event, error)

	// Favorites
	ToggleFavorite(userID string, eventID string) (bool, error)
	GetFavoritedEvents(userID string) ([]events.Event, error)

	// Organization
	Create(orgID string, req events.EventCreateReq, imageUrl string) (events.Event, error)
	Update(id string, orgID string, req events.EventUpdateReq) (events.Event, error)
	SubmitForReview(id string, orgID string) error
	GetByOrganization(orgID string) ([]events.Event, error)
	GetBookingSummary(eventID string, orgID string) ([]BookingSummaryRow, error)

	// Seat management (org)
	AddSeats(eventID string, orgID string, seats []events.SeatInput) error

	// Admin review
	Review(id string, status events.PublishStatus, rejectReason string) error
	GetPendingEvents() ([]events.Event, error)
	GetAllAdminEvents() ([]map[string]interface{}, error)
	AdminEditEvent(id string, req events.AdminEditEventReq) error

	// Payout
	CreatePayout(orgID string, req events.PayoutReq) (events.Payout, error)
	GetPayouts(orgID string) ([]events.Payout, error)
	GetAllPayouts() ([]map[string]interface{}, error)
	ProcessPayout(payoutID string, status events.PayoutStatus, note string) error
	GetTotalIncomeByOrgID(orgID string) (int, error)

	// Dashboard
	OrganizerGetOverallSummary(orgID string) (map[string]interface{}, error)
	OrganizerGetEventSales(eventID string, days int) ([]map[string]interface{}, error)
	OrganizerGetEventSeatsSummary(eventID string) (map[string]interface{}, error)
	OrganizerGetEventsCompare(orgID string) ([]map[string]interface{}, error)

	// Admin Dashboard
	AdminGetOverallSummary() (map[string]interface{}, error)
	AdminGetOrderStats(days int) ([]map[string]interface{}, error)
	AdminGetOrdersSummaryStatus() (map[string]interface{}, error)
	AdminGetTopSellingEvents(limit int) ([]map[string]interface{}, error)

	// Locations
	GetAllLocations() ([]events.Location, error)
}

type BookingSummaryRow struct {
	BookingID  string    `json:"bookingId"   gorm:"column:booking_id"`
	Buyer      string    `json:"buyer"       gorm:"column:buyer"`
	BuyerEmail string    `json:"buyerEmail"  gorm:"column:buyer_email"`
	BuyerPhone string    `json:"buyerPhone"  gorm:"column:buyer_phone"`
	SeatNumber string    `json:"seatNumber"  gorm:"column:seat_number"`
	TicketType string    `json:"ticketType"  gorm:"column:ticket_type"`
	Amount     int       `json:"amount"      gorm:"column:amount"`
	Status     string    `json:"status"      gorm:"column:status"`
	BookedAt   time.Time `json:"bookedAt"    gorm:"column:booked_at"`
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

// ── Public (approved only) ────────────────────────────────────────────────────

func (r *eventRepository) populateFavorites(evts []events.Event, userID string) {
	if userID == "" || len(evts) == 0 {
		return
	}
	var eventIDs []string
	for _, e := range evts {
		eventIDs = append(eventIDs, e.ID)
	}
	var favIDs []string
	r.db.Table("user_favorites").Where("user_id = ? AND event_id IN ?", userID, eventIDs).Pluck("event_id", &favIDs)
	favMap := make(map[string]bool)
	for _, id := range favIDs {
		favMap[id] = true
	}
	for i := range evts {
		evts[i].IsFav = favMap[evts[i].ID]
	}
}

func (r *eventRepository) GetBanner(userID string) ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").
		Where("is_banner = ? AND publish_status = ?", true, events.PublishStatusApproved).
		Find(&evts).Error
	if err == nil {
		r.populateFavorites(evts, userID)
	}
	return evts, err
}

func (r *eventRepository) GetAll(userID string, search string) ([]events.Event, error) {
	var evts []events.Event
	db := r.db.Preload("Location").Where("publish_status = ?", events.PublishStatusApproved)

	if search != "" {
		searchTerm := "%" + search + "%"
		db = db.Joins("LEFT JOIN locations ON locations.id = events.location_id").
			Where("(events.title ILIKE ? OR events.description ILIKE ? OR locations.name ILIKE ? OR locations.city ILIKE ?)",
				searchTerm, searchTerm, searchTerm, searchTerm)
	}

	err := db.Find(&evts).Error
	if err == nil {
		r.populateFavorites(evts, userID)
	}
	return evts, err
}

func (r *eventRepository) GetRecommend(userID string) ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").
		Where("is_recommend = ? AND publish_status = ?", true, events.PublishStatusApproved).
		Find(&evts).Error
	if err == nil {
		r.populateFavorites(evts, userID)
	}
	return evts, err
}

func (r *eventRepository) GetComingSoon(userID string) ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").
		Where("is_coming_soon = ? AND publish_status = ?", true, events.PublishStatusApproved).
		Find(&evts).Error
	if err == nil {
		r.populateFavorites(evts, userID)
	}
	return evts, err
}

func (r *eventRepository) GetByID(id string, userID string) (events.Event, error) {
	var evt events.Event
	err := r.db.Preload("Location").Where("id = ?", id).First(&evt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return evt, ErrEventNotFound
	}
	if err == nil {
		var count int64
		r.db.Table("user_favorites").Where("user_id = ? AND event_id = ?", userID, id).Count(&count)
		evt.IsFav = count > 0
	}
	return evt, err
}

func (r *eventRepository) ToggleFavorite(userID string, eventID string) (bool, error) {
	var fav events.UserFavorite
	err := r.db.Where("user_id = ? AND event_id = ?", userID, eventID).First(&fav).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Like: Create record
		newFav := events.UserFavorite{
			UserID:    userID,
			EventID:   eventID,
			CreatedAt: time.Now(),
		}
		if err := r.db.Create(&newFav).Error; err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}

	// Unlike: Remove record
	if err := r.db.Delete(&fav).Error; err != nil {
		return false, err
	}
	return false, nil
}

func (r *eventRepository) GetFavoritedEvents(userID string) ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").
		Joins("JOIN user_favorites uf ON uf.event_id = events.id").
		Where("uf.user_id = ?", userID).
		Find(&evts).Error

	if err == nil {
		for i := range evts {
			evts[i].IsFav = true
		}
	}
	return evts, err
}

// ── Organization ──────────────────────────────────────────────────────────────

func (r *eventRepository) Create(orgID string, req events.EventCreateReq, imageUrl string) (events.Event, error) {
	loc := events.Location{
		Name:          req.LocationName,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		City:          req.City,
		StateProvince: req.StateProvince,
		Country:       req.Country,
		PostCode:      req.PostCode,
		Type:          req.LocationType,
		IsActive:      true,
	}
	if err := r.db.Create(&loc).Error; err != nil {
		return events.Event{}, err
	}

	var eventTime *string
	if req.EventTime != "" {
		eventTime = &req.EventTime
	}

	evt := events.Event{
		Title:          req.Title,
		Detail:         req.Description,
		Image:          imageUrl,
		LocationID:     loc.ID,
		Date:           req.EventDate,
		Time:           eventTime,
		Price:          req.Price,
		IsBanner:       req.IsBanner,
		IsRecommend:    req.IsRecommend,
		IsComingSoon:   req.IsComingSoon,
		OrganizationID: &orgID,
		PublishStatus:  events.PublishStatusPending,
	}
	if err := r.db.Create(&evt).Error; err != nil {
		return events.Event{}, err
	}
	return evt, nil
}

func (r *eventRepository) Update(id string, orgID string, req events.EventUpdateReq) (events.Event, error) {
	var evt events.Event
	if err := r.db.Where("id = ? AND organization_id = ?", id, orgID).First(&evt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return evt, ErrNotOwner
		}
		return evt, err
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Detail != nil {
		updates["description"] = *req.Detail
	}
	if req.Image != nil {
		updates["image_url"] = *req.Image
	}
	if req.LocationID != nil {
		updates["location_id"] = *req.LocationID
	}
	if req.Date != nil {
		updates["event_date"] = *req.Date
	}
	if req.Time != nil {
		updates["event_time"] = *req.Time
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.IsBanner != nil {
		updates["is_banner"] = *req.IsBanner
	}
	if req.IsRecommend != nil {
		updates["is_recommend"] = *req.IsRecommend
	}
	if req.IsComingSoon != nil {
		updates["is_coming_soon"] = *req.IsComingSoon
	}

	if err := r.db.Model(&evt).Updates(updates).Error; err != nil {
		return evt, err
	}
	return evt, nil
}

func (r *eventRepository) SubmitForReview(id string, orgID string) error {
	result := r.db.Model(&events.Event{}).
		Where("id = ? AND organization_id = ?", id, orgID).
		Update("publish_status", events.PublishStatusPending)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotOwner
	}
	return nil
}

func (r *eventRepository) GetByOrganization(orgID string) ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").Where("organization_id = ?", orgID).Find(&evts).Error
	return evts, err
}

func (r *eventRepository) GetBookingSummary(eventID string, orgID string) ([]BookingSummaryRow, error) {
	// Verify event belongs to org first
	var count int64
	r.db.Model(&events.Event{}).Where("id = ? AND organization_id = ?", eventID, orgID).Count(&count)
	if count == 0 {
		return nil, ErrNotOwner
	}

	var rows []BookingSummaryRow
	err := r.db.Raw(`
		SELECT b.id AS booking_id, u.username AS buyer, u.email AS buyer_email,
		       s.seat_number, s.seat_type AS ticket_type, b.total_amount AS amount,
		       b.status, b.created_at AS booked_at
		FROM bookings b
		JOIN users u ON u.id = b.user_id
		JOIN booking_seats bs ON bs.booking_id = b.id
		JOIN seats s ON s.id = bs.seat_id
		WHERE s.event_id = ?
		ORDER BY b.created_at DESC
	`, eventID).Scan(&rows).Error
	return rows, err
}

func (r *eventRepository) AddSeats(eventID string, orgID string, seatInputs []events.SeatInput) error {
	// Verify ownership
	var count int64
	r.db.Model(&events.Event{}).Where("id = ? AND organization_id = ?", eventID, orgID).Count(&count)
	if count == 0 {
		return ErrNotOwner
	}

	type seatRow struct {
		EventID    string `gorm:"column:event_id"`
		SeatNumber string `gorm:"column:seat_number"`
		Price      int    `gorm:"column:price"`
		SeatType   string `gorm:"column:seat_type"`
		Status     string `gorm:"column:status"`
	}

	rows := make([]seatRow, 0)
	for _, s := range seatInputs {
		if s.Capacity > 0 {
			// Auto-generate seats for capacity
			prefix := s.SeatType
			if prefix == "" {
				prefix = "Zone"
			}
			for i := 1; i <= s.Capacity; i++ {
				rows = append(rows, seatRow{
					EventID:    eventID,
					SeatNumber: fmt.Sprintf("%s-%d", prefix, i),
					Price:      s.Price,
					SeatType:   s.SeatType,
					Status:     "available",
				})
			}
		} else {
			// Normal individual seat
			rows = append(rows, seatRow{
				EventID:    eventID,
				SeatNumber: s.SeatNumber,
				Price:      s.Price,
				SeatType:   s.SeatType,
				Status:     "available",
			})
		}
	}
	return r.db.Table("seats").Create(&rows).Error
}

// ── Admin ─────────────────────────────────────────────────────────────────────

func (r *eventRepository) Review(id string, status events.PublishStatus, rejectReason string) error {
	updates := map[string]interface{}{
		"publish_status": status,
		"reject_reason":  nil,
		"published_at":   nil,
	}
	if status == events.PublishStatusApproved {
		now := time.Now()
		updates["published_at"] = now
	}
	if status == events.PublishStatusRejected {
		updates["reject_reason"] = rejectReason
	}
	result := r.db.Model(&events.Event{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrEventNotFound
	}
	return nil
}

func (r *eventRepository) GetPendingEvents() ([]events.Event, error) {
	var evts []events.Event
	err := r.db.Preload("Location").
		Where("publish_status = ?", events.PublishStatusPending).
		Order("created_at ASC").
		Find(&evts).Error
	return evts, err
}

func (r *eventRepository) GetAllAdminEvents() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT 
			e.*,
			e.publish_status as "publishStatus",
			e.reject_reason as "rejectReason",
			e.published_at as "publishedAt",
			e.image_url as "imageUrl",
			e.event_date as "eventDate",
			e.event_time as "eventTime",
			l.name as "locationName",
			l.type as "locationType",
			COALESCE(u_usr.username, u_org.username, 'Unknown') as "organizerName"
		FROM events e
		LEFT JOIN locations l ON e.location_id = l.id
		LEFT JOIN users u_usr ON e.user_id = u_usr.id
		LEFT JOIN users u_org ON e.organization_id = u_org.id
		ORDER BY e.created_at DESC
	`
	err := r.db.Raw(query).Scan(&results).Error
	return results, err
}

func (r *eventRepository) AdminEditEvent(id string, req events.AdminEditEventReq) error {
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.LocationID != nil {
		updates["location_id"] = *req.LocationID
	}
	if req.EventDate != nil {
		updates["event_date"] = *req.EventDate
	}
	if req.EventTime != nil {
		updates["event_time"] = *req.EventTime
	}
	if req.IsBanner != nil {
		updates["is_banner"] = *req.IsBanner
	}
	if req.IsRecommend != nil {
		updates["is_recommend"] = *req.IsRecommend
	}
	if req.IsComingSoon != nil {
		updates["is_coming_soon"] = *req.IsComingSoon
	}
	if req.PublishStatus != nil {
		updates["publish_status"] = string(*req.PublishStatus)
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.Model(&events.Event{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrEventNotFound
	}
	return nil
}

// ── Payout ────────────────────────────────────────────────────────────────────

func (r *eventRepository) CreatePayout(orgID string, req events.PayoutReq) (events.Payout, error) {
	var eventIDPtr *string
	if req.EventID != "" {
		eventIDPtr = &req.EventID
	}

	payout := events.Payout{
		OrganizationID: orgID,
		EventID:        eventIDPtr,
		Amount:         req.Amount,
		AccountName:    req.AccountName,
		BankAccount:    req.BankAccount,
		BankName:       req.BankName,
		Status:         events.PayoutStatusRequested,
		RequestedAt:    time.Now(),
	}
	if err := r.db.Create(&payout).Error; err != nil {
		return events.Payout{}, err
	}
	return payout, nil
}

func (r *eventRepository) GetPayouts(orgID string) ([]events.Payout, error) {
	var payouts []events.Payout
	err := r.db.Where("organization_id = ?", orgID).Order("requested_at DESC").Find(&payouts).Error
	return payouts, err
}

func (r *eventRepository) GetAllPayouts() ([]map[string]interface{}, error) {
	var payouts []map[string]interface{}
	query := `
		SELECT 
			p.id, 
			p.organization_id as "organizationId",
			u.username as "organizationName",
			p.amount,
			p.bank_name as "bankName",
			p.bank_account as "bankAccount",
			p.status,
			p.reject_reason as "rejectReason",
			p.requested_at as "requestedAt",
			p.processed_at as "processedAt"
		FROM payouts p
		LEFT JOIN users u ON p.organization_id = u.id
		ORDER BY p.requested_at DESC
	`
	if err := r.db.Raw(query).Scan(&payouts).Error; err != nil {
		return nil, err
	}

	// Map DB status 'requested' to 'pending' for UI consistency as per design
	for _, p := range payouts {
		if p["status"] == "requested" {
			p["status"] = "pending"
		}
	}

	return payouts, nil
}

func (r *eventRepository) ProcessPayout(payoutID string, status events.PayoutStatus, note string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"processed_at": now,
	}
	if note != "" {
		updates["reject_reason"] = note
	}

	result := r.db.Model(&events.Payout{}).Where("id = ?", payoutID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("payout not found")
	}
	return nil
}
func (r *eventRepository) GetAllLocations() ([]events.Location, error) {
	var locations []events.Location
	if err := r.db.Where("is_active = ?", true).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *eventRepository) GetTotalIncomeByOrgID(orgID string) (int, error) {
	var totalIncome int
	err := r.db.Table("bookings").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("events.organization_id = ? AND bookings.status = 'confirmed'", orgID).
		Select("COALESCE(SUM(bookings.total_amount), 0)").
		Scan(&totalIncome).Error
	return totalIncome, err
}

func (r *eventRepository) OrganizerGetOverallSummary(orgID string) (map[string]interface{}, error) {
	var summary struct {
		TotalRevenue int
		TicketsSold  int64
		MyEvents     int64
	}

	// 1. Count Total Events owned by this organizer (Should NOT be 0 if events exist)
	r.db.Model(&events.Event{}).Where("organization_id = ?", orgID).Count(&summary.MyEvents)

	// 2. Total Revenue (Summarize from confirmed bookings only)
	r.db.Table("bookings").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("events.organization_id = ? AND bookings.status = 'confirmed'", orgID).
		Select("COALESCE(SUM(bookings.total_amount), 0)").
		Scan(&summary.TotalRevenue)

	// 3. Tickets Sold (Count seats in confirmed bookings)
	r.db.Table("booking_seats").
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("events.organization_id = ? AND bookings.status = 'confirmed'", orgID).
		Count(&summary.TicketsSold)

	return map[string]interface{}{
		"totalRevenue": summary.TotalRevenue,
		"ticketsSold":  summary.TicketsSold,
		"paidSeats":    summary.TicketsSold,
		"myEvents":     summary.MyEvents,
	}, nil
}

func (r *eventRepository) OrganizerGetEventSales(eventID string, days int) ([]map[string]interface{}, error) {
	stats := []map[string]interface{}{} // Initialize as empty slice to avoid 'null' in JSON

	// Query daily revenue and orders for the last N days
	err := r.db.Raw(`
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as date,
			COALESCE(SUM(total_amount), 0) as revenue,
			COUNT(id) as orders
		FROM bookings
		WHERE event_id = ? AND status = 'confirmed'
		  AND created_at >= CURRENT_DATE - (INTERVAL '1 day' * ?)
		GROUP BY date, DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at) ASC
	`, eventID, days).Scan(&stats).Error
	return stats, err
}

func (r *eventRepository) OrganizerGetEventSeatsSummary(eventID string) (map[string]interface{}, error) {
	var counts []struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}

	r.db.Raw(`
		SELECT status, COUNT(*) as count
		FROM seats
		WHERE event_id = ?
		GROUP BY status
	`, eventID).Scan(&counts)

	res := map[string]interface{}{
		"paid":      int64(0),
		"reserved":  int64(0),
		"available": int64(0),
	}

	for _, c := range counts {
		switch c.Status {
		case "sold":
			res["paid"] = c.Count
		case "reserved":
			res["reserved"] = c.Count
		case "available":
			res["available"] = c.Count
		}
	}

	return res, nil
}

func (r *eventRepository) OrganizerGetEventsCompare(orgID string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := r.db.Raw(`
		SELECT 
			e.title as name,
			COALESCE(SUM(b.total_amount), 0) as revenue
		FROM events e
		LEFT JOIN bookings b ON b.event_id = e.id AND b.status = 'confirmed'
		WHERE e.organization_id = ?
		GROUP BY e.id, e.title
		ORDER BY revenue DESC
	`, orgID).Scan(&result).Error
	return result, err
}

func (r *eventRepository) AdminGetOverallSummary() (map[string]interface{}, error) {
	var summary struct {
		TotalRevenue int   `gorm:"column:total_revenue"`
		TotalOrders  int64 `gorm:"column:total_orders"`
		PaidOrders   int64 `gorm:"column:paid_orders"`
	}

	// 1. Total Revenue & Order counts
	r.db.Raw(`
		SELECT 
			COALESCE(SUM(CASE WHEN status = 'confirmed' THEN total_amount ELSE 0 END), 0) as total_revenue,
			COUNT(*) as total_orders,
			COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as paid_orders
		FROM bookings
	`).Scan(&summary)

	// 2. Top Event Revenue
	var topEvent struct {
		Title   string `gorm:"column:title"`
		Revenue int    `gorm:"column:revenue"`
	}
	r.db.Raw(`
		SELECT e.title, SUM(b.total_amount) as revenue
		FROM events e
		JOIN bookings b ON b.event_id = e.id
		WHERE b.status = 'confirmed'
		GROUP BY e.id, e.title
		ORDER BY revenue DESC
		LIMIT 1
	`).Scan(&topEvent)

	return map[string]interface{}{
		"totalRevenue":    summary.TotalRevenue,
		"totalOrders":     summary.TotalOrders,
		"paidOrders":      summary.PaidOrders,
		"topEventRevenue": topEvent.Revenue,
		"topEventName":    topEvent.Title,
	}, nil
}

func (r *eventRepository) AdminGetOrderStats(days int) ([]map[string]interface{}, error) {
	var stats []map[string]interface{}
	err := r.db.Raw(`
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as date,
			COALESCE(SUM(total_amount), 0) as revenue,
			COUNT(id) as orders
		FROM bookings
		WHERE status = 'confirmed'
		  AND created_at >= CURRENT_DATE - (INTERVAL '1 day' * ?)
		GROUP BY date, DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at) ASC
	`, days).Scan(&stats).Error
	return stats, err
}

func (r *eventRepository) AdminGetOrdersSummaryStatus() (map[string]interface{}, error) {
	var counts []struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}

	r.db.Raw(`
		SELECT status, COUNT(*) as count
		FROM bookings
		GROUP BY status
	`).Scan(&counts)

	res := map[string]interface{}{
		"paid":    int64(0),
		"pending": int64(0),
		"failed":  int64(0),
	}

	for _, c := range counts {
		switch c.Status {
		case "confirmed":
			res["paid"] = c.Count
		case "pending":
			res["pending"] = c.Count
		case "cancelled":
			res["failed"] = c.Count
		}
	}

	return res, nil
}

func (r *eventRepository) AdminGetTopSellingEvents(limit int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := r.db.Raw(`
		SELECT 
			e.title as name,
			COALESCE(SUM(b.total_amount), 0) as revenue
		FROM events e
		LEFT JOIN bookings b ON b.event_id = e.id AND b.status = 'confirmed'
		GROUP BY e.id, e.title
		ORDER BY revenue DESC
		LIMIT ?
	`, limit).Scan(&result).Error
	return result, err
}
