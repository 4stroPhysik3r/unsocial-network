package db

import (
	"database/sql"
	"fmt"
)

func GetUserIDFromSessionToken(token string) (int, error) {
	var userID int
	err := DB.QueryRow("SELECT user_id FROM sessions WHERE session_key = ?", token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("session token not found")
		}
		return 0, err
	}
	return userID, nil
}

func UserDataFromID(userID int) (User, error) {
	var user User
	query := `SELECT * FROM users WHERE user_id = ?`

	err := DB.QueryRow(query, userID).Scan(&user.UserID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.ProfilePublic, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}

func UpdateUserProfile(userID int, profilePublic bool) error {
	query := `"UPDATE users SET profile_public = ? WHERE user_id = ?"`

	_, err := DB.Exec(query, profilePublic, userID)
	if err != nil {
		return fmt.Errorf("error updating profile_public field: %w", err)
	}
	return nil
}
