package handler

import (
	"backend/internal/seats/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetSeatsHandler struct {
	svc service.SeatService
}

func NewGetSeatsHandler(svc service.SeatService) *GetSeatsHandler {
	return &GetSeatsHandler{svc: svc}
}

func (h *GetSeatsHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Query("event_id")
	}
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "event_id is required"})
		return
	}
	data, err := h.svc.GetByEventID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
