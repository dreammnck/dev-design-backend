package handler

import (
	"backend/internal/events"
	evtSvc "backend/internal/events/service"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeatAdminHandler struct {
	evtSvc evtSvc.EventService
}

func NewSeatAdminHandler(evtSvc evtSvc.EventService) *SeatAdminHandler {
	return &SeatAdminHandler{evtSvc: evtSvc}
}

// POST /events/:id/seats
func (h *SeatAdminHandler) AddSeats(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	id := c.Param("id")
	var req events.SeatBatchCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := h.evtSvc.AddSeats(id, claims.UserID, req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "seats added successfully"})
}
