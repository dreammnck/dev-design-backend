package handler

import (
	"backend/internal/payment/service"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MyBookingHandler struct {
	svc service.PaymentService
}

func NewMyBookingHandler(svc service.PaymentService) *MyBookingHandler {
	return &MyBookingHandler{svc: svc}
}

func (h *MyBookingHandler) GetMyBookings(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	bookings, err := h.svc.GetMyBookings(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   bookings,
	})
}
func (h *MyBookingHandler) GetTicketByID(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID is required"})
		return
	}

	bookings, err := h.svc.GetTicketByID(id, claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if len(bookings) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Ticket not found or ownership mismatch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   bookings,
	})
}
