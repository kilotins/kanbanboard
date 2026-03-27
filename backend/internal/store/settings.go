package store

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetSetting retrieves a setting value by key. Returns the fallback if not found.
func GetSetting(db *sql.DB, key, fallback string) (string, error) {
	var value string
	err := db.QueryRow("SELECT value FROM app_settings WHERE key = $1", key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return fallback, nil
	}
	if err != nil {
		return "", fmt.Errorf("get setting %s: %w", key, err)
	}
	return value, nil
}

// SetSetting creates or updates a setting.
func SetSetting(db *sql.DB, key, value string) error {
	_, err := db.Exec(`
		INSERT INTO app_settings (key, value) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
	`, key, value)
	if err != nil {
		return fmt.Errorf("set setting %s: %w", key, err)
	}
	return nil
}
