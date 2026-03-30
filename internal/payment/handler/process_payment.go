package handler

import (
	"backend/internal/payment"
	"backend/internal/payment/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProcessPaymentHandler struct {
	svc service.PaymentService
}

func NewProcessPaymentHandler(svc service.PaymentService) *ProcessPaymentHandler {
	return &ProcessPaymentHandler{svc: svc}
}

func (h *ProcessPaymentHandler) Handle(c *gin.Context) {
	var req payment.PaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid request body",
		})
		return
	}

	res, err := h.svc.Process(req)
	if err != nil {
		c.JSON(http.StatusPaymentRequired, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
