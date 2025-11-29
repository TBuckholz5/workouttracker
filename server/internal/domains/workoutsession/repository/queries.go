package repository

const createSessionQuery = `INSERT INTO sessions (name, user_id, description, duration)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, user_id, description, duration, created_at, updated_at;`

const createWorkoutQuery = `INSERT INTO workouts (exercise_id, description, session_id)
	VALUES ($1, $2, $3)
	RETURNING id, exercise_id, description, session_id, created_at, updated_at;`

const createSetQuery = `INSERT INTO workout_sets (workout_id, reps, weight, set_type, set_order)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, reps, weight, set_type, set_order, workout_id, created_at, updated_at;`
