package service

import (
	"context"

	"github.com/TBuckholz5/workouttracker/internal/exercise/models"
	repo "github.com/TBuckholz5/workouttracker/internal/exercise/repository"
)

type ExerciseService interface {
	CreateExercise(reqContext context.Context, params *models.CreateExerciseForUserParams) (models.Exercise, error)
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

func (s *Service) CreateExercise(reqContext context.Context, dto *models.CreateExerciseForUserParams) (models.Exercise, error) {
	return s.repo.CreateExercise(reqContext, &repo.CreateExerciseParams{
		Name:         dto.Name,
		Description:  dto.Description,
		TargetMuscle: dto.TargetMuscle,
		PictureURL:   dto.PictureURL,
		UserID:       dto.UserID,
	})
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
