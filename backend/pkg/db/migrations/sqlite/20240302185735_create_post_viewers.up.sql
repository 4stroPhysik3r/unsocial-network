CREATE TABLE
    IF NOT EXISTS post_viewers (
        post_id INTEGER NOT NULL,
        viewer_id INTEGER NOT NULL,
        PRIMARY KEY (post_id, viewer_id),
        FOREIGN KEY (post_id) REFERENCES posts (post_id),
        FOREIGN KEY (viewer_id) REFERENCES users (user_id)
    );