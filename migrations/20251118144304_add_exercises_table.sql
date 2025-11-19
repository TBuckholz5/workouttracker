-- +goose Up
CREATE TABLE exercises (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    description TEXT,
    target_muscle TEXT,
    picture_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id BIGINT REFERENCES users(id)
);

-- +goose Down
DROP TABLE exercises;
