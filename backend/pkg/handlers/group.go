package handlers

import (
	"backend/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
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

	var groupData db.Group
	groupData.UserID = userID

	err := json.NewDecoder(r.Body).Decode(&groupData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	groupID, err := db.InsertGroup(groupData)
	if err != nil {
		http.Error(w, "Failed to insert group data", http.StatusInternalServerError)
		log.Printf("Failed to insert group data: %v", err)
		return
	}

	// Insert the creator as a member with status 'accepted'
	err = db.InsertGroupMember(groupID, userID, "accepted")
	if err != nil {
		http.Error(w, "Failed to insert creator data", http.StatusInternalServerError)
		log.Printf("Failed to insert creator data: %v", err)
		return
	}

	// Insert other selected members with status 'invited'
	for _, memberID := range groupData.Members {

		err = db.InsertGroupMember(groupID, memberID, "invited")
		if err != nil {
			http.Error(w, "Failed to insert member data", http.StatusInternalServerError)
			log.Printf("Failed to insert member data: %v", err)
			return
		}

		err = db.CreateGroupInvitationNotification(userID, memberID, groupID)
		if err != nil {
			log.Printf("Error creating group invitation notification: %v", err)
			return
		}
	}

	response := map[string]string{"message": "Group created successfully, and members are sent a notification to join to group"}
	json.NewEncoder(w).Encode(response)
}

func InviteToGroupHandler(w http.ResponseWriter, r *http.Request) {
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

	groupIDStr := r.URL.Path[len("/api/invite-to-group/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	// Decode the request body into an InviteRequest struct
	var inviteReq db.InviteRequest
	err = json.NewDecoder(r.Body).Decode(&inviteReq)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Now you can access the selectedMembers array from inviteReq
	selectedMembers := inviteReq.SelectedMembers

	// Insert other selected members with status 'invited'
	for _, memberID := range selectedMembers {

		err = db.InsertGroupMember(groupID, memberID, "invited")
		if err != nil {
			http.Error(w, "Failed to insert member data", http.StatusInternalServerError)
			log.Printf("Failed to insert member data: %v", err)
			return
		}

		err = db.CreateGroupInvitationNotification(userID, memberID, groupID)
		if err != nil {
			log.Printf("Error creating group invitation notification: %v", err)
			return
		}
	}

	// Send success response
	response := map[string]string{"message": "Users invited successfully"}
	json.NewEncoder(w).Encode(response)

}

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse the dynamic part of the URL
	groupIDStr := r.URL.Path[len("/api/join-group/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	err = db.JoinGroup(userID, groupID)
	if err != nil {
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		log.Printf("Error joining group: %v", err)
		return
	}

	err = db.CreateJoinGroupRequestNotification(userID, groupID)
	if err != nil {
		log.Printf("Error creating join group invitation notification: %v", err)
		return
	}

	response := map[string]string{"message": "Successfully requested to join the group"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse the dynamic part of the URL
	groupIDStr := r.URL.Path[len("/api/leave-group/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	err = db.LeaveGroup(userID, groupID)
	if err != nil {
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		log.Printf("Error leaving group: %v", err)
		return
	}

	response := map[string]string{"message": "Successfully requested to leave the group"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetMyGroupsHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	groups, err := db.GetMyGroups(userID)
	if err != nil {
		http.Error(w, "GetMyGroups: Failed to fetch groups", http.StatusInternalServerError)
		log.Printf("GetMyGroups: Failed to fetch groups %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groups, err := db.GetAllGroups()
	if err != nil {
		http.Error(w, "GetAllGroups: Failed to fetch groups", http.StatusInternalServerError)
		log.Printf("GetAllGroups: Failed to fetch groups %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetGroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	groupIDStr := r.URL.Path[len("/api/get-group-posts/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	// Call your DB function to fetch the post by its ID
	post, err := db.GetPostsByGroupID(groupID)
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

func GetGroupInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the dynamic part of the URL
	groupIDStr := r.URL.Path[len("/api/get-group-info/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	// Call your DB function to fetch the post by its ID
	info, err := db.GetGroupByID(groupID)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	// Serialize the post to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CheckMembershipHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromSession(r)
	if userID == 0 {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse the dynamic part of the URL
	groupIDStr := r.URL.Path[len("/api/check-membership/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	response, err := db.CheckMembership(userID, groupID)
	if err != nil {
		http.Error(w, "GetMyGroups: Failed to fetch groups", http.StatusInternalServerError)
		log.Printf("GetMyGroups: Failed to fetch groups %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
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
	groupIDStr := r.URL.Path[len("/api/create-group-post/"):]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	var groupPostData db.GroupPost
	err = json.NewDecoder(r.Body).Decode(&groupPostData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	groupPostData.UserID = userID
	groupPostData.GroupID = groupID

	if groupPostData.PostImage != nil && *groupPostData.PostImage != "" {

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Post: Error getting current working directory: %v", err)
		}

		dirPath := cwd + "/uploads/post-image"

		fileName, err := saveBase64File(*groupPostData.PostImage, dirPath)
		if err != nil {
			http.Error(w, "Post: Failed to save post image", http.StatusInternalServerError)
			return
		}
		webPath := "http://localhost:8000/uploads/post-image/" + fileName
		groupPostData.PostImage = &webPath
	}

	err = db.InsertGroupPost(groupPostData)
	if err != nil {
		http.Error(w, "Failed to insert post data", http.StatusInternalServerError)
		log.Printf("Failed to insert post data: %v", err)
		return
	}

	response := map[string]string{"message": "Post created successfully"}
	json.NewEncoder(w).Encode(response)
}
