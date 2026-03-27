package store

import (
	"database/sql"
	"fmt"

	"kanbanboard/internal/model"
)

// CountUsers returns the total number of users in the database.
func CountUsers(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

// CreateUser inserts a new user and returns it with the generated ID and timestamps.
func CreateUser(db *sql.DB, user model.User) (model.User, error) {
	err := db.QueryRow(`
		INSERT INTO users (name, email, password_hash, is_admin, is_team_manager, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, user.Name, user.Email, user.PasswordHash, user.IsAdmin, user.IsTeamManager, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}
