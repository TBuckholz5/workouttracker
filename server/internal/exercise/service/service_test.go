package service

import (
	"context"
	"errors"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/exercise/models"
	"github.com/TBuckholz5/workouttracker/internal/exercise/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExerciserepository struct {
	mock.Mock
}

func (m *mockExerciserepository) CreateExercise(ctx context.Context, params *repository.CreateExerciseParams) (models.Exercise, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.Exercise), args.Error(1)
}

func (m *mockExerciserepository) GetExercisesForUser(ctx context.Context, params *repository.GetExerciseForUserParams) ([]models.Exercise, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.Exercise), args.Error(1)
}

func TestCreateExercise_Success(t *testing.T) {
	mockrepository := new(mockExerciserepository)
	req := &models.CreateExerciseForUserParams{
		UserID:       1,
		Name:         "Bench Press",
		Description:  "Chest exercise",
		TargetMuscle: "Chest",
		PictureURL:   "",
	}
	expected := models.Exercise{ID: 1, Name: req.Name, Description: req.Description, TargetMuscle: req.TargetMuscle}
	mockrepository.On("CreateExercise", mock.Anything, mock.MatchedBy(func(p *repository.CreateExerciseParams) bool {
		return p.Name == req.Name && p.Description == req.Description && p.TargetMuscle == req.TargetMuscle && p.UserID == req.UserID
	})).Return(expected, nil)

	svc := NewService(mockrepository)
	result, err := svc.CreateExercise(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	mockrepository.AssertNumberOfCalls(t, "CreateExercise", 1)
}

func TestCreateExercise_repositoryError(t *testing.T) {
	mockrepository := new(mockExerciserepository)
	req := &models.CreateExerciseForUserParams{
		UserID:       1,
		Name:         "Bench Press",
		Description:  "Chest exercise",
		TargetMuscle: "Chest",
		PictureURL:   "",
	}
	mockrepository.On("CreateExercise", mock.Anything, mock.Anything).Return(models.Exercise{}, errors.New("db error"))

	svc := NewService(mockrepository)
	result, err := svc.CreateExercise(context.Background(), req)
	assert.NotNil(t, err)
	assert.Equal(t, models.Exercise{}, result)
	mockrepository.AssertNumberOfCalls(t, "CreateExercise", 1)
}

func TestGetExercisesForUser_Success(t *testing.T) {
	mockrepository := new(mockExerciserepository)
	req := &models.GetExerciseForUserParams{
		UserID: 1,
		Offset: 0,
		Limit:  10,
	}
	expected := []models.Exercise{{ID: 1, Name: "Bench Press"}}
	mockrepository.On("GetExercisesForUser", mock.Anything, mock.MatchedBy(func(p *repository.GetExerciseForUserParams) bool {
		return p.UserID == req.UserID && p.Offset == req.Offset && p.Limit == req.Limit
	})).Return(expected, nil)

	svc := NewService(mockrepository)
	result, err := svc.GetExercisesForUser(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	mockrepository.AssertNumberOfCalls(t, "GetExercisesForUser", 1)
}

func TestGetExercisesForUser_EmptyResult(t *testing.T) {
	mockrepository := new(mockExerciserepository)
	req := &models.GetExerciseForUserParams{
		UserID: 1,
		Offset: 0,
		Limit:  10,
	}
	mockrepository.On("GetExercisesForUser", mock.Anything, mock.Anything).Return([]models.Exercise{}, nil)

	svc := NewService(mockrepository)
	result, err := svc.GetExercisesForUser(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, []models.Exercise{}, result)
	mockrepository.AssertNumberOfCalls(t, "GetExercisesForUser", 1)
}

func TestGetExercisesForUser_repositoryError(t *testing.T) {
	mockrepository := new(mockExerciserepository)
	req := &models.GetExerciseForUserParams{
		UserID: 1,
		Offset: 0,
		Limit:  10,
	}
	mockrepository.On("GetExercisesForUser", mock.Anything, mock.Anything).Return([]models.Exercise{}, errors.New("db error"))

	svc := NewService(mockrepository)
	result, err := svc.GetExercisesForUser(context.Background(), req)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	mockrepository.AssertNumberOfCalls(t, "GetExercisesForUser", 1)
}
