package user

import (
	"github.com/TBuckholz5/workouttracker/internal/api/v1/exercise/dto"
	service "github.com/TBuckholz5/workouttracker/internal/service/exercise"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ExerciseService
}

func NewHandler(s service.ExerciseService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateExercise(c *gin.Context) {
	var payload dto.CreateExerciseRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateExercise(c.Request.Context(), &payload); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) GetExerciseForUser(c *gin.Context) {
	var payload dto.GetExerciseForUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	exercises, err := h.service.GetExercisesForUser(c.Request.Context(), &payload)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"exercises": exercises})
}
