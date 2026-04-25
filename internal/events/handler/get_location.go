package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *EventHandler) GetLocations(c *gin.Context) {
	locations, err := h.svc.GetAllLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   locations,
	})
}
