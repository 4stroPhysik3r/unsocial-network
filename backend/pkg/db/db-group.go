package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Group struct {
	GroupID          int    `json:"group_id,omitempty"`
	UserID           int    `json:"user_id"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	CreatorName      string `json:"creator_name"`
	CreatorFirstname string
	CreatorLastname  string
	Members          []int     `json:"members"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type GroupPost struct {
	PostID    int       `json:"post_id,omitempty"`
	UserID    int       `json:"user_id"`
	GroupID   int       `json:"group_id,omitempty"`
	Content   string    `json:"content"`
	PostImage *string   `json:"post_image,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type InviteRequest struct {
	SelectedMembers []int `json:"selectedMembers"`
}

func InsertGroup(group Group) (groupID int, err error) {
	statement, err := DB.Prepare(`INSERT INTO groups (user_id, title, content) VALUES (?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing insert statement for group: %v", err)
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(group.UserID, group.Title, group.Content)
	if err != nil {
		log.Printf("Error executing insert statement for group: %v", err)
		return 0, err
	}

	GroupID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return 0, err
	}

	var chatID int64
	chatResult, err := DB.Exec("INSERT INTO chats (group_id) VALUES (?)", GroupID)
	if err != nil {
		log.Printf("Error creating new chat: %v", err)
		return 0, err
	}

	chatID, err = chatResult.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving chat ID: %v", err)
		return 0, err
	}

	_, err = DB.Exec(`INSERT INTO chat_participants (chat_id, participant_id) 
        VALUES (?, ?)`, chatID, group.UserID)
	if err != nil {
		log.Printf("Error adding users to chat participants: %v", err)
		return 0, err
	}

	return int(GroupID), nil
}

func InsertGroupMember(groupID, memberID int, status string) error {
	// Prepare the SQL statement to insert a new group member with status
	stmt, err := DB.Prepare("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to insert the group member
	_, err = stmt.Exec(groupID, memberID, status)
	if err != nil {
		return err
	}

	return nil
}

func JoinGroup(userID, groupID int) error {
	var status string

	err := DB.QueryRow("SELECT status FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID).Scan(&status)
	switch {
	case err == sql.ErrNoRows:
		_, err := DB.Exec("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)", groupID, userID, "request")
		if err != nil {
			log.Printf("Error inserting join group request into database: %v", err)
			return err
		}

	case err != nil:
		log.Printf("Error querying group_members table: %v", err)
		return err

	default:
		if status != "request" {
			_, err := DB.Exec("UPDATE group_members SET status = 'request' WHERE group_id = ? AND user_id = ?", groupID, userID)
			if err != nil {
				log.Printf("Error updating join group request status in database: %v", err)
				return err
			}
		}
	}

	return nil
}

func LeaveGroup(userID, groupID int) error {
	// Update status in the group_members table to 'rejected'
	_, err := DB.Exec("UPDATE group_members SET status = 'rejected' WHERE group_id = ? AND user_id = ?", groupID, userID)
	if err != nil {
		log.Printf("Error updating group membership status: %v", err)
		return err
	}

	// leave the group chat
	err = removeUserFromGroupChat(groupID, userID)
	if err != nil {
		log.Printf("Error removing user from group chat: %v", err)
		return err
	}

	return nil
}

// GetMyGroups fetches all groups that the user is a member of
func GetMyGroups(userID int) ([]Group, error) {
	query := `
        SELECT g.group_id, g.user_id, g.title, g.content, u.firstname, u.lastname, g.created_at
        FROM groups g
        INNER JOIN users u ON g.user_id = u.user_id
        INNER JOIN group_members gm ON g.group_id = gm.group_id
        WHERE gm.user_id = ? AND gm.status = 'accepted'`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetMyGroups: failed to query groups: %v", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.GroupID, &group.UserID, &group.Title, &group.Content, &group.CreatorFirstname, &group.CreatorLastname, &group.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetMyGroups: failed to scan group row: %v", err)
		}
		group.CreatorName = fmt.Sprintf("%s %s", group.CreatorFirstname, group.CreatorLastname)
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetMyGroups: error iterating over group rows: %v", err)
	}

	return groups, nil
}

func GetAllGroups() ([]Group, error) {
	query := `
        SELECT g.group_id, g.user_id, g.title, g.content, u.firstname, u.lastname, g.created_at
        FROM groups g
        INNER JOIN users u ON g.user_id = u.user_id`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllGroups: failed to query groups: %v", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.GroupID, &group.UserID, &group.Title, &group.Content, &group.CreatorFirstname, &group.CreatorLastname, &group.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetAllGroups: failed to scan group row: %v", err)
		}
		group.CreatorName = fmt.Sprintf("%s %s", group.CreatorFirstname, group.CreatorLastname)
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllGroups: error iterating over group rows: %v", err)
	}

	return groups, nil
}

