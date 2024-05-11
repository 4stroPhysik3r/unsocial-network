package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	// Parse the dynamic part of the URL
	postIDStr := r.URL.Path[len("/api/add-comment/"):]

	// Convert userID from string to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	var commentData db.Comment
	err = json.NewDecoder(r.Body).Decode(&commentData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	commentData.UserID = userID
	commentData.PostID = postID

	if commentData.CommentImage != nil && *commentData.CommentImage != "" {

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Post: Error getting current working directory: %v", err)
		}

		dirPath := cwd + "/uploads/comment-image"

		fileName, err := saveBase64File(*commentData.CommentImage, dirPath)
		if err != nil {
			http.Error(w, "Failed to save comment image", http.StatusInternalServerError)
			return
		}
		webPath := "http://localhost:8000/uploads/comment-image/" + fileName
		commentData.CommentImage = &webPath
	}

	err = db.InsertComment(commentData)
	if err != nil {
		http.Error(w, "Failed to insert comment data", http.StatusInternalServerError)
		log.Printf("Failed to insert comment data: %v", err)
		return
	}

	response := map[string]string{"message": "Comment created successfully"}
	json.NewEncoder(w).Encode(response)
}

func GetCommentsFromPostIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	postIDStr := r.URL.Path[len("/api/get-comments-for-post/"):]

	// Convert userID from string to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	comments, err := db.GetCommentsForPost(postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}