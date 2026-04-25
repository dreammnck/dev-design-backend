package handler

import (
	"backend/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/organizer/summary
func (h *EventHandler) GetOrganizerOverallSummary(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	data, err := h.svc.GetOrganizerOverallSummary(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// GET /api/organizer/events/:id/sales?period=30d
func (h *EventHandler) GetOrganizerEventSales(c *gin.Context) {
	// claims check ignored for simplicity in ID-based lookup if needed,
	// but normally we should check if org owns this event.
	eventID := c.Param("id")
	periodStr := c.DefaultQuery("period", "30d")

	// Simple parsing for "30d" -> 30
	days := 30
	if len(periodStr) > 1 {
		if d, err := strconv.Atoi(periodStr[:len(periodStr)-1]); err == nil {
			days = d
		}
	}

	data, err := h.svc.GetOrganizerEventSales(eventID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// GET /api/organizer/events/:id/seats/summary
func (h *EventHandler) GetOrganizerEventSeatsDashboard(c *gin.Context) {
	eventID := c.Param("id")

	data, err := h.svc.GetOrganizerEventSeatsSummary(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// GET /api/organizer/events/compare
func (h *EventHandler) GetOrganizerEventsCompare(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	data, err := h.svc.GetOrganizerEventsCompare(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// GetAdminOverallSummary returns global system metrics
func (h *EventHandler) GetAdminOverallSummary(c *gin.Context) {
	summary, err := h.svc.GetAdminOverallSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary, "status": "success"})
}

// GetAdminOrderStats returns daily sales/orders stats for the whole system
func (h *EventHandler) GetAdminOrderStats(c *gin.Context) {
	periodStr := c.DefaultQuery("period", "30d")
	days := 30
	if len(periodStr) > 1 {
		if d, err := strconv.Atoi(periodStr[:len(periodStr)-1]); err == nil {
			days = d
		}
	}

	stats, err := h.svc.GetAdminOrderStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats, "status": "success"})
}

// GetAdminOrdersSummaryStatus returns order status breakdown (Doughnut Chart)
func (h *EventHandler) GetAdminOrdersSummaryStatus(c *gin.Context) {
	summary, err := h.svc.GetAdminOrdersSummaryStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary, "status": "success"})
}

// GetAdminTopSellingEvents returns top 5 selling events (Bar Chart)
func (h *EventHandler) GetAdminTopSellingEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	result, err := h.svc.GetAdminTopSellingEvents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "status": "success"})
}
