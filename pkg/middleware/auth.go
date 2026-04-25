package middleware

import (
	"backend/internal/auth"
	jwtPkg "backend/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ClaimsKey = "claims"
)

// AuthRequired validates the Bearer JWT in the Authorization header.
// On success it stores *auth.JWTClaims in the gin context under ClaimsKey.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Authorization header format must be: Bearer <token>",
			})
			return
		}

		claims, err := jwtPkg.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
			})
			return
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

// OptionalAuth attempts to parse JWT but does not fail if it's missing or invalid
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			claims, err := jwtPkg.ValidateToken(parts[1])
			if err == nil {
				c.Set(ClaimsKey, claims)
			}
		}
		c.Next()
	}
}

// RolesAllowed returns a middleware that allows only the specified roles.
// Must be chained after AuthRequired().
func RolesAllowed(roles ...auth.UserRole) gin.HandlerFunc {
	allowed := make(map[auth.UserRole]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		claimsRaw, exists := c.Get(ClaimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
			return
		}

		claims, ok := claimsRaw.(*auth.JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to parse claims",
			})
			return
		}

		if _, ok := allowed[claims.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "You do not have permission to access this resource",
			})
			return
		}

		c.Next()
	}
}

// GetClaims is a helper to retrieve the JWT claims from gin context
func GetClaims(c *gin.Context) (*auth.JWTClaims, bool) {
	raw, exists := c.Get(ClaimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := raw.(*auth.JWTClaims)
	return claims, ok
}

// GetOptionalClaims returns claims if exists, otherwise returns empty one (not nil)
func GetOptionalClaims(c *gin.Context) *auth.JWTClaims {
	claims, ok := GetClaims(c)
	if !ok {
		return &auth.JWTClaims{}
	}
	return claims
}
