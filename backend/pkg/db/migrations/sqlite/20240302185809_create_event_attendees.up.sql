CREATE TABLE
    IF NOT EXISTS event_attendees (
        attendee_id INTEGER NOT NULL,
        event_id INTEGER NOT NULL,
        status TEXT CHECK (status IN ('going', 'not going', 'not replied')) NOT NULL DEFAULT 'not replied',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (attendee_id, event_id),
        FOREIGN KEY (attendee_id) REFERENCES users (user_id),
        FOREIGN KEY (event_id) REFERENCES events (event_id)
    );