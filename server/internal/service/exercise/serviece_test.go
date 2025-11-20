package exercise

import (
	"context"
	"errors"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/api/v1/exercise/dto"
	repo "github.com/TBuckholz5/workouttracker/internal/repository/exercise"
	"github.com/TBuckholz5/workouttracker/internal/service/exercise/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExerciseRepo struct {
	mock.Mock
}

func (m *mockExerciseRepo) CreateExercise(ctx context.Context, params *repo.CreateExerciseParams) (models.Exercise, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.Exercise), args.Error(1)
}

func (m *mockExerciseRepo) GetExercisesForUser(ctx context.Context, params *repo.GetExerciseForUserParams) ([]models.Exercise, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.Exercise), args.Error(1)
}

func TestCreateExercise_Success(t *testing.T) {
	mockRepo := new(mockExerciseRepo)
	req := &dto.CreateExerciseRequest{
		Name:         "Bench Press",
		Description:  "Chest exercise",
		TargetMuscle: "Chest",
	}
	mockRepo.On("CreateExercise", mock.Anything, mock.MatchedBy(func(p *repo.CreateExerciseParams) bool {
		return p.Name == req.Name && p.Description == req.Description && p.TargetMuscle == req.TargetMuscle
	})).Return(models.Exercise{ID: 1, Name: req.Name}, nil)

	svc := NewService(mockRepo)
	err := svc.CreateExercise(context.Background(), req)
	assert.Nil(t, err)
	mockRepo.AssertNumberOfCalls(t, "CreateExercise", 1)
}

func TestCreateExercise_RepoError(t *testing.T) {
	mockRepo := new(mockExerciseRepo)
	req := &dto.CreateExerciseRequest{
		Name:         "Bench Press",
		Description:  "Chest exercise",
		TargetMuscle: "Chest",
	}
	mockRepo.On("CreateExercise", mock.Anything, mock.Anything).Return(models.Exercise{}, errors.New("db error"))

	svc := NewService(mockRepo)
	err := svc.CreateExercise(context.Background(), req)
	assert.NotNil(t, err)
	mockRepo.AssertNumberOfCalls(t, "CreateExercise", 1)
}

func TestGetExercisesForUser_Success(t *testing.T) {
	mockRepo := new(mockExerciseRepo)
	req := &dto.GetExerciseForUserRequest{
		UserID: 1,
		Offset: 0,
		Limit:  10,
	}
	expected := []models.Exercise{{ID: 1, Name: "Bench Press"}}
	mockRepo.On("GetExercisesForUser", mock.Anything, mock.MatchedBy(func(p *repo.GetExerciseForUserParams) bool {
		return p.UserID == req.UserID && p.Offset == req.Offset && p.Limit == req.Limit
	})).Return(expected, nil)

	svc := NewService(mockRepo)
	result, err := svc.GetExercisesForUser(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertNumberOfCalls(t, "GetExercisesForUser", 1)
}

func TestGetExercisesForUser_RepoError(t *testing.T) {
	mockRepo := new(mockExerciseRepo)
	req := &dto.GetExerciseForUserRequest{
		UserID: 1,
		Offset: 0,
		Limit:  10,
	}
	mockRepo.On("GetExercisesForUser", mock.Anything, mock.Anything).Return([]models.Exercise{}, errors.New("db error"))

	svc := NewService(mockRepo)
	result, err := svc.GetExercisesForUser(context.Background(), req)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	mockRepo.AssertNumberOfCalls(t, "GetExercisesForUser", 1)
}
