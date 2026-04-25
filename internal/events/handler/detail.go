package handler

import (
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *EventHandler) Detail(c *gin.Context) {
	id := c.Param("id")
	var userID string
	if claims, ok := middleware.GetClaims(c); ok {
		userID = claims.UserID
	}
	data, err := h.svc.GetByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": flattenEvent(data)})
}
