package handler

import (
	"backend/internal/auth"
	"backend/internal/auth/service"
	"backend/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.AuthService
}

func NewUserHandler(svc service.AuthService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	user, err := h.svc.GetUser(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req auth.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	res, err := h.svc.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": res})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req auth.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := h.svc.UpdateRole(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "role updated successfully"})
}

func (h *UserHandler) AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req auth.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := h.svc.AdminUpdateUser(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "user deleted successfully"})
}
