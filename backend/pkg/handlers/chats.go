package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Map to keep track of connections by chat room
var chatRooms = make(map[string][]*websocket.Conn)

func GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	chatInfo, err := db.GetChatParticipantNames(userID)
	if err != nil {
		log.Fatal("Error getting chat participant names:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chatInfo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var chatReq db.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatReq); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	chatID := chatReq.ChatID
	if chatID <= 0 {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query("SELECT message_id, chat_id, sender_id, content, emoji, created_at FROM messages WHERE chat_id = ?", chatID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []db.ChatMessage
	for rows.Next() {
		var msg db.ChatMessage
		if err := rows.Scan(&msg.MessageID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.Emoji, &msg.CreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := db.DeleteUnreadMessagesForUser(userID, chatID); err != nil {
		log.Printf("Error deleting unread messages for user %d in chat %d: %v", userID, chatID, err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func ChatWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse the message to get the chat ID
		var chatMessage db.ChatMessage
		err = json.Unmarshal(message, &chatMessage)
		if err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		strChatID := strconv.Itoa(chatMessage.ChatID)

		if _, ok := chatRooms[strChatID]; !ok {
			chatRooms[strChatID] = []*websocket.Conn{}
		}

		if !containsConnection(chatRooms[strChatID], conn) {
			chatRooms[strChatID] = append(chatRooms[strChatID], conn)
		}

		if !(chatMessage.Content == "" || chatMessage.SenderID == 0) {
			if err := db.InsertChatMessage(chatMessage); err != nil {
				log.Printf("Error inserting chat message into the database: %v", err)
				continue
			}

			// Add the message to the unread_messages table for each participant in the chat
			participants, err := db.GetChatParticipants(chatMessage.ChatID)
			if err != nil {
				log.Printf("Error getting chat participants: %v", err)
				continue
			}

			// Delete unread messages for all participants in the chat
			if err := db.DeleteUnreadMessages(chatMessage.ChatID); err != nil {
				log.Printf("Error deleting unread messages for chat %d: %v", chatMessage.ChatID, err)
				continue
			}

			for _, participant := range participants {
				if participant.UserID != chatMessage.SenderID {
					timestamp := time.Now() // Generate the current timestamp
					if err := db.AddUnreadMessage(participant.UserID, chatMessage.ChatID, timestamp); err != nil {
						log.Printf("Error adding unread message for participant %d: %v", participant.UserID, err)
						continue
					}
				}
			}

			// Broadcast message to all participants in the chat room
			for _, participant := range chatRooms[strChatID] {
				// Now, we send the message to every participant, including the sender
				if err := participant.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// containsConnection checks if the connection is already in the list of connections
func containsConnection(connections []*websocket.Conn, conn *websocket.Conn) bool {
	for _, c := range connections {
		if c == conn {
			return true
		}
	}
	return false
}

func GetUnreadMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get unread chat IDs for the user
	chatIDs, err := db.GetUnreadChatIDs(userID)
	if err != nil {
		log.Printf("Error getting unread chat IDs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode chat IDs as JSON and send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chatIDs); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
