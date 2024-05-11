package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Follower struct {
	FollowerID  int    `json:"follower_id"`
	FollowingID int    `json:"following_id"`
	FullName    string `json:"full_name"`
}

func GetFollowing(userID int) ([]Follower, error) {
	var followings []Follower
	query := `SELECT f.follower_id, f.following_id, u.firstname || ' ' || u.lastname AS full_name FROM follows f
			  JOIN users u ON f.following_id = u.user_id
			  WHERE f.follower_id = ? AND f.status = 'accepted'`

	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Printf("Error querying following: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var following Follower
		if err := rows.Scan(&following.FollowingID, &following.FollowerID, &following.FullName); err != nil {
			log.Printf("Error scanning follower: %v", err)
			return nil, err
		}
		followings = append(followings, following)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating following: %v", err)
		return nil, err
	}

	return followings, nil
}

func GetFollowers(userID int) ([]Follower, error) {
	var followers []Follower

	query := `SELECT f.following_id, u.firstname || ' ' || u.lastname AS full_name FROM follows f
			  JOIN users u ON f.follower_id = u.user_id
			  WHERE f.following_id = ? AND f.status = 'accepted'`

	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Printf("Error querying followers: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower Follower
		if err := rows.Scan(&follower.FollowerID, &follower.FullName); err != nil {
			log.Printf("Error scanning follower: %v", err)
			return nil, err
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating followers: %v", err)
		return nil, err
	}

	return followers, nil
}

func GetFollowStatus(followerID, followingID int) (string, error) {
	var status string

	err := DB.QueryRow("SELECT status FROM follows WHERE follower_id = ? AND following_id = ?", followerID, followingID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			status = "rejected"
			return status, nil
		}
		return "", fmt.Errorf("GetFollowStatus: error querying follow status: %w", err)
	}
	return status, nil
}

func FollowUserRequest(followerID, followingID int) error {
	// Check if the following person has a public profile
	var profilePublic bool
	err := DB.QueryRow("SELECT profile_public FROM users WHERE user_id = ?", followingID).Scan(&profilePublic)
	if err != nil {
		log.Printf("Error checking profile public status: %v", err)
		return err
	}

	if !profilePublic {

		query := `
		INSERT INTO follows (follower_id, following_id, status) 
		VALUES (?, ?, 'pending') 
		ON CONFLICT (follower_id, following_id) 
		DO UPDATE SET status = 'pending';`

		_, err = DB.Exec(query, followerID, followingID)
		if err != nil {
			log.Printf("Failed to insert or update follow request to pending: %v", err)
			return err
		}

		err = CreateFollowRequestNotification(followerID, followingID)
		if err != nil {
			log.Printf("Error creating follower request notification: %v", err)
			return err
		}

	} else {
		query := `
		INSERT INTO follows (follower_id, following_id, status) 
		VALUES (?, ?, 'accepted') 
		ON CONFLICT (follower_id, following_id) 
		DO UPDATE SET status = 'accepted';`

		_, err = DB.Exec(query, followerID, followingID)
		if err != nil {
			log.Printf("Failed to insert or update follow request to pending: %v", err)
			return err
		}
		// Check if a chat already exists between the two users
		var existingChatID int
		query = `
		SELECT c.chat_id 
		FROM chats c
		JOIN chat_participants cp1 ON c.chat_id = cp1.chat_id AND cp1.participant_id = ?
		JOIN chat_participants cp2 ON c.chat_id = cp2.chat_id AND cp2.participant_id = ?
		WHERE c.group_id IS NULL
		LIMIT 1;
`
		err := DB.QueryRow(query, followerID, followingID).Scan(&existingChatID)
		if err == nil {
			// A chat already exists, so no need to create a new one
			log.Printf("Chat already exists between user %d and user %d with chat ID %d", followerID, followingID, existingChatID)
			return nil
		} else if err != sql.ErrNoRows {
			// An actual error occurred
			log.Printf("Error checking for existing chat: %v", err)
			return err
		}

		// No existing chat found, proceed to create a new chat
		var chatID int64
		result, err := DB.Exec("INSERT INTO chats DEFAULT VALUES")
		if err != nil {
			log.Printf("Error creating new chat: %v", err)
			return err
		}

		chatID, err = result.LastInsertId()
		if err != nil {
			log.Printf("Error retrieving new chat ID: %v", err)
			return err
		}

		// Add both users to the chat_participants table
		_, err = DB.Exec(`
			INSERT INTO chat_participants (chat_id, participant_id) 
			VALUES (?, ?), (?, ?)
		`, chatID, followerID, chatID, followingID)
		if err != nil {
			log.Printf("Error adding users to chat participants: %v", err)
			return err
		}

	}

	return nil
}

func UnfollowUser(followerID, followingID int) error {
	query := `UPDATE follows SET status = 'rejected' WHERE follower_id = ? AND following_id = ?;`

	_, err := DB.Exec(query, followerID, followingID)
	if err != nil {
		log.Printf("Failed to update follow status to not_following: %v", err)
		return err
	}

	// first we check if there other user still follows
	followStatus, err := GetFollowStatus(followingID, followerID)
	if err != nil {
		log.Printf("Failed to GetFollowStatus: %v", err)
		return err
	}

	// in case the other user is not following, we delete the chat
	if followStatus != "accepted" {
		var chatID int

		// Search for the chat_id where both followerID and followingID are present and it's not a group chat
		chatIDQuery := `
			SELECT cp.chat_id 
			FROM chat_participants cp
			JOIN chats c ON cp.chat_id = c.chat_id
			WHERE cp.participant_id IN (?, ?) 
			AND c.group_id IS NULL
			GROUP BY cp.chat_id
			HAVING COUNT(*) = 2;
				`
		err = DB.QueryRow(chatIDQuery, followerID, followingID).Scan(&chatID)
		if err != nil {
			log.Printf("Failed to find chat ID: %v", err)
			return err
		}

		// Delete that chat from the chats table
		_, err = DB.Exec(`DELETE FROM chats WHERE chat_id = ?;`, chatID)
		if err != nil {
			log.Printf("Failed to delete chat: %v", err)
			return err
		}

		// Delete all messages in the messages table where this chat id is present
		_, err = DB.Exec(`DELETE FROM messages WHERE chat_id = ?;`, chatID)
		if err != nil {
			log.Printf("Failed to delete messages: %v", err)
			return err
		}

		// Delete all records from chat_participants table where this chat id is present
		_, err = DB.Exec(`DELETE FROM chat_participants WHERE chat_id = ?;`, chatID)
		if err != nil {
			log.Printf("Failed to delete chat participants: %v", err)
			return err
		}
	}

	return nil
}

func handleFollowRequestResponse(tx *sql.Tx, followerID, followingID int, status string) error {

	// Update the follows table with the new status
	_, err := tx.Exec("UPDATE follows SET status = ? WHERE follower_id = ? AND following_id = ?", status, followerID, followingID)
	if err != nil {
		log.Printf("Failed to update follows status: %v", err)
		return err
	}

	// If the follow request is accepted, create a chat
	if status == "accepted" {
		// creates chat between two user if it does not exists yet
		return createChatForUsers(tx, followerID, followingID)
	}

	return nil
}
