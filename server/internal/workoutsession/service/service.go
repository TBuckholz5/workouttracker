package service

import (
	"context"

	"github.com/TBuckholz5/workouttracker/internal/workoutsession/models"
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/repository"
)

type WorkoutSessionService interface {
	Create(reqContext context.Context, session *models.WorkoutSession) (*models.WorkoutSession, error)
}

type Service struct {
	repo repository.WorkoutSessionRepository
}

func NewService(r repository.WorkoutSessionRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(reqContext context.Context, session *models.WorkoutSession) (*models.WorkoutSession, error) {
	repositorySession, repositoryWorkouts, repositorySets, err := s.repo.Create(reqContext, session)
	if err != nil {
		return nil, err
	}
	serviceSession := repositoryToModels(repositorySession, repositoryWorkouts, repositorySets)
	return serviceSession, nil
}
