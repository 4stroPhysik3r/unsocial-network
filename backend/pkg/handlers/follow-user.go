package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetFollowStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is authenticated
	followerID := userIDFromSession(r)
	if followerID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		log.Println("UserProfile: User not authenticated")
		return
	}

	// Extract the user ID of the user to be followed from the URL
	followingIDStr := r.URL.Path[len("/api/get-follow-status/"):]
	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Update the user profile in the database
	var status string
	status, err = db.GetFollowStatus(followerID, followingID)
	if err != nil {
		http.Error(w, "Failed to update user status", http.StatusInternalServerError)
		log.Printf("FollowStatus: Failed to update user status: %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	followerID := userIDFromSession(r)
	if followerID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	followingIDStr := r.URL.Path[len("/api/follow-user/"):]
	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	err = db.FollowUserRequest(followerID, followingID)
	if err != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		log.Printf("Error following user: %v", err)
		return
	}

	

	response := map[string]string{"message": "Follow request sent successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	followerID := userIDFromSession(r) 
	if followerID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	followingIDStr := r.URL.Path[len("/api/unfollow-user/"):]
	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	err = db.UnfollowUser(followerID, followingID)
	if err != nil {
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		log.Printf("Error unfollowing user: %v", err)
		return
	}

	response := map[string]string{"message": "Unfollowed user successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetFollowingListForPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	followers, err := db.GetFollowing(userID)
	if err != nil {
		http.Error(w, "Failed to fetch followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	userIDStr := r.URL.Path[len("/api/following/"):]

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	following, err := db.GetFollowing(userID)
	if err != nil {
		http.Error(w, "Failed to fetch following", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func GetFollower(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	userIDStr := r.URL.Path[len("/api/follower/"):]

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	followers, err := db.GetFollowers(userID)
	if err != nil {
		http.Error(w, "Failed to fetch followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
