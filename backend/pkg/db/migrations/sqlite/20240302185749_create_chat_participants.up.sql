CREATE TABLE
    IF NOT EXISTS chat_participants (
        chat_id INTEGER NOT NULL,
        participant_id INTEGER NOT NULL,
        PRIMARY KEY (chat_id, participant_id),
        FOREIGN KEY (chat_id) REFERENCES chats (chat_id),
        FOREIGN KEY (participant_id) REFERENCES users (user_id)
    );