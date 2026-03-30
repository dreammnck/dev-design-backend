package handler

import (
	"backend/internal/payment"
	"backend/internal/payment/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfirmPaymentHandler struct {
	svc service.PaymentService
}

func NewConfirmPaymentHandler(svc service.PaymentService) *ConfirmPaymentHandler {
	return &ConfirmPaymentHandler{svc: svc}
}

func (h *ConfirmPaymentHandler) Handle(c *gin.Context) {
	var req payment.PaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid request body",
		})
		return
	}

	res, err := h.svc.Confirm(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
