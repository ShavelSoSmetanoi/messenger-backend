-- +goose Up
CREATE TABLE IF NOT EXISTS Chats (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE Chats;
