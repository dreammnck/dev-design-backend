package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *EventHandler) Detail(c *gin.Context) {
	id := c.Param("id")
	data, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}
