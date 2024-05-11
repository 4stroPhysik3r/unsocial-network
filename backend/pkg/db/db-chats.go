package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ChatInfo struct {
	ChatID   int    `json:"chat_id"`
	FullName string `json:"full_name"`
}

type ChatMessage struct {
	MessageID int            `json:"message_id,omitempty"`
	ChatID    int            `json:"chat_id"`
	SenderID  int            `json:"sender_id"`
	Content   string         `json:"content"`
	Emoji     sql.NullString `json:"emoji,omitempty"`
	CreatedAt string         `json:"created_at,omitempty"`
}

type ChatRequest struct {
	ChatID int `json:"chat_id"`
}

type Participant struct {
	UserID int
}

func createChatForUsers(tx *sql.Tx, user1ID, user2ID int) error {

	// Check if a chat already exists between the two users
	var existingChatID int
	query := `
	SELECT c.chat_id 
	FROM chats c
	JOIN chat_participants cp1 ON c.chat_id = cp1.chat_id AND cp1.participant_id = ?
	JOIN chat_participants cp2 ON c.chat_id = cp2.chat_id AND cp2.participant_id = ?
	WHERE c.group_id IS NULL
	LIMIT 1;
    `
	err := tx.QueryRow(query, user1ID, user2ID).Scan(&existingChatID)
	if err == nil {
		// A chat already exists, so no need to create a new one
		return nil
	} else if err != sql.ErrNoRows {
		// An actual error occurred
		log.Printf("Error checking for existing chat: %v", err)
		return err
	}

	// No existing chat found, proceed to create a new chat
	var chatID int64
	result, err := tx.Exec("INSERT INTO chats DEFAULT VALUES")
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
	_, err = tx.Exec(`
        INSERT INTO chat_participants (chat_id, participant_id) 
        VALUES (?, ?), (?, ?)
    `, chatID, user1ID, chatID, user2ID)
	if err != nil {
		log.Printf("Error adding users to chat participants: %v", err)
		return err
	}

	return nil
}

func GetChat(user1ID, user2ID int) (int, error) {
	var chatID int

	// Check if a chat already exists between the two users
	err := DB.QueryRow(`
        SELECT c.chat_id
        FROM chats c
        JOIN chat_participants cp1 ON c.chat_id = cp1.chat_id AND cp1.participant_id = ?
        JOIN chat_participants cp2 ON c.chat_id = cp2.chat_id AND cp2.participant_id = ?
    `, user1ID, user2ID).Scan(&chatID)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying for existing chat: %v", err)
		return 0, err
	}
	return chatID, nil

}

func GetChatParticipantNames(excludeUserID int) ([]ChatInfo, error) {
	var chatNames []ChatInfo

	chatIDs, err := GetUserChatIDs(excludeUserID)
	if err != nil {
		log.Fatal("Error getting user chat IDs:", err)
	}
	chatNames, err = GetChatInfo(chatIDs, excludeUserID)
	if err != nil {
		log.Fatal("Error getting GetChatInfo IDs:", err)
	}

	return chatNames, nil

}

func GetUserChatIDs(userID int) ([]int, error) {
	var chatIDs []int

	// Use DISTINCT to ensure each chat_id is unique in the result set
	query := `
        SELECT DISTINCT chat_id
        FROM chat_participants
        WHERE participant_id = ?
    `
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chat IDs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chatID int
		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("error scanning chat ID: %v", err)
		}
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat IDs: %v", err)
	}

	return chatIDs, nil
}

func GetChatInfo(chatIds []int, userId int) ([]ChatInfo, error) {
	var chatInfos []ChatInfo

	for _, chatId := range chatIds {
		var chatInfo ChatInfo
		chatInfo.ChatID = chatId

		// Check if the chat is a group chat
		var groupId sql.NullInt64
		err := DB.QueryRow("SELECT group_id FROM chats WHERE chat_id = ?", chatId).Scan(&groupId)
		if err != nil {
			return nil, err
		}

		if groupId.Valid { // It's a group chat
			err := DB.QueryRow("SELECT title FROM groups WHERE group_id = ?", groupId.Int64).Scan(&chatInfo.FullName)
			if err != nil {
				return nil, err
			}
		} else { // It's a one-on-one chat
			var otherParticipantId int
			err := DB.QueryRow(`
				SELECT participant_id FROM chat_participants 
				WHERE chat_id = ? AND participant_id != ?`, chatId, userId).Scan(&otherParticipantId)
			if err != nil {
				return nil, err
			}

			var firstName, lastName string
			err = DB.QueryRow("SELECT firstname, lastname FROM users WHERE user_id = ?", otherParticipantId).Scan(&firstName, &lastName)
			if err != nil {
				return nil, err
			}
			chatInfo.FullName = firstName + " " + lastName
		}

		chatInfos = append(chatInfos, chatInfo)
	}

	return chatInfos, nil
}

func InsertChatMessage(message ChatMessage) error {
	statement, err := DB.Prepare(`INSERT INTO messages (chat_id, sender_id, content, emoji, created_at) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing insert statement: %v", err)
		return err
	}
	defer statement.Close()

	// Assuming message.CreatedAt is a string in the correct TIMESTAMP format for your database
	_, err = statement.Exec(message.ChatID, message.SenderID, message.Content, message.Emoji, message.CreatedAt)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		return err
	}

	return nil
}

// GetChatParticipants retrieves the list of participants in a chat.
func GetChatParticipants(chatID int) ([]Participant, error) {
	var participants []Participant

	rows, err := DB.Query("SELECT participant_id FROM chat_participants WHERE chat_id = ?", chatID)
	if err != nil {
		log.Printf("Error querying participants for chat %d: %v", chatID, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var participantID int
		if err := rows.Scan(&participantID); err != nil {
			log.Printf("Error scanning participant ID: %v", err)
			continue
		}
		participants = append(participants, Participant{UserID: participantID})
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating participant rows: %v", err)
		return nil, err
	}

	return participants, nil
}

// AddUnreadMessage adds an unread message for a participant in a specific chat.
func AddUnreadMessage(userID, chatID int, timestamp time.Time) error {
	_, err := DB.Exec("INSERT INTO unread_messages (user_id, chat_id, timestamp) VALUES (?, ?, ?)", userID, chatID, timestamp)
	if err != nil {
		log.Printf("Error adding unread message for user %d in chat %d: %v", userID, chatID, err)
		return err
	}
	return nil
}

func DeleteUnreadMessages(chatID int) error {
	_, err := DB.Exec("DELETE FROM unread_messages WHERE chat_id = ?", chatID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUnreadMessagesForUser(userID, chatID int) error {
	_, err := DB.Exec("DELETE FROM unread_messages WHERE user_id = ? AND chat_id = ?", userID, chatID)
	if err != nil {
		return err
	}
	return nil
}

func GetUnreadChatIDs(userID int) ([]int, error) {
	// Query the database to get all chat IDs with unread messages for the user
	rows, err := DB.Query("SELECT chat_id FROM unread_messages WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error querying unread_messages table: %v", err)
	}
	defer rows.Close()

	// Store chat IDs with unread messages
	var chatIDs []int
	for rows.Next() {
		var chatID int
		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		chatIDs = append(chatIDs, chatID)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return chatIDs, nil
}
