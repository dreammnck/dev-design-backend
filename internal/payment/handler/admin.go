package handler

import (
	"backend/internal/payment/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GlobalAdminHandler struct {
	paymentSvc service.PaymentService
}

func NewGlobalAdminHandler(paymentSvc service.PaymentService) *GlobalAdminHandler {
	return &GlobalAdminHandler{paymentSvc: paymentSvc}
}

func (h *GlobalAdminHandler) GetAllBookings(c *gin.Context) {
	bookings, err := h.paymentSvc.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": bookings})
}

func (h *GlobalAdminHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.paymentSvc.GetAllPayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": payments})
}
func (h *GlobalAdminHandler) ScanQR(c *gin.Context) {
	var req struct {
		QRPayload string `json:"qrPayload" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := h.paymentSvc.RedeemTicket(req.QRPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
