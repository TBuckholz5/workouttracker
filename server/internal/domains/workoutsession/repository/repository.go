package repository

import (
	"context"
	"fmt"

	"github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutSessionRepository interface {
	Create(ctx context.Context, session *models.WorkoutSession) (*WorkoutSession, []*Workout, []*WorkoutSet, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

type WorkoutParam = map[*Workout][]*WorkoutSet

func (r *Repository) Create(ctx context.Context, session *models.WorkoutSession) (*WorkoutSession, []*Workout, []*WorkoutSet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var newSession WorkoutSession
	err = tx.QueryRow(ctx, createSessionQuery, session.Name, session.UserID, session.Description, session.Duration).Scan(
		&newSession.ID,
		&newSession.Name,
		&newSession.UserID,
		&newSession.Description,
		&newSession.Duration,
		&newSession.CreatedAt,
		&newSession.UpdatedAt,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	newWorkouts := make([]*Workout, len(session.Workouts))
	var newSets []*WorkoutSet
	for i, w := range session.Workouts {
		var workout Workout

		err := tx.QueryRow(ctx, createWorkoutQuery, w.ExerciseID, w.Description, newSession.ID).Scan(
			&workout.ID,
			&workout.ExerciseID,
			&workout.Description,
			&workout.SessionId,
			&workout.CreatedAt,
			&workout.UpdatedAt,
		)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create workout: %w", err)
		}
		newWorkouts[i] = &workout

		for _, s := range w.Sets {
			var set WorkoutSet
			err := tx.QueryRow(ctx, createSetQuery, workout.ID, s.Reps, s.Weight, s.SetType, s.SetOrder).Scan(
				&set.ID,
				&set.Reps,
				&set.Weight,
				&set.SetType,
				&set.SetOrder,
				&set.WorkoutID,
				&set.CreatedAt,
				&set.UpdatedAt,
			)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to create workout set: %w", err)
			}
			newSets = append(newSets, &set)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &newSession, newWorkouts, newSets, nil
}
