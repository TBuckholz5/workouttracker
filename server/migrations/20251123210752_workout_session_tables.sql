-- +goose Up
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name text NOT NULL,
    description TEXT,
    duration INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE workouts (
    id BIGSERIAL PRIMARY KEY,
    exercise_id BIGINT REFERENCES exercises(id),
    session_id BIGINT REFERENCES sessions(id),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE set_type as ENUM('dropset', 'superset', 'normal', 'failure');

CREATE TABLE workout_sets (
    id BIGSERIAL PRIMARY KEY,
    workout_id BIGINT REFERENCES workouts(id) ON DELETE CASCADE,
    reps INT NOT NULL,
    weight DECIMAL(6,2),
    set_type set_type NOT NULL DEFAULT 'normal',
    set_order INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE sessions CASCADE;
DROP TABLE workouts CASCADE;
DROP TABLE workout_sets CASCADE;
DROP TYPE set;
