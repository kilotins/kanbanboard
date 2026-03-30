package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"kanbanboard/internal/store"

	"golang.org/x/crypto/bcrypt"
)

const sessionDuration = 7 * 24 * time.Hour // 7 days
const cookieName = "session_token"

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleLogin authenticates a user and creates a session.
func HandleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		// Find user
		user, err := store.GetUserByEmail(db, req.Email)
		if errors.Is(err, store.ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to authenticate")
			return
		}

		// Check active
		if !user.IsActive {
			writeError(w, http.StatusUnauthorized, "Account is deactivated")
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			writeError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Create session
		token, err := store.CreateSession(db, user.ID, sessionDuration)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to create session")
			return
		}

		// Set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(sessionDuration.Seconds()),
		})

		writeJSON(w, http.StatusOK, user)
	}
}

// HandleLogout destroys the current session.
func HandleLogout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err == nil {
			_ = store.DeleteSession(db, cookie.Value)
		}

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// HandleMe returns the currently authenticated user.
func HandleMe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Not authenticated")
			return
		}

		session, err := store.GetSession(db, cookie.Value)
		if err != nil {
			// Clear invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   -1,
			})
			writeError(w, http.StatusUnauthorized, "Not authenticated")
			return
		}

		user, err := store.GetUserByID(db, session.UserID)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Not authenticated")
			return
		}

		writeJSON(w, http.StatusOK, user)
	}
}
