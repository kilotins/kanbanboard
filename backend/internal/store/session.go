package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"kanbanboard/internal/model"
)

// ErrSessionExpired is returned when a session token has expired.
var ErrSessionExpired = errors.New("session expired")

// ErrSessionNotFound is returned when a session token is not found.
var ErrSessionNotFound = errors.New("session not found")

// CreateSession creates a new session for the given user.
func CreateSession(db *sql.DB, userID string, duration time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	expiresAt := time.Now().Add(duration)
	_, err = db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)",
		token, userID, expiresAt,
	)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	// Clean up expired sessions on login
	_ = DeleteExpiredSessions(db)

	return token, nil
}

// GetSession retrieves a session by token. Returns ErrSessionNotFound if not found,
// or ErrSessionExpired if the session has expired.
func GetSession(db *sql.DB, token string) (model.Session, error) {
	var s model.Session
	err := db.QueryRow(
		"SELECT token, user_id, created_at, expires_at FROM sessions WHERE token = $1",
		token,
	).Scan(&s.Token, &s.UserID, &s.CreatedAt, &s.ExpiresAt)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Session{}, ErrSessionNotFound
	}
	if err != nil {
		return model.Session{}, fmt.Errorf("get session: %w", err)
	}

	if time.Now().After(s.ExpiresAt) {
		_ = DeleteSession(db, token)
		return model.Session{}, ErrSessionExpired
	}

	return s, nil
}

// DeleteSession removes a session by token.
func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE token = $1", token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func DeleteExpiredSessions(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM sessions WHERE expires_at < NOW()")
	if err != nil {
		return fmt.Errorf("delete expired sessions: %w", err)
	}
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
