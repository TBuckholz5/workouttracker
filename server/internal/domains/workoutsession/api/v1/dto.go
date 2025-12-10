package v1

import "github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/models"

type CreateWorkoutSessionResponse struct {
	Session models.WorkoutSession `json:"session"`
}
