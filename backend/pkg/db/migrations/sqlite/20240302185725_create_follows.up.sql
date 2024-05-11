CREATE TABLE
    IF NOT EXISTS follows (
        follower_id INTEGER NOT NULL, --sending following request
        following_id INTEGER NOT NULL, -- is being followed by the follower
        status TEXT CHECK (status IN ('pending', 'accepted', 'rejected')) NOT NULL DEFAULT 'pending',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (follower_id, following_id),
        FOREIGN KEY (follower_id) REFERENCES users (user_id),
        FOREIGN KEY (following_id) REFERENCES users (user_id)
    );