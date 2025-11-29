package v1

import (
	"github.com/TBuckholz5/workouttracker/internal/domains/user/api/v1/dto"
	"github.com/TBuckholz5/workouttracker/internal/domains/user/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.UserService
}

func NewHandler(s service.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	var payload dto.RegisterRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateUser(c.Request.Context(), &payload); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) Login(c *gin.Context) {
	var payload dto.LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.AuthenticateUser(c.Request.Context(), &payload)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
