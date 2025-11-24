-- +goose Up
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    description TEXT,
    duration INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE set_type as ENUM('dropset', 'superset', 'normal', 'failure');

CREATE TYPE set AS (
    reps INT,
    weight FLOAT,
    type set_type
);

CREATE TABLE workouts (
    id BIGSERIAL PRIMARY KEY,
    exercise_id BIGINT REFERENCES exercises(id),
    description TEXT,
    sets set[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE sessions;
DROP TABLE workouts;
DROP TYPE set;
DROP TYPE set_type;
