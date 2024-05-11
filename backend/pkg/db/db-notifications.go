package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Notification struct {
	NotificationID int       `json:"notification_id"`
	UserID         int       `json:"user_id"`
	Type           string    `json:"type"`
	Message        string    `json:"message"`
	Status         string    `json:"status"`
	ReferenceID    int       `json:"reference_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type NotificationResponse struct {
	Action         string `json:"action"`
	NotificationID int    `json:"notification_id"`
	Accepted       bool   `json:"accepted"`
}

func FetchNewNotifications(userID int, lastCheck time.Time) ([]Notification, time.Time, error) {

	var notifications []Notification

	// Adjust the query based on the value of lastCheck
	query := ""
	if lastCheck.IsZero() {
		query = `SELECT notification_id, user_id, type, message, status, reference_id, created_at FROM notifications WHERE user_id = ? AND status = 'unread' ORDER BY created_at ASC`
	} else {
		// Fetch notifications newer than lastCheck
		query = `SELECT notification_id, user_id, type, message, status, reference_id, created_at FROM notifications WHERE user_id = ? AND created_at > ? ORDER BY created_at ASC`
	}

	// Prepare and execute the query based on whether lastCheck is zero
	var rows *sql.Rows
	var err error
	if lastCheck.IsZero() {
		rows, err = DB.Query(query, userID)
	} else {
		rows, err = DB.Query(query, userID, lastCheck)
	}

	if err != nil {
		return nil, lastCheck, err
	}
	defer rows.Close()

	var maxTime time.Time
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.NotificationID, &n.UserID, &n.Type, &n.Message, &n.Status, &n.ReferenceID, &n.CreatedAt); err != nil {
			log.Printf("Error scanning notification: %v", err)
			continue
		}
		notifications = append(notifications, n)
		if n.CreatedAt.After(maxTime) {
			maxTime = n.CreatedAt
		}
	}

	if err = rows.Err(); err != nil {
		return nil, lastCheck, err
	}

	// If no new notifications were fetched, return the original lastCheck
	if maxTime.IsZero() {
		maxTime = lastCheck
	}

	return notifications, maxTime, nil
}

func HandleNotificationResponse(userID, notificationID int, accepted bool) error {

	// Determine the new status based on the 'accepted' flag
	status := "rejected"
	if accepted {
		status = "accepted"
	}

	// Begin a transaction
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	// Update the notification's status
	_, err = tx.Exec("UPDATE notifications SET status = ? WHERE notification_id = ? AND user_id = ?", status, notificationID, userID)
	if err != nil {
		tx.Rollback() // Roll back in case of error
		log.Printf("Failed to update notification status: %v", err)
		return err
	}

	// Fetch the type and reference ID of the notification
	var notifType string
	var referenceID int
	var secondReferenceID sql.NullInt64 // Use sql.NullInt64 for nullable integer columns
	err = tx.QueryRow("SELECT type, reference_id, second_reference_id FROM notifications WHERE notification_id = ?", notificationID).Scan(&notifType, &referenceID, &secondReferenceID)
	if err != nil {
		tx.Rollback() // Roll back in case of error
		log.Printf("Failed to fetch notification type and reference ID: %v", err)
		return err
	}

	// Handle the follow request response
	if notifType == "follow_request" {
		err = handleFollowRequestResponse(tx, referenceID, userID, status)
		if err != nil {
			tx.Rollback() // Roll back in case of error
			return err
		}
	}

	// Handle group_invitation notifications
	if notifType == "group_invitation" {
		err = handleGroupInvitation(tx, referenceID, userID, status)
		if err != nil {
			tx.Rollback() // Roll back in case of error
			return err
		}
	}


	// Handle join_group_request notifications
	if notifType == "join_group_request" && secondReferenceID.Valid {
		err = handleJoinGroup(tx, referenceID, secondReferenceID, status)
		if err != nil {
			tx.Rollback() // Roll back in case of error
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func CreateFollowRequestNotification(followerID, followingID int) error {

	// First, retrieve the full name of the following
	var fullName string
	err := DB.QueryRow("SELECT firstname || ' ' || lastname FROM users WHERE user_id = ?", followerID).Scan(&fullName)
	if err != nil {
		log.Printf("Error querying database for user's full name: %v", err)
		return err
	}

	// Then, create the notification with the full name included in the message
	message := fmt.Sprintf("%s has sent you a follow request.", fullName)
	_, err = DB.Exec(`INSERT INTO notifications (user_id, type, message, reference_id)
                      VALUES (?, 'follow_request', ?, ?)`,
		followingID, message, followerID)
	if err != nil {
		log.Printf("Error inserting follow request notification: %v", err)
		return err
	}

	return nil
}

func CreateGroupInvitationNotification(invitingUserID, invitedUserID, groupID int) error {
	// First, retrieve the title of the group
	var groupTitle string
	err := DB.QueryRow("SELECT title FROM groups WHERE group_id = ?", groupID).Scan(&groupTitle)
	if err != nil {
		log.Printf("Error querying database for group title: %v", err)
		return err
	}

	// Then, retrieve the full name of the inviting user
	var fullName string
	err = DB.QueryRow("SELECT firstname || ' ' || lastname FROM users WHERE user_id = ?", invitingUserID).Scan(&fullName)
	if err != nil {
		log.Printf("Error querying database for user's full name: %v", err)
		return err
	}

	// Construct the invitation message
	message := fmt.Sprintf("%s has invited you to join the group '%s'.", fullName, groupTitle)

	// Insert the invitation notification into the notifications table
	_, err = DB.Exec(`INSERT INTO notifications (user_id, type, message, reference_id)
                      VALUES (?, 'group_invitation', ?, ?)`,
		invitedUserID, message, groupID)
	if err != nil {
		log.Printf("Error inserting group invitation notification: %v", err)
		return err
	}

	return nil
}

func CreateJoinGroupRequestNotification(requestingUserID, groupID int) error {

	var groupTitle string
	var creatorUserID int
	err := DB.QueryRow("SELECT title, user_id FROM groups WHERE group_id = ?", groupID).Scan(&groupTitle, &creatorUserID)
	if err != nil {
		log.Printf("Error querying database for group title and creator ID: %v", err)
		return err
	}

	var fullName string
	err = DB.QueryRow("SELECT firstname || ' ' || lastname FROM users WHERE user_id = ?", requestingUserID).Scan(&fullName)
	if err != nil {
		log.Printf("Error querying database for requesting user's full name: %v", err)
		return err
	}

	message := fmt.Sprintf("%s has requested to join your group '%s'.", fullName, groupTitle)

	// Including both reference_id (for the group) and second_reference_id (for the requesting user)
	_, err = DB.Exec(`INSERT INTO notifications (user_id, type, message, reference_id, second_reference_id)
                      VALUES (?, 'join_group_request', ?, ?, ?)`,
		creatorUserID, message, groupID, requestingUserID)
	if err != nil {
		log.Printf("Error inserting join group request notification: %v", err)
		return err
	}

	return nil
}

func CreateEventNotification(groupID int) error {
	var groupTitle string
	err := DB.QueryRow("SELECT title FROM groups WHERE group_id = ?", groupID).Scan(&groupTitle)
	if err != nil {
		log.Printf("Error querying database for group title and creator ID: %v", err)
		return err
	}

	message := fmt.Sprintf("There is a new event in the group '%s'.", groupTitle)

	memberIDs, err := GetAcceptedGroupMembers(groupID, DB)
	if err != nil {
		log.Printf("Error fetching accepted group members: %v", err)
		return err
	}

	for _, memberID := range memberIDs {
		// Including both reference_id (for the group) and second_reference_id (for the requesting user)
		_, err = DB.Exec(`INSERT INTO notifications (user_id, type, message, reference_id)
			VALUES (?, 'new_event', ?, ?)`,
			memberID, message, groupID)
		if err != nil {
			log.Printf("Error inserting join group request notification: %v", err)
			return err
		}

		fmt.Printf("Notification created for user %d: %s\n", memberID, message)
	}

	return nil
}

func GetAcceptedGroupMembers(groupID int, DB *sql.DB) ([]int, error) {
	var memberIDs []int
	rows, err := DB.Query("SELECT user_id FROM group_members WHERE group_id = ? AND status = 'accepted'", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memberIDs, nil
}
