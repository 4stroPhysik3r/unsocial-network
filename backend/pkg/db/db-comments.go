package db

import (
	"fmt"
	"log"
	"time"
)

type Comment struct {
	CommentID    int       `json:"comment_id,omitempty"`
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id,omitempty"`
	Content      string    `json:"content"`
	CommentImage *string   `json:"comment_image,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	FullName     string    `json:"full_name"`
}

func InsertComment(comment Comment) error {
	statement, err := DB.Prepare(`INSERT INTO comments (post_id, user_id, content, comment_image) VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing insert statement for comment: %v", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(comment.PostID, comment.UserID, comment.Content, comment.CommentImage)
	if err != nil {
		log.Printf("Error executing insert statement for comment: %v", err)
		return err
	}

	fmt.Println("Comment inserted to DB successfully", comment.CommentID)

	return nil
}

func GetCommentsForPost(postID int) ([]Comment, error) {
	var comments []Comment

	// Adjust the query to select comments and user names based on postID
	query := `SELECT c.comment_id, c.post_id, c.user_id, c.content, c.comment_image, c.created_at, 
                     u.firstname || ' ' || u.lastname AS full_name
              FROM comments c
              JOIN users u ON c.user_id = u.user_id
              WHERE c.post_id = ?`

	rows, err := DB.Query(query, postID)
	if err != nil {
		log.Printf("Error querying database for comments: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the retrieved comments
	for rows.Next() {
		var comment Comment
		// Adjust the Scan to include the full_name
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CommentImage, &comment.CreatedAt, &comment.FullName)
		if err != nil {
			log.Printf("Error scanning comment row: %v", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over comment rows: %v", err)
		return nil, err
	}

	return comments, nil
}
