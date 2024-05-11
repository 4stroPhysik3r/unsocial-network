package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if the user is authenticated
	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		log.Println("UserProfile: User not authenticated")
		return
	}

	// Decode the request body to get the updated profilePublic value
	var updatedProfilePublic bool
	err := json.NewDecoder(r.Body).Decode(&updatedProfilePublic)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("UserProfile: Failed to decode request body: %v", err)
		return
	}

	// Update the user profile in the database
	err = db.UpdateUserProfile(userID, updatedProfilePublic)
	if err != nil {
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		log.Printf("UserProfile: Failed to update user profile: %v", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User profile updated successfully")
}
