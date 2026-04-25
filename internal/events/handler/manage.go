package handler

import (
	"backend/internal/events"
	"backend/pkg/middleware"
	"backend/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /events
func (h *EventHandler) CreateEvent(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	var req events.EventCreateReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var imageUrl string
	if req.ImageFile != nil {
		url, err := storage.UploadFile(req.ImageFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to upload image: " + err.Error()})
			return
		}
		imageUrl = url
	}

	// Map old fields if needed locally or directly pass them
	req.ImageFile = nil
	// (We will temporarily patch EventCreateReq DTO inside to also set Image if needed later)
	
	evt, err := h.svc.CreateEvent(claims.UserID, req, imageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   gin.H{"id": evt.ID},
	})
}

// PUT /events/:id
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	id := c.Param("id")
	var req events.EventUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	evt, err := h.svc.UpdateEvent(id, claims.UserID, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": evt})
}

// POST /events/:id/submit
func (h *EventHandler) SubmitForReview(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	id := c.Param("id")
	if err := h.svc.SubmitForReview(id, claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event submitted for review"})
}

// GET /events/my
func (h *EventHandler) GetMyEvents(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	evts, err := h.svc.GetMyEvents(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": evts})
}
