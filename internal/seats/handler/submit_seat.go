package handler

import (
	"backend/internal/seats/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubmitSeatHandler struct {
	svc service.SeatService
}

type SubmitSeatRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
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

	var req SubmitSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "customer_id is required"})
		return
	}

	// Reserve seat via service directly
	err := h.svc.ReserveSeat(id, req.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to reserve seat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "seat reserved successfully",
	})
}
