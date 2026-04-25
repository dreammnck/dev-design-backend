package handler

import (
	"backend/internal/seats/service"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubmitSeatHandler struct {
	svc service.SeatService
}

func NewSubmitSeatHandler(svc service.SeatService) *SubmitSeatHandler {
	return &SubmitSeatHandler{svc: svc}
}

func (h *SubmitSeatHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "seat id is required"})
		return
	}

	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	// Reserve seat via service directly using the authenticated user's ID
	if err := h.svc.ReserveSeat(id, claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to reserve seat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "seat reserved successfully",
	})
}
