package handler

import (
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToggleFavReq struct {
	EventID string `json:"eventId" binding:"required"`
}

func (h *EventHandler) ToggleFavorite(c *gin.Context) {
	var req ToggleFavReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "eventId is required in body"})
		return
	}

	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	isFavorite, err := h.svc.ToggleFavorite(claims.UserID, req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	message := "Added to favorites"
	if !isFavorite {
		message = "Removed from favorites"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data": gin.H{
			"isFav": isFavorite,
		},
	})
}

func (h *EventHandler) ListFavorist(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	data, err := h.svc.GetFavoritedEvents(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   flattenEvents(data),
	})
}
