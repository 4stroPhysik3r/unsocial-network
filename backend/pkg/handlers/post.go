package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var postData db.Post
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	postData.UserID = userID

	if postData.PostImage != nil && *postData.PostImage != "" {

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Post: Error getting current working directory: %v", err)
		}

		dirPath := cwd + "/uploads/post-image"

		fileName, err := saveBase64File(*postData.PostImage, dirPath)
		if err != nil {
			http.Error(w, "Post: Failed to save post image", http.StatusInternalServerError)
			return
		}
		webPath := "http://localhost:8000/uploads/post-image/" + fileName
		postData.PostImage = &webPath
	}

	postID, err := db.InsertPost(postData)
	if err != nil {
		http.Error(w, "Failed to insert post data", http.StatusInternalServerError)
		log.Printf("Failed to insert post data: %v", err)
		return
	}

	if len(postData.ViewerIDs) > 0 {
		err = db.InsertPostViewers(postID, postData.ViewerIDs)
		if err != nil {
			http.Error(w, "Failed to insert post viewers", http.StatusInternalServerError)
			log.Printf("Failed to insert post viewers: %v", err)
			return
		}
	}

	response := map[string]string{"message": "Post created successfully"}
	json.NewEncoder(w).Encode(response)
}

func GetPostsHandlerForFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	posts, err := db.GetPostsForFeed(userID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPostsFromSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	posts, err := db.GetPostsForProfile(userID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPostsFromIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	loggedID := userIDFromSession(r)
	if loggedID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse the dynamic part of the URL
	userIDStr := r.URL.Path[len("/api/postsFromID/"):]

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	posts, err := db.GetPostsForUser(userID, loggedID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPostFromPostID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	postIDStr := r.URL.Path[len("/api/get-post/"):]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	// Call your DB function to fetch the post by its ID
	post, err := db.GetPostByPostID(postID)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	// Serialize the post to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ViewerStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	postIDStr := r.URL.Path[len("/api/viewer-status/"):]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Call your DB function to fetch the post by its ID
	response, err := db.GetViewerStatus(postID, userID)
	if err != nil {
		http.Error(w, "Failed to fetch viewer status for post", http.StatusInternalServerError)
		return
	}

	// Serialize the post to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
