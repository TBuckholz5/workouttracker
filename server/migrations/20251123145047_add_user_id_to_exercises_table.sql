-- +goose Up
ALTER TABLE exercises
ADD COLUMN user_id BIGINT REFERENCES users(id);

-- +goose Down
ALTER TABLE exercises
DROP COLUMN user_id;