// Modify the GetGroupByID function to fetch group members
func GetGroupByID(groupID int) (Group, error) {
	query := `
		SELECT g.group_id, g.user_id, g.title, g.content, u.firstname, u.lastname, g.created_at, g.user_id,
		m.user_id
		FROM groups g
		INNER JOIN users u ON g.user_id = u.user_id
		LEFT JOIN group_members m ON g.group_id = m.group_id
		WHERE g.group_id = ? AND m.status IN ('accepted', 'request', 'invited')`

	var group Group
	rows, err := DB.Query(query, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Group{}, fmt.Errorf("GetGroupByID: group with ID %d not found", groupID)
		}
		return Group{}, fmt.Errorf("GetGroupByID: failed to fetch group: %v", err)
	}
	defer rows.Close()

	// Map to store member IDs temporarily
	memberIDsMap := make(map[int]bool)

	for rows.Next() {
		var memberID int
		err := rows.Scan(&group.GroupID, &group.UserID, &group.Title, &group.Content, &group.CreatorFirstname,
			&group.CreatorLastname, &group.CreatedAt, &group.UserID, &memberID)
		if err != nil {
			return Group{}, fmt.Errorf("GetGroupByID: failed to scan row: %v", err)
		}
		group.CreatorName = fmt.Sprintf("%s %s", group.CreatorFirstname, group.CreatorLastname)

		// Add member ID to the map if not already present
		memberIDsMap[memberID] = true
	}

	// Convert map keys to a slice of member IDs
	for memberID := range memberIDsMap {
		group.Members = append(group.Members, memberID)
	}

	return group, nil
}

func CheckMembership(userID, groupID int) (string, error) {
	var status string

	// Prepare the SQL query to select the status
	query := "SELECT status FROM group_members WHERE user_id = ? AND group_id = ? LIMIT 1"

	// Execute the query
	err := DB.QueryRow(query, userID, groupID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are found, the user is not a member of the group
			return "not_member", nil
		}
		// Return an error for any other database error
		return "", fmt.Errorf("CheckMembership: failed to check membership: %v", err)
	}

	// Return the membership status
	return status, nil
}

func InsertGroupPost(groupPost GroupPost) error {
	statement, err := DB.Prepare(`
        INSERT INTO posts (user_id, group_id, content, post_image)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		log.Printf("Error preparing insert statement: %v", err)
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(groupPost.UserID, groupPost.GroupID, groupPost.Content, groupPost.PostImage)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return err
	}

	fmt.Println("Group Post inserted to DB successfully! postID: ", postID)

	return nil
}

func handleGroupInvitation(tx *sql.Tx, groupID, userID int, status string) error {

	_, err := tx.Exec("UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?", status, groupID, userID)
	if err != nil {
		tx.Rollback() // Roll back in case of error
		log.Printf("Failed to update group membership status: %v", err)
		return err
	}

	// If the group invitation is accepted, add chat participant
	if status == "accepted" {
		return addUserToGroupChat(tx, groupID, userID)
	}

	return nil
}

func handleJoinGroup(tx *sql.Tx, groupID int, secondReferenceID sql.NullInt64, status string) error {

	// Update the group_members table to reflect the decision for the user who requested to join
	_, err := tx.Exec("UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ?", status, groupID, secondReferenceID.Int64)
	if err != nil {
		tx.Rollback() // Roll back in case of error
		log.Printf("Failed to update join group request status: %v", err)
		return err
	}

	// If the group invitation is accepted, add chat participant
	if status == "accepted" {
		return addUserToGroupChat(tx, groupID, int(secondReferenceID.Int64))
	}

	return nil
}

func addUserToGroupChat(tx *sql.Tx, groupID int, userID int) error {
	var chatID int64

	// Query to select the chat_id from the chats table where the group_id matches
	query := "SELECT chat_id FROM chats WHERE group_id = ? LIMIT 1"
	err := tx.QueryRow(query, groupID).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			// The case where no matching chat is found
			log.Printf("No chat found for group ID: %v", groupID)
			return err
		} else {
			// Handle other errors
			log.Printf("Error querying chat by group ID: %v", err)
			return err
		}
	}

	// Insert the user into the chat_participants table
	_, err = tx.Exec("INSERT INTO chat_participants (chat_id, participant_id) VALUES (?, ?)", chatID, userID)
	if err != nil {
		log.Printf("Error adding user to chat participants: %v", err)
		return err
	}

	return nil
}

func removeUserFromGroupChat(groupID int, userID int) error {
	var chatID int64

	// Query to select the chat_id from the chats table where the group_id matches
	query := "SELECT chat_id FROM chats WHERE group_id = ? LIMIT 1"
	err := DB.QueryRow(query, groupID).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			// The case where no matching chat is found
			log.Printf("No chat found for group ID: %v", groupID)
			return err
		} else {
			// Handle other errors
			log.Printf("Error querying chat by group ID: %v", err)
			return err
		}
	}

	// Delete the user from the chat_participants table
	_, err = DB.Exec("DELETE FROM chat_participants WHERE chat_id = ? AND participant_id = ?", chatID, userID)
	if err != nil {
		log.Printf("Error removing user from chat participants: %v", err)
		return err
	}

	return nil
}
