package user

import (
	"github.com/TBuckholz5/workouttracker/internal/api/v1/user/dto"
	"github.com/TBuckholz5/workouttracker/internal/service/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *user.Service
}

func NewHandler(s *user.Service) *Handler {
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
