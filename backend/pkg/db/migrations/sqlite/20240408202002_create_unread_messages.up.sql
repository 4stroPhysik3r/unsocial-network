CREATE TABLE IF NOT EXISTS unread_messages (
    user_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (chat_id) REFERENCES chats (chat_id)
);
