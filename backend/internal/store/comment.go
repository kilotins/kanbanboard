package store

import (
	"database/sql"
	"errors"
	"fmt"

	"kanbanboard/internal/model"
)

// ErrCommentNotFound is returned when a comment is not found.
var ErrCommentNotFound = errors.New("comment not found")

// CreateComment inserts a new comment.
func CreateComment(db *sql.DB, comment model.Comment) (model.Comment, error) {
	err := db.QueryRow(`
		INSERT INTO comments (task_id, author_id, text)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, comment.TaskID, comment.AuthorID, comment.Text,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return model.Comment{}, fmt.Errorf("create comment: %w", err)
	}
	return comment, nil
}

// CommentWithAuthor extends Comment with the author's name.
type CommentWithAuthor struct {
	model.Comment
	AuthorName string `json:"authorName"`
}

// ListCommentsForTask returns all comments for a task with author names.
func ListCommentsForTask(db *sql.DB, taskID string) ([]CommentWithAuthor, error) {
	rows, err := db.Query(`
		SELECT c.id, c.task_id, c.author_id, c.text, c.created_at, c.updated_at, u.name
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.task_id = $1
		ORDER BY c.created_at
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []CommentWithAuthor
	for rows.Next() {
		var c CommentWithAuthor
		if err := rows.Scan(&c.ID, &c.TaskID, &c.AuthorID, &c.Text, &c.CreatedAt, &c.UpdatedAt, &c.AuthorName); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// GetComment retrieves a comment by ID.
func GetComment(db *sql.DB, commentID string) (model.Comment, error) {
	var c model.Comment
	err := db.QueryRow(`
		SELECT id, task_id, author_id, text, created_at, updated_at
		FROM comments WHERE id = $1
	`, commentID).Scan(&c.ID, &c.TaskID, &c.AuthorID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Comment{}, ErrCommentNotFound
	}
	if err != nil {
		return model.Comment{}, fmt.Errorf("get comment: %w", err)
	}
	return c, nil
}

// UpdateComment updates a comment's text.
func UpdateComment(db *sql.DB, commentID, text string) (model.Comment, error) {
	var c model.Comment
	err := db.QueryRow(`
		UPDATE comments SET text = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, task_id, author_id, text, created_at, updated_at
	`, text, commentID).Scan(&c.ID, &c.TaskID, &c.AuthorID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return model.Comment{}, fmt.Errorf("update comment: %w", err)
	}
	return c, nil
}

// DeleteComment removes a comment by ID.
func DeleteComment(db *sql.DB, commentID string) error {
	_, err := db.Exec("DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}
