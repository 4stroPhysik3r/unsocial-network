package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Event struct {
	EventID          int    `json:"event_id,omitempty"`
	UserID           int    `json:"user_id"`
	GroupID          int    `json:"group_id"`
	Date             string `json:"date"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	CreatorName      string `json:"creator_name"`
	CreatorFirstname string
	CreatorLastname  string
	Members          []int     `json:"members"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type AttendeesStatus struct {
	AttendeesStatus string `json:"attendees_status"`
}

func InsertEvent(eventData Event) error {
	// Prepare the SQL statement
	statement, err := DB.Prepare(`
		INSERT INTO events (group_id, user_id, title, content, date)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer statement.Close()

	// Execute the SQL statement
	result, err := statement.Exec(eventData.GroupID, eventData.UserID, eventData.Title, eventData.Content, eventData.Date)
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %v", err)
	}

	// Get the ID of the newly inserted row
	eventID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}
	fmt.Println("Successfully added event with ID: ", eventID)

	return nil
}

func GetEvents(groupID int) ([]Event, error) {
	query := `
        SELECT e.event_id, e.user_id, e.group_id, e.date, e.title, e.content, 
               CONCAT(u.firstname, ' ', u.lastname) AS creator_name, e.created_at
        FROM events e
        INNER JOIN users u ON e.user_id = u.user_id
        WHERE e.group_id = ?`

	rows, err := DB.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("GetEvents: failed to query events: %v", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.EventID, &event.UserID, &event.GroupID, &event.Date, &event.Title, &event.Content, &event.CreatorName, &event.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetEvents: failed to scan event row: %v", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetEvents: error iterating over event rows: %v", err)
	}

	return events, nil
}

func UpdateAttendeesStatus(attendeeID int, eventID int, status string) error {
	// Check if the attendee's status already exists in the database
	var existingStatus string
	err := DB.QueryRow("SELECT status FROM event_attendees WHERE attendee_id = ? AND event_id = ?", attendeeID, eventID).Scan(&existingStatus)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing status: %w", err)
	}

	// Decide whether to insert or update based on the existence of the record
	if err == sql.ErrNoRows {
		// If the record doesn't exist, insert a new one
		_, err := DB.Exec("INSERT INTO event_attendees (attendee_id, event_id, status) VALUES (?, ?, ?)", attendeeID, eventID, status)
		if err != nil {
			return fmt.Errorf("failed to insert attendee's status: %w", err)
		}
	} else {
		// If the record already exists, update the status
		_, err := DB.Exec("UPDATE event_attendees SET status = ? WHERE attendee_id = ? AND event_id = ?", status, attendeeID, eventID)
		if err != nil {
			return fmt.Errorf("failed to update attendee's status: %w", err)
		}
	}

	return nil
}

func GetAttendeesStatus(attendeeID int, eventID int) (string, error) {
	var status string
	err := DB.QueryRow("SELECT status FROM event_attendees WHERE attendee_id = ? AND event_id = ?", attendeeID, eventID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return "nil" string if no rows are found
			return "nil", nil
		}
		return "", fmt.Errorf("failed to fetch status: %w", err)
	}

	return status, nil
}
