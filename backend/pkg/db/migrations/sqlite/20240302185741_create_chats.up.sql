CREATE TABLE
    IF NOT EXISTS chats (
        chat_id INTEGER PRIMARY KEY AUTOINCREMENT,
        group_id INTEGER,
        FOREIGN KEY (group_id) REFERENCES groups (group_id)
    );