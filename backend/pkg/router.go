package pkg

import (
	"backend/pkg/db"
	"backend/pkg/handlers"
	"net/http"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/get-posts-feed", handlers.GetPostsHandlerForFeed)               // fetches posts to display on the feed
	mux.HandleFunc("/api/register", handlers.RegisterHandler)                            // gets data from form to register a new user
	mux.HandleFunc("/api/login", handlers.LoginHandler)                                  // gets data from form to check credentials, and if matching creates a session
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)                                // deletes cookie and session record from db
	mux.HandleFunc("/api/auth/status", handlers.AuthStatusHandler)                       // checks if the user is authenticated
	mux.HandleFunc("/api/create-post/", handlers.CreatePostHandler)                      // gets data from form to create a new post and save in db
	mux.HandleFunc("/api/add-comment/", handlers.CreateCommentHandler)                   // gets data from form to create a new comment and save in db
	mux.HandleFunc("/api/get-comments-for-post/", handlers.GetCommentsFromPostIDHandler) // fetches comments for a post
	mux.HandleFunc("/api/posts/", handlers.GetPostsFromSessionHandler)                   // fetches posts to display on the My profile based on ID from session
	mux.HandleFunc("/api/postsFromID/", handlers.GetPostsFromIDHandler)                  // fetches posts to display on the Other users profile based on ID from front end sent in URL
	mux.HandleFunc("/api/userid/", handlers.UserIDHandler)                               // displays user (profile) page based on id
	mux.HandleFunc("/api/my-profile", handlers.UserHandler)                              // display my profile information, and update privacy status
	mux.HandleFunc("/api/get-users", handlers.UsersHandler)                              // fetches data about all users
	mux.HandleFunc("/api/get-post/", handlers.GetPostFromPostID)                         // fetches single post based on postID from frontend

	mux.HandleFunc("/api/get-follow-status/", handlers.GetFollowStatusHandler)
	mux.HandleFunc("/api/follow-user/", handlers.FollowUserHandler)
	mux.HandleFunc("/api/unfollow-user/", handlers.UnfollowUserHandler)
	mux.HandleFunc("/api/following-list/", handlers.GetFollowingListForPost)
	mux.HandleFunc("/api/following/", handlers.GetFollowing)
	mux.HandleFunc("/api/follower/", handlers.GetFollower)

	mux.HandleFunc("/api/create-group/", handlers.CreateGroupHandler)
	mux.HandleFunc("/api/invite-to-group/", handlers.InviteToGroupHandler)
	mux.HandleFunc("/api/join-group/", handlers.JoinGroupHandler)
	mux.HandleFunc("/api/leave-group/", handlers.LeaveGroupHandler)
	mux.HandleFunc("/api/get-my-groups", handlers.GetMyGroupsHandler)
	mux.HandleFunc("/api/get-all-groups", handlers.GetAllGroupsHandler)
	mux.HandleFunc("/api/get-group-posts/", handlers.GetGroupPostsHandler)
	mux.HandleFunc("/api/get-group-info/", handlers.GetGroupInfoHandler)
	mux.HandleFunc("/api/check-membership/", handlers.CheckMembershipHandler)
	mux.HandleFunc("/api/create-group-post/", handlers.CreateGroupPostHandler)
	mux.HandleFunc("/api/viewer-status/", handlers.ViewerStatusHandler)

	mux.HandleFunc("/api/create-event/", handlers.CreateEventHandler)
	mux.HandleFunc("/api/get-events/", handlers.GetEventsHandler)
	mux.HandleFunc("/api/update-attendees-status/", handlers.UpdateAttendeesStatus)
	mux.HandleFunc("/api/get-attendees-status/", handlers.GetAttendeesStatus)
	mux.HandleFunc("/api/unread-messages", handlers.GetUnreadMessagesHandler)

	// WebSocket endpoint for notifications
	mux.HandleFunc("/api/notifications/ws", handlers.NotificationWebSocketHandler)
	mux.HandleFunc("/api/chat-notifications/ws", handlers.ChatNotificationWebSocketHandler)

	mux.HandleFunc("/api/get-chats", handlers.GetChatsHandler)
	mux.HandleFunc("/api/get-messages", handlers.GetMessagesHandler)
	mux.HandleFunc("/api/chat/ws", handlers.ChatWebSocketHandler)

	// Serve static files from the ./uploads directory without directory listing
	uploadsDir := http.Dir("./uploads")
	// Ensure the handler is mounted at /uploads/ to serve all subdirectories
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", handlers.CustomFileServer(uploadsDir)))

	return corsMiddleware(sessionMiddleware(mux))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Pre-flight request handling
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		noAuthRequired := []string{"/api/login", "/api/register", "/uploads/", "/api/notifications/ws"}

		path := r.URL.Path
		for _, p := range noAuthRequired {
			if path == p {
				next.ServeHTTP(w, r)
				return
			}
		}

		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		isValid := db.ValidateSessionToken(sessionToken.Value)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
