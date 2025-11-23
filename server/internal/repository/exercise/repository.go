package exercise

import (
	"context"
	"fmt"

	serviceModels "github.com/TBuckholz5/workouttracker/internal/service/exercise/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateExerciseParams struct {
	Name         string
	Description  string
	TargetMuscle string
	PictureURL   string
	UserID       int64
}

type GetExerciseForUserParams struct {
	UserID int64
	Limit  int
	Offset int
}

type ExerciseRepository interface {
	CreateExercise(ctx context.Context, params *CreateExerciseParams) (serviceModels.Exercise, error)
	GetExercisesForUser(ctx context.Context, params *GetExerciseForUserParams) ([]serviceModels.Exercise, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) CreateExercise(ctx context.Context, params *CreateExerciseParams) (serviceModels.Exercise, error) {
	var exercise exercise
	err := r.pool.QueryRow(ctx, createExerciseQuery,
		params.Name,
		params.Description,
		params.TargetMuscle,
		params.PictureURL,
		params.UserID,
	).Scan(
		&exercise.id,
		&exercise.name,
		&exercise.description,
		&exercise.targetMuscle,
		&exercise.pictureUrl,
		&exercise.createdAt,
		&exercise.updatedAt,
		&exercise.userId,
	)
	if err != nil {
		return serviceModels.Exercise{}, fmt.Errorf("error creating exercise: %w", err)
	}
	return serviceModels.Exercise{
		ID:           exercise.id,
		Name:         exercise.name,
		Description:  exercise.description,
		TargetMuscle: exercise.targetMuscle,
		PictureURL:   exercise.pictureUrl,
	}, nil
}

func (r *Repository) GetExercisesForUser(ctx context.Context, params *GetExerciseForUserParams) ([]serviceModels.Exercise, error) {
	rows, err := r.pool.Query(ctx, getExercisesForUserQuery, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching exercises for user: %w", err)
	}
	defer rows.Close()

	exercises := make([]serviceModels.Exercise, 0)
	for rows.Next() {
		var exercise exercise
		err := rows.Scan(
			&exercise.id,
			&exercise.name,
			&exercise.description,
			&exercise.targetMuscle,
			&exercise.pictureUrl,
			&exercise.createdAt,
			&exercise.updatedAt,
			&exercise.userId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning exercise row: %w", err)
		}
		exercises = append(exercises, serviceModels.Exercise{
			ID:           exercise.id,
			Name:         exercise.name,
			Description:  exercise.description,
			TargetMuscle: exercise.targetMuscle,
			PictureURL:   exercise.pictureUrl,
		})
	}
	return exercises, nil
}
