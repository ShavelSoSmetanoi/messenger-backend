-- +goose Up
CREATE TABLE IF NOT EXISTS Messages (
        id SERIAL PRIMARY KEY,
        chat_id INT NOT NULL,
        user_id INT NOT NULL,
        type VARCHAR(10) NOT NULL DEFAULT 'text',
        content TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        is_read BOOLEAN DEFAULT FALSE,
        read_at TIMESTAMP,
        FOREIGN KEY (chat_id) REFERENCES Chats(id),
        FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- +goose Down
DROP TABLE IF EXISTS Messages;
