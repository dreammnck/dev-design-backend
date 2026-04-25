package handler

import (
	"backend/internal/events"
	evtSvc "backend/internal/events/service"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayoutHandler struct {
	evtSvc evtSvc.EventService
}

func NewPayoutHandler(evtSvc evtSvc.EventService) *PayoutHandler {
	return &PayoutHandler{evtSvc: evtSvc}
}

// POST /payouts
func (h *PayoutHandler) RequestPayout(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	var req events.PayoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	payout, err := h.evtSvc.RequestPayout(claims.UserID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": payout})
}

// GET /payouts
func (h *PayoutHandler) GetPayouts(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	payouts, err := h.evtSvc.GetPayouts(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": payouts})
}

// GET /admin/payouts
func (h *PayoutHandler) GetAllPayouts(c *gin.Context) {
	payouts, err := h.evtSvc.GetAllPayouts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": payouts})
}

// POST /admin/payouts/:id/process
func (h *PayoutHandler) ProcessPayout(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"` // approve, reject
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	payoutStatus := events.PayoutStatusCompleted
	if req.Status == "reject" {
		payoutStatus = events.PayoutStatusRejected
	} else if req.Status != "approve" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid status, expected approve or reject"})
		return
	}

	if err := h.evtSvc.ProcessPayout(id, payoutStatus, req.Note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Fetch updated list to match design: data: [ ... ]
	payouts, err := h.evtSvc.GetAllPayouts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   payouts,
	})
}
