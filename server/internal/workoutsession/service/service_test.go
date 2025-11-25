package service

import (
	"context"
	"errors"
	"testing"

	"github.com/TBuckholz5/workouttracker/internal/workoutsession/models"
	"github.com/TBuckholz5/workouttracker/internal/workoutsession/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWorkoutSessionRepository struct {
	mock.Mock
}

func (m *MockWorkoutSessionRepository) Create(ctx context.Context, session *models.WorkoutSession) (*repository.WorkoutSession, []*repository.Workout, []*repository.WorkoutSet, error) {
	args := m.Called(ctx, session)
	return args.Get(0).(*repository.WorkoutSession),
		args.Get(1).([]*repository.Workout),
		args.Get(2).([]*repository.WorkoutSet),
		args.Error(3)
}

func TestService_Create_Success(t *testing.T) {
	mockRepo := new(MockWorkoutSessionRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	inputSession := &models.WorkoutSession{
		Name:        "Morning Workout",
		UserID:      42,
		Description: "Upper body workout",
		Duration:    60,
		Workouts: []models.Workout{
			{
				ExerciseID:  1,
				Description: "Bench press",
				Sets: []models.WorkoutSet{
					{
						Reps:     10,
						Weight:   135.5,
						SetType:  "working",
						SetOrder: 1,
					},
				},
			},
		},
	}

	expectedRepoSession := &repository.WorkoutSession{
		ID:          1,
		Name:        "Morning Workout",
		UserID:      42,
		Description: "Upper body workout",
		Duration:    60,
	}

	expectedRepoWorkouts := []*repository.Workout{
		{
			ID:          1,
			ExerciseID:  1,
			Description: "Bench press",
			SessionId:   1,
		},
	}

	expectedRepoSets := []*repository.WorkoutSet{
		{
			ID:        1,
			Reps:      10,
			Weight:    135.5,
			SetType:   "working",
			SetOrder:  1,
			WorkoutID: 1,
		},
	}

	mockRepo.On("Create", ctx, inputSession).Return(expectedRepoSession, expectedRepoWorkouts, expectedRepoSets, nil)

	result, err := service.Create(ctx, inputSession)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Morning Workout", result.Name)
	assert.Equal(t, int64(42), result.UserID)
	assert.Equal(t, "Upper body workout", result.Description)
	assert.Equal(t, 60, result.Duration)
	assert.Len(t, result.Workouts, 1)
	assert.Equal(t, int64(1), result.Workouts[0].ID)
	assert.Equal(t, int64(1), result.Workouts[0].ExerciseID)
	assert.Equal(t, "Bench press", result.Workouts[0].Description)
	assert.Len(t, result.Workouts[0].Sets, 1)
	assert.Equal(t, int64(1), result.Workouts[0].Sets[0].ID)
	assert.Equal(t, 10, result.Workouts[0].Sets[0].Reps)
	assert.Equal(t, 135.5, result.Workouts[0].Sets[0].Weight)
	assert.Equal(t, "working", result.Workouts[0].Sets[0].SetType)
	assert.Equal(t, 1, result.Workouts[0].Sets[0].SetOrder)

	mockRepo.AssertExpectations(t)
}

func TestService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(MockWorkoutSessionRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	inputSession := &models.WorkoutSession{
		Name:        "Morning Workout",
		UserID:      42,
		Description: "Upper body workout",
		Duration:    60,
	}

	mockRepo.On("Create", ctx, inputSession).Return(
		(*repository.WorkoutSession)(nil),
		([]*repository.Workout)(nil),
		([]*repository.WorkoutSet)(nil),
		errors.New("database connection failed"))

	result, err := service.Create(ctx, inputSession)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection failed")
	mockRepo.AssertExpectations(t)
}

func TestService_Create_EmptyWorkouts(t *testing.T) {
	mockRepo := new(MockWorkoutSessionRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	inputSession := &models.WorkoutSession{
		Name:        "Empty Workout",
		UserID:      42,
		Description: "No workouts",
		Duration:    30,
		Workouts:    []models.Workout{},
	}

	expectedRepoSession := &repository.WorkoutSession{
		ID:          1,
		Name:        "Empty Workout",
		UserID:      42,
		Description: "No workouts",
		Duration:    30,
	}

	mockRepo.On("Create", ctx, inputSession).Return(
		expectedRepoSession,
		[]*repository.Workout{},
		[]*repository.WorkoutSet{},
		nil)

	result, err := service.Create(ctx, inputSession)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Empty Workout", result.Name)
	assert.Len(t, result.Workouts, 0)
	mockRepo.AssertExpectations(t)
}
