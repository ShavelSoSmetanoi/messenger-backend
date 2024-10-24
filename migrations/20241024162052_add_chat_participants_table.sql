-- +goose Up
CREATE TABLE IF NOT EXISTS ChatParticipants (
    chat_id INT NOT NULL,
    user_id INT NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES Chats(id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- +goose Down
DROP TABLE IF EXISTS ChatParticipants;
