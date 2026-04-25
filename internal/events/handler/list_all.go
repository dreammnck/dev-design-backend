package handler

import (
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *EventHandler) ListAll(c *gin.Context) {
	search := c.Query("search")
	claims := middleware.GetOptionalClaims(c)

	evts, err := h.svc.GetAll(claims.UserID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": flattenEvents(evts)})
}
