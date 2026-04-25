package handler

import (
	evtSvc "backend/internal/events/service"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingAdminHandler struct {
	evtSvc evtSvc.EventService
}

func NewBookingAdminHandler(evtSvc evtSvc.EventService) *BookingAdminHandler {
	return &BookingAdminHandler{evtSvc: evtSvc}
}

// GET /events/:id/bookings
func (h *BookingAdminHandler) GetBookingSummary(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	id := c.Param("id")
	rows, err := h.evtSvc.GetBookingSummary(id, claims.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": rows})
}
