package db

import (
	"database/sql"
	"log"
)

type Post struct {
	PostID       int     `json:"post_id,omitempty"`
	UserID       int     `json:"user_id"`
	GroupID      *int    `json:"group_id,omitempty"`
	Content      string  `json:"content"`
	PostImage    *string `json:"post_image,omitempty"`
	PrivacyLevel *string `json:"privacy_level,omitempty"`
	CreatedAt    string  `json:"created_at"`
	ViewerIDs    []int   `json:"viewer_ids,omitempty"`
	FullName     string  `json:"full_name,omitempty"`
}

func InsertPost(post Post) (int, error) {
	statement, err := DB.Prepare(`INSERT INTO posts (user_id, group_id, content, post_image, privacy_level) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing insert statement: %v", err)
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(post.UserID, post.GroupID, post.Content, post.PostImage, post.PrivacyLevel)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return 0, err
	}

	// fmt.Println("Post inserted to DB successfully with postID:", postID)

	return int(postID), nil
}

func GetPostsForFeed(userID int) ([]Post, error) {
	var posts []Post

	// Fetch public posts
	publicPostsQuery := `SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.privacy_level, p.created_at, 
                         (u.firstname || ' ' || u.lastname) AS full_name
                         FROM posts p
                         JOIN users u ON p.user_id = u.user_id
                         WHERE p.privacy_level = 'public'`

	// Use a helper function to reduce code duplication
	publicPosts, err := fetchPosts(publicPostsQuery, userID)
	if err != nil {
		return nil, err
	}

	// Append the results of public posts to the initial posts slice
	posts = append(posts, publicPosts...)

	// Fetch private posts
	privatePostsQuery := `SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.privacy_level, p.created_at, 
                      (u.firstname || ' ' || u.lastname) AS full_name
                      FROM posts p
                      JOIN follows f ON p.user_id = f.following_id
                      JOIN users u ON p.user_id = u.user_id
                      WHERE f.follower_id = ? AND f.status = 'accepted' AND p.privacy_level = 'private' AND p.group_id IS NULL`

	privatePosts, err := fetchPosts(privatePostsQuery, userID)
	if err != nil {
		return nil, err
	}

	// Append the results of private posts to the initial posts slice
	posts = append(posts, privatePosts...)

	// Fetch friends (viewer) posts
	viewerPostsQuery := `SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.privacy_level, p.created_at, 
                         (u.firstname || ' ' || u.lastname) AS full_name
                         FROM posts p
                         JOIN post_viewers pv ON p.post_id = pv.post_id
                         JOIN users u ON p.user_id = u.user_id
                         WHERE pv.viewer_id = ?`

	viewerPosts, err := fetchPosts(viewerPostsQuery, userID)
	if err != nil {
		return nil, err
	}

	// Append the results of friends (viewer) posts to the initial posts slice
	posts = append(posts, viewerPosts...)

	// Fetch friends (viewer) posts
	creatorPostsQuery := `SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.privacy_level, p.created_at,
	(u.firstname || ' ' || u.lastname) AS full_name
	FROM posts p
	JOIN users u ON p.user_id = u.user_id
	WHERE p.user_id = ? AND p.privacy_level = 'friends'
	`

	creatorPosts, err := fetchPosts(creatorPostsQuery, userID)
	if err != nil {
		return nil, err
	}

	// Append the results of friends (viewer) posts to the initial posts slice
	posts = append(posts, creatorPosts...)

	// // Fetch friends (viewer) posts
	// creatorPrivateQuery := `SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.privacy_level, p.created_at,
	// (u.firstname || ' ' || u.lastname) AS full_name
	// FROM posts p
	// JOIN users u ON p.user_id = u.user_id
	// WHERE p.user_id = ? AND p.privacy_level = 'private' AND p.group_id IS NULL
	// `

	// creatorPrivate, err := fetchPosts(creatorPrivateQuery, userID)
	// if err != nil {
	// 	return nil, err
	// }

	// // Append the results of friends (viewer) posts to the initial posts slice
	// posts = append(posts, creatorPrivate...)

	return posts, nil
}

// fetchPosts is a helper function to execute the provided query and fetch posts
func fetchPosts(query string, userID int) ([]Post, error) {
	var posts []Post

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.GroupID, &post.Content, &post.PostImage, &post.PrivacyLevel, &post.CreatedAt, &post.FullName)
		if err != nil {
			return nil, err
		}

		// Check if the privacy level exists (not a group post)
		if post.PrivacyLevel != nil {
			posts = append(posts, post)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostsForProfile(userID int) ([]Post, error) {
	var posts []Post

	query := "SELECT post_id, user_id, group_id, content, post_image, privacy_level, created_at FROM posts WHERE user_id = ?"

	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Printf("Error querying database for posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.GroupID, &post.Content, &post.PostImage, &post.PrivacyLevel, &post.CreatedAt)
		if err != nil {
			log.Printf("Error scanning post row: %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over post rows: %v", err)
		return nil, err
	}

	return posts, nil
}

func GetPostsForUser(userID, loggedID int) ([]Post, error) {
	var posts []Post

	query := `
	SELECT post_id, user_id, group_id, content, post_image, privacy_level, created_at
	FROM posts
	WHERE user_id = ? 
	  AND privacy_level IS NOT NULL 
	  AND (privacy_level = 'public' 
	       OR user_id = ? 
	       OR EXISTS (SELECT 1 FROM follows WHERE follower_id = ? AND following_id = posts.user_id AND status = 'accepted') 
	       OR EXISTS (SELECT 1 FROM post_viewers WHERE post_id = posts.post_id AND viewer_id = ?))
	`
	// Query for posts belonging to the specified user and considering privacy settings
	rows, err := DB.Query(query, userID, userID, loggedID, loggedID)
	if err != nil {
		log.Printf("Error querying database for posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the retrieved posts
	for rows.Next() {
		var post Post

		err := rows.Scan(&post.PostID, &post.UserID, &post.GroupID, &post.Content, &post.PostImage, &post.PrivacyLevel, &post.CreatedAt)
		if err != nil {
			log.Printf("Error scanning post row: %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over post rows: %v", err)
		return nil, err
	}

	return posts, nil
}

func InsertPostViewers(postID int, viewerIDs []int) error {
	statement, err := DB.Prepare("INSERT INTO post_viewers (post_id, viewer_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	for _, viewerID := range viewerIDs {
		_, err := statement.Exec(postID, viewerID)
		if err != nil {
			return err
		}
	}

	return nil
}

// to display one post
func GetPostByPostID(postID int) (Post, error) {
	var post Post

	query := `
	SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.created_at, p.privacy_level, CONCAT(u.firstname, ' ', u.lastname) AS full_name
	FROM posts p
	JOIN users u ON p.user_id = u.user_id
	WHERE p.post_id = ?
`
	err := DB.QueryRow(query, postID).Scan(&post.PostID, &post.UserID, &post.GroupID, &post.Content, &post.PostImage, &post.CreatedAt, &post.PrivacyLevel, &post.FullName)
	if err != nil {
		log.Printf("Error querying database for post: %v", err)
		return Post{}, err
	}

	return post, nil
}

// GetPostsByGroupID retrieves posts that belong to a specific group by its groupID
func GetPostsByGroupID(groupID int) ([]Post, error) {
	var posts []Post

	query := `
    SELECT p.post_id, p.user_id, p.group_id, p.content, p.post_image, p.created_at, CONCAT(u.firstname, ' ', u.lastname) AS full_name
    FROM posts p
    JOIN users u ON p.user_id = u.user_id
    WHERE p.group_id = ?`

	rows, err := DB.Query(query, groupID)
	if err != nil {
		log.Printf("Error querying database for posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.GroupID, &post.Content, &post.PostImage, &post.CreatedAt, &post.FullName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return posts, nil
}

func GetViewerStatus(postID, userID int) (bool, error) {
	var postUserID int

	err := DB.QueryRow("SELECT user_id FROM posts WHERE post_id = ?", postID).Scan(&postUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if postUserID == userID {
		return true, nil
	}

	query := "SELECT EXISTS(SELECT 1 FROM post_viewers WHERE post_id = ? AND viewer_id = ?)"
	var exists bool
	err = DB.QueryRow(query, postID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
