package v1

import (
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/models"
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.WorkoutSessionService
}

func NewHandler(s service.WorkoutSessionService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c *gin.Context) {
	var payload models.WorkoutSession
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(500, gin.H{"error": "userID not found in context"})
		return
	}
	payload.UserID = userID.(int64)
	session, err := h.service.Create(c.Request.Context(), &payload)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"session": session})
}
