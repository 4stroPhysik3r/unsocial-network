package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData db.User
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userData.Avatar != "" {
		//getting current directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Register: Error getting current working directory: %v", err)
		}

		dirPath := cwd + "/uploads/avatars"

		fileName, err := saveBase64File(userData.Avatar, dirPath)
		if err != nil {
			log.Printf("Register: Error saving avatar: %v\n", err)
			http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
			return
		}

		webPath := "http://localhost:8000/uploads/avatars/" + fileName
		userData.Avatar = webPath
	} else if userData.Avatar == "" {
		webPath := "http://localhost:8000/uploads/avatars/default-avatar-profile.jpg"
		userData.Avatar = webPath
	}

	// Check if the email is already taken
	err = db.RegisterUser(userData)
	if err != nil {
		log.Printf("Register: Email already taken: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest) // Return an error response
		return
	}

	// Send a response back to the client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User registered successfully")
}
