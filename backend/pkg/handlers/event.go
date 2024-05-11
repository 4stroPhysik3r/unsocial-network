package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var eventData db.Event
	eventData.UserID = userID

	err := json.NewDecoder(r.Body).Decode(&eventData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	err = db.InsertEvent(eventData)
	if err != nil {
		http.Error(w, "Failed to insert group data", http.StatusInternalServerError)
		log.Printf("Failed to insert group data: %v", err)
		return
	}

	groupIDStr := r.URL.Path[len("/api/create-event/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	err = db.CreateEventNotification(groupID)
	if err != nil {
		log.Printf("Error creating new event notification: %v", err)
		return
	}

	response := map[string]string{"message": "Event created successfully"}
	json.NewEncoder(w).Encode(response)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	groupIDStr := r.URL.Path[len("/api/get-events/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	events, err := db.GetEvents(groupID)
	if err != nil {
		http.Error(w, "GetEvents: Failed to fetch events", http.StatusInternalServerError)
		log.Printf("GetEvents: Failed to fetch events %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateAttendeesStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Check if the user is authenticated
	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		log.Println("UserProfile: User not authenticated")
		return
	}

	// Parse the dynamic part of the URL
	eventIDStr := r.URL.Path[len("/api/update-attendees-status/"):]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated status
	var updatedStatus db.AttendeesStatus
	err = json.NewDecoder(r.Body).Decode(&updatedStatus)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("UpdateAttendeesStatus: Failed to decode request body: %v", err)
		return
	}

	// Update the attendees' status for the event in the database
	err = db.UpdateAttendeesStatus(userID, eventID, updatedStatus.AttendeesStatus)
	if err != nil {
		http.Error(w, "Failed to update attendees' status", http.StatusInternalServerError)
		log.Printf("UpdateAttendeesStatus: Failed to update attendees' status: %v", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User profile updated successfully")
}

func GetAttendeesStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		log.Println("UserProfile: User not authenticated")
		return
	}

	eventIDStr := r.URL.Path[len("/api/get-attendees-status/"):]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid eventID", http.StatusBadRequest)
		return
	}

	status, err := db.GetAttendeesStatus(userID, eventID)
	if err != nil {
		http.Error(w, "GetStatus: Failed to fetch status", http.StatusInternalServerError)
		log.Printf("GetStatus: Failed to fetch status %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
