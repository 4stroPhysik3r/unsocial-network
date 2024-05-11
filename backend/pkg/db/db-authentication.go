package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateUser(creds LoginCredentials) (bool, int) {
	var userID int
	var hashedPassword string

	err := DB.QueryRow("SELECT user_id, password FROM users WHERE email = ?", creds.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user with the email %s\n", creds.Email)
			return false, 0
		}
		log.Printf("Database error: %v\n", err)
		return false, 0
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		return false, 0
	}

	return true, userID
}

func StoreSession(sessionKey string, userID int) error {
	// Check if a session for the user already exists
	var existingSessionID int
	err := DB.QueryRow("SELECT session_id FROM sessions WHERE user_id = ?", userID).Scan(&existingSessionID)

	if err == nil {
		// If a session exists, delete it
		_, delErr := DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if delErr != nil {
			log.Printf("Error deleting existing session for user %d: %v", userID, delErr)
			return delErr
		}
	} else if err != sql.ErrNoRows {
		// If we encounter an error other than ErrNoRows, report it
		log.Printf("Error checking for existing session for user %d: %v", userID, err)
		return err
	}

	// No existing session or deletion successful, proceed to insert the new session
	_, insErr := DB.Exec("INSERT INTO sessions (session_key, user_id) VALUES (?, ?)", sessionKey, userID)
	if insErr != nil {
		log.Printf("Error inserting new session into database for user %d: %v", userID, insErr)
		return insErr
	}

	return nil
}

func ValidateSessionToken(sessionToken string) bool {
	var isValid bool

	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE session_key = ? LIMIT 1)`
	err := DB.QueryRow(query, sessionToken).Scan(&isValid)
	if err != nil {
		log.Printf("Error checking session token: %v", err)
		return false
	}

	return isValid
}

func DeleteSession(sessionToken string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE session_key = ?", sessionToken)
	return err
}
