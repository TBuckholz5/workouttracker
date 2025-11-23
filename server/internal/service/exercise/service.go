package exercise

import (
	"context"

	repo "github.com/TBuckholz5/workouttracker/internal/repository/exercise"
	"github.com/TBuckholz5/workouttracker/internal/service/exercise/models"
)

type ExerciseService interface {
	CreateExercise(reqContext context.Context, params *models.CreateExerciseForUserParams) error
	GetExercisesForUser(reqContext context.Context, params *models.GetExerciseForUserParams) ([]models.Exercise, error)
}

type Service struct {
	repo repo.ExerciseRepository
}

func NewService(r repo.ExerciseRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateExercise(reqContext context.Context, dto *models.CreateExerciseForUserParams) error {
	_, err := s.repo.CreateExercise(reqContext, &repo.CreateExerciseParams{
		Name:         dto.Name,
		Description:  dto.Description,
		TargetMuscle: dto.TargetMuscle,
		PictureURL:   dto.PictureURL,
		UserID:       dto.UserID,
	})
	return err
}

func (s *Service) GetExercisesForUser(reqContext context.Context, dto *models.GetExerciseForUserParams) ([]models.Exercise, error) {
	exercises, err := s.repo.GetExercisesForUser(reqContext, &repo.GetExerciseForUserParams{
		UserID: dto.UserID,
		Offset: dto.Offset,
		Limit:  dto.Limit,
	})
	if err != nil {
		return nil, err
	}
	return exercises, nil
}
