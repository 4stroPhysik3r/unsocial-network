CREATE TABLE
    IF NOT EXISTS posts (
        post_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        group_id INTEGER,
        content TEXT NOT NULL,
        post_image TEXT,
        privacy_level VARCHAR CHECK (privacy_level IN ('private', 'public', 'friends')) NOT NULL DEFAULT 'private',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (user_id),
        FOREIGN KEY (group_id) REFERENCES groups (group_id)
    );