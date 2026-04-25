package handler

import (
	"backend/internal/payment"
	"backend/internal/payment/service"
	"backend/pkg/middleware"
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
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	var req payment.PaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid request body",
		})
		return
	}

	req.UserID = claims.UserID

	res, err := h.svc.Process(req)
	if err != nil {
		c.JSON(http.StatusPaymentRequired, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
