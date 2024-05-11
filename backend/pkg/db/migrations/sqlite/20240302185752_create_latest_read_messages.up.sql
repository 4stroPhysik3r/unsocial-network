CREATE TABLE
    IF NOT EXISTS latest_read_messages (
        chat_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        message_id INTEGER NOT NULL,
        PRIMARY KEY (chat_id, user_id),
        FOREIGN KEY (chat_id) REFERENCES chats (chat_id),
        FOREIGN KEY (user_id) REFERENCES users (user_id),
        FOREIGN KEY (message_id) REFERENCES messages (message_id)
    );