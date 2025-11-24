-- +goose Up
ALTER TABLE workouts
ADD COLUMN session_id BIGINT REFERENCES sessions(id);

-- +goose Down
ALTER TABLE workouts
DROP COLUMN session_id;
