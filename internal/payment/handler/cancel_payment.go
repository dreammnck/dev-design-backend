package handler

import (
	"backend/internal/payment"
	"backend/internal/payment/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CancelPaymentHandler struct {
	svc service.PaymentService
}

func NewCancelPaymentHandler(svc service.PaymentService) *CancelPaymentHandler {
	return &CancelPaymentHandler{svc: svc}
}

func (h *CancelPaymentHandler) Handle(c *gin.Context) {
	var req payment.PaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid request body",
		})
		return
	}

	res, err := h.svc.Cancel(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
