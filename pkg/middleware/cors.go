package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS handles Cross-Origin Resource Sharing (CORS) headers enabling frontend access
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// allow all origins for development, you may want to restrict this in production
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
