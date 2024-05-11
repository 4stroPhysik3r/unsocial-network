CREATE TABLE IF NOT EXISTS group_members (
    group_id INTEGER,
    user_id INTEGER,
    status TEXT CHECK (status IN ('invited', 'accepted', 'request', 'rejected')) NOT NULL DEFAULT 'request',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES groups (group_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);
