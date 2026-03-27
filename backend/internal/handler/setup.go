package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kanbanboard/internal/model"
	"kanbanboard/internal/store"
	"kanbanboard/internal/validate"

	"golang.org/x/crypto/bcrypt"
)

type setupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	AppTitle string `json:"appTitle"`
}

// HandleSetupStatus returns whether initial setup is required.
func HandleSetupStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := store.CountUsers(db)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to check setup status")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"setupRequired": count == 0,
		})
	}
}

// HandleSetup creates the initial admin user and sets the application title.
func HandleSetup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check that no users exist
		count, err := store.CountUsers(db)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to check setup status")
			return
		}
		if count > 0 {
			writeError(w, http.StatusConflict, "Setup has already been completed")
			return
		}

		// Parse request
		var req setupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate
		if req.Name == "" || req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "Name, email, and password are required")
			return
		}

		if msg := validate.Password(req.Password); msg != "" {
			writeError(w, http.StatusBadRequest, msg)
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		// Create admin user
		user := model.User{
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: string(hash),
			IsAdmin:      true,
			IsActive:     true,
		}

		user, err = store.CreateUser(db, user)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to create admin user")
			return
		}

		// Set application title
		title := req.AppTitle
		if title == "" {
			title = "Kanban Board"
		}
		if err := store.SetSetting(db, "app_title", title); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save application title")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// HandleAppTitle returns the application title.
func HandleAppTitle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := store.GetSetting(db, "app_title", "Kanban Board")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to get application title")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"title": title,
		})
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
