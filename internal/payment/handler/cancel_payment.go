package handler

import (
	"backend/internal/payment"
	"backend/internal/payment/service"
	"backend/pkg/middleware"
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

	res, err := h.svc.Cancel(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
