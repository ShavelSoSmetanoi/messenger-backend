-- +goose Up
CREATE TABLE IF NOT EXISTS user_tokens (
        id SERIAL PRIMARY KEY,
        user_id VARCHAR(255) NOT NULL,
        token TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS user_tokens;