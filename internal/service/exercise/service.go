package exercise

import (
	"context"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/exercise/dto"
	repo "github.com/TBuckholz5/workouttracker/internal/repository/exercise"
	"github.com/TBuckholz5/workouttracker/internal/service/exercise/models"
)

type ExerciseService interface {
	CreateExercise(reqContext context.Context, dto *dto.CreateExerciseRequest) error
	GetExercisesForUser(reqContext context.Context, dto *dto.GetExerciseForUserRequest) ([]models.Exercise, error)
}

type Service struct {
	repo repo.ExerciseRepository
}

func NewService(r repo.ExerciseRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreateExercise(reqContext context.Context, dto *dto.CreateExerciseRequest) error {
	_, err := s.repo.CreateExercise(reqContext, &repo.CreateExerciseParams{
		Name:         dto.Name,
		Description:  dto.Description,
		TargetMuscle: dto.TargetMuscle,
	})
	return err
}

func (s *Service) GetExercisesForUser(reqContext context.Context, dto *dto.GetExerciseForUserRequest) ([]models.Exercise, error) {
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
