CREATE TABLE IF NOT EXISTS groups (
    group_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    -- status TEXT CHECK (status IN ('invited', 'accepted', 'request')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- PRIMARY KEY (group_id, user_id),
    -- FOREIGN KEY (group_id) REFERENCES groups (group_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);