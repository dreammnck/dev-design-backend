package handler

import (
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *EventHandler) ListRecommend(c *gin.Context) {
	var userID string
	if claims, ok := middleware.GetClaims(c); ok {
		userID = claims.UserID
	}
	data, err := h.svc.GetRecommend(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": flattenEvents(data)})
}
