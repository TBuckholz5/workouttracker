package service

import (
	"context"

	"github.com/TBuckholz5/workouttracker/internal/domains/exercise/models"
	repo "github.com/TBuckholz5/workouttracker/internal/domains/exercise/repository"
)

type ExerciseService interface {
	CreateExercise(reqContext context.Context, params *CreateExerciseForUserParams) (models.Exercise, error)
	GetExercisesForUser(reqContext context.Context, params *GetExerciseForUserParams) ([]models.Exercise, error)
}

type Service struct {
	repo repo.ExerciseRepository
}

func NewService(r repo.ExerciseRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateExercise(reqContext context.Context, params *CreateExerciseForUserParams) (models.Exercise, error) {
	return s.repo.CreateExercise(reqContext, &repo.CreateExerciseParams{
		Name:         params.Name,
		Description:  params.Description,
		TargetMuscle: params.TargetMuscle,
		PictureURL:   params.PictureURL,
		UserID:       params.UserID,
	})
}

func (s *Service) GetExercisesForUser(reqContext context.Context, params *GetExerciseForUserParams) ([]models.Exercise, error) {
	exercises, err := s.repo.GetExercisesForUser(reqContext, &repo.GetExerciseForUserParams{
		UserID: params.UserID,
		Offset: params.Offset,
		Limit:  params.Limit,
	})
	if err != nil {
		return nil, err
	}
	return exercises, nil
}
