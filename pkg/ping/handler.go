package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is a simple health-check handler
func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
