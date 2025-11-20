package exercise

const createExerciseQuery = `INSERT INTO exercises (name, description, target_muscle, picture_url, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, name, description, target_muscle, picture_url, created_at, updated_at, user_id;`

const getExercisesForUserQuery = `SELECT id, name, description, target_muscle, picture_url, created_at, updated_at, user_id
	FROM exercises WHERE user_id = $1
	LIMIT $2 OFFSET $3;`
