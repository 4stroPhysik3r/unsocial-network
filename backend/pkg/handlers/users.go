package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check the HTTP method
	if r.Method == http.MethodGet {
		getUserData(w, r)
	} else if r.Method == http.MethodPost || r.Method == http.MethodPut {
		updateUserProfile(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if the user is authenticated
	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		log.Println("UserData: User not authenticated")
		return
	}

	// Retrieve user data from the database
	userData, err := db.UserDataFromID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		log.Printf("UserData: Failed to fetch user data: %v", err)
		return
	}

	usersJSON, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		log.Printf("UserData: Failed to encode user data: %v", err)
		return
	}

	_, err = w.Write(usersJSON)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check the HTTP method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	getUsersData(w)
}

func getUsersData(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve all user data from the database
	allUsersData, err := db.GetUsersFromDB()
	if err != nil {
		http.Error(w, "Failed to fetch all users data", http.StatusInternalServerError)
		log.Printf("UserData: Failed to fetch all user data: %v", err)
		return
	}

	usersJSON, err := json.Marshal(allUsersData)
	if err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		log.Printf("UserData: Failed to encode users data: %v", err)
		return
	}

	_, err = w.Write(usersJSON)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func UserIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check the HTTP method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	userIDStr := r.URL.Path[len("/api/userid/"):]

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Retrieve user data from the database
	userData, err := db.UserDataFromID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		log.Printf("UserData: Failed to fetch user data: %v", err)
		return
	}

	// Encode user data to JSON
	err = json.NewEncoder(w).Encode(userData)
	if err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		log.Printf("UserData: Failed to encode user data: %v", err)
		return
	}
}
