package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func NotificationWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	lastCheck := time.Time{}

	// Sending notifications in a separate goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				notifications, newLastCheck, err := db.FetchNewNotifications(userID, lastCheck)
				if err != nil {
					log.Println("Error fetching notifications:", err)
					continue
				}

				if len(notifications) > 0 {
					if err := conn.WriteJSON(notifications); err != nil {
						log.Println("Error writing notifications:", err)
						return
					}
					lastCheck = newLastCheck
				}
			}
		}
	}()

	// Reading and handling incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break // Exit the loop (and thus close the connection) on read error
		}

		var resp db.NotificationResponse
		if err = json.Unmarshal(msg, &resp); err != nil {
			log.Println("Unmarshal error:", err)
			continue // Continue reading next message despite the current error
		}

		if resp.Action == "notification_response" {
			if err := db.HandleNotificationResponse(userID, resp.NotificationID, resp.Accepted); err != nil {
				log.Println("Error handling notification response:", err)
			}
		}
	}
}

func ChatNotificationWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Sending unread message IDs in a separate goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-ticker.C:
				unreadMessageIDs, err := db.GetUnreadChatIDs(userID)
				if err != nil {
					log.Println("Error fetching unread message IDs:", err)
					continue
				}

				if err := conn.WriteJSON(unreadMessageIDs); err != nil {
					return
				}
			}
		}
	}()
}
