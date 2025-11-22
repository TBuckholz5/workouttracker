package exercise

import (
	"strconv"

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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(500, gin.H{"error": "userID not found in context"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	payload := dto.GetExerciseForUserRequest{
		UserID: userID.(int64),
		Offset: offset,
		Limit:  limit,
	}
	exercises, err := h.service.GetExercisesForUser(c.Request.Context(), &payload)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"exercises": exercises})
}
