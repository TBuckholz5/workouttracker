package service

import (
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/models"
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/repository"
)

func repositoryToModels(session *repository.WorkoutSession, workouts []*repository.Workout, sets []*repository.WorkoutSet) *models.WorkoutSession {
	workoutMap := make(map[int64]*models.Workout)
	for _, workout := range workouts {
		workoutMap[workout.ID] = &models.Workout{
			ID:          workout.ID,
			ExerciseID:  workout.ExerciseID,
			Description: workout.Description,
			Sets:        []models.WorkoutSet{},
		}
	}
	for _, set := range sets {
		workoutMap[set.WorkoutID].Sets = append(workoutMap[set.WorkoutID].Sets, models.WorkoutSet{
			ID:       set.ID,
			Reps:     set.Reps,
			Weight:   set.Weight,
			SetType:  set.SetType,
			SetOrder: set.SetOrder,
		})
	}
	var modelWorkouts []models.Workout
	for _, workout := range workoutMap {
		modelWorkouts = append(modelWorkouts, *workout)
	}
	modelSession := &models.WorkoutSession{
		ID:          session.ID,
		UserID:      session.UserID,
		Name:        session.Name,
		Description: session.Description,
		Duration:    session.Duration,
		Workouts:    modelWorkouts,
		CreatedAt:   session.CreatedAt,
	}
	return modelSession
}
