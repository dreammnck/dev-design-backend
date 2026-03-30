package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *EventHandler) ListAll(c *gin.Context) {
	data, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}
