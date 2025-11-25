package repository

import (
	"time"
)

type WorkoutSet struct {
	ID        int64     `db:"id"`
	WorkoutID int64     `db:"workout_id"`
	Reps      int       `db:"reps"`
	Weight    float64   `db:"weight"`
	SetType   string    `db:"set_type"`
	SetOrder  int       `db:"set_order"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Workout struct {
	ID          int64        `db:"id"`
	ExerciseID  int64        `db:"exercise_id"`
	SessionId   int64        `db:"session_id"`
	Description string       `db:"description"`
	Sets        []WorkoutSet `db:"sets"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}

type WorkoutSession struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	UserID      int64     `db:"user_id"`
	Description string    `db:"description"`
	Duration    int       `db:"duration"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
