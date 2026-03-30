package store

import (
	"database/sql"
	"errors"
	"fmt"

	"kanbanboard/internal/model"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

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

// GetUserByEmail retrieves a user by email address.
// Excludes soft-deleted users (prevents login as deleted user).
func GetUserByEmail(db *sql.DB, email string) (model.User, error) {
	var u model.User
	err := db.QueryRow(`
		SELECT id, name, email, password_hash, is_admin, is_team_manager, is_active, deleted_at, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
	`, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.IsTeamManager, &u.IsActive, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

// GetUserByID retrieves a user by ID.
// Includes soft-deleted users (needed for displaying creator/author names).
func GetUserByID(db *sql.DB, id string) (model.User, error) {
	var u model.User
	err := db.QueryRow(`
		SELECT id, name, email, password_hash, is_admin, is_team_manager, is_active, deleted_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.IsTeamManager, &u.IsActive, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

// ListUsers returns all users including soft-deleted (for admin listing).
func ListUsers(db *sql.DB) ([]model.User, error) {
	rows, err := db.Query(`
		SELECT id, name, email, password_hash, is_admin, is_team_manager, is_active, deleted_at, created_at, updated_at
		FROM users ORDER BY (deleted_at IS NOT NULL), name
	`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.IsTeamManager, &u.IsActive, &u.DeletedAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// BasicUser holds the minimal user fields for listings.
type BasicUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ListActiveUsersBasic returns active, non-deleted users with only id, name, and email.
func ListActiveUsersBasic(db *sql.DB) ([]BasicUser, error) {
	rows, err := db.Query(`
		SELECT id, name, email FROM users WHERE is_active = true AND deleted_at IS NULL ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("list active users: %w", err)
	}
	defer rows.Close()

	var users []BasicUser
	for rows.Next() {
		var u BasicUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// UpdateUserAdmin updates a user's name, email, active status, and roles (admin operation).
func UpdateUserAdmin(db *sql.DB, user model.User) (model.User, error) {
	err := db.QueryRow(`
		UPDATE users SET name = $1, email = $2, is_active = $3, is_admin = $4, is_team_manager = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`, user.Name, user.Email, user.IsActive, user.IsAdmin, user.IsTeamManager, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("update user admin: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user's name and email.
func UpdateUser(db *sql.DB, user model.User) (model.User, error) {
	err := db.QueryRow(`
		UPDATE users SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`, user.Name, user.Email, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}
	return user, nil
}

// UpdatePassword updates a user's password hash.
func UpdatePassword(db *sql.DB, userID, passwordHash string) error {
	_, err := db.Exec(
		"UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2",
		passwordHash, userID,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

// SoftDeleteUser marks a user as deleted, clears their email and password, and deactivates them.
// Email is cleared to allow reuse. Name is preserved for historical references.
func SoftDeleteUser(db *sql.DB, userID string) error {
	_, err := db.Exec(`
		UPDATE users SET deleted_at = NOW(), email = 'deleted_' || id, password_hash = '', is_active = false, updated_at = NOW()
		WHERE id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	return nil
}

// ListProjectsOwnedByUser returns all projects directly owned by a user.
func ListProjectsOwnedByUser(db *sql.DB, userID string) ([]model.Project, error) {
	rows, err := db.Query(`
		SELECT id, name, visibility, tag, next_task_number, owner_user_id, owner_team_id, created_at, updated_at
		FROM projects WHERE owner_user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list projects owned by user: %w", err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Visibility, &p.Tag, &p.NextTaskNumber, &p.OwnerUserID, &p.OwnerTeamID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// UnassignTasksForUser clears the assignee on all tasks assigned to a user.
func UnassignTasksForUser(db *sql.DB, userID string) error {
	_, err := db.Exec(
		"UPDATE tasks SET assignee_id = NULL, updated_at = NOW() WHERE assignee_id = $1",
		userID,
	)
	if err != nil {
		return fmt.Errorf("unassign tasks: %w", err)
	}
	return nil
}

// RemoveUserFromAllTeams removes a user from all team memberships.
func RemoveUserFromAllTeams(db *sql.DB, userID string) error {
	_, err := db.Exec("DELETE FROM team_members WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("remove from all teams: %w", err)
	}
	return nil
}
