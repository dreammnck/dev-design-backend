package handler

import (
	"backend/internal/auth"
	"backend/internal/auth/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	svc service.AuthService
}

func NewLoginHandler(svc service.AuthService) *LoginHandler {
	return &LoginHandler{svc: svc}
}

// Handle POST /auth/login
func (h *LoginHandler) Handle(c *gin.Context) {
	var req auth.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body: email and password are required",
		})
		return
	}

	res, err := h.svc.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid email or password",
			})
			return
		}
		if errors.Is(err, service.ErrUserInactive) {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Your account has been deactivated",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   res,
	})
}
