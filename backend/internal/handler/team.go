package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"kanbanboard/internal/middleware"
	"kanbanboard/internal/model"
	"kanbanboard/internal/store"
)

type teamRequest struct {
	Name string `json:"name"`
}

type addMemberRequest struct {
	UserID string `json:"userId"`
}

// requireTeamOwner loads a team and checks that the user is the owner.
// Returns the team and true if allowed, or writes an error response and returns false.
func requireTeamOwner(db *sql.DB, w http.ResponseWriter, teamID, userID string) (model.Team, bool) {
	team, err := store.GetTeam(db, teamID)
	if errors.Is(err, store.ErrTeamNotFound) {
		writeError(w, http.StatusNotFound, "Team not found")
		return model.Team{}, false
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get team")
		return model.Team{}, false
	}
	if !IsTeamOwner(team, userID) {
		writeError(w, http.StatusForbidden, "Only the team owner can perform this action")
		return model.Team{}, false
	}
	return team, true
}

// HandleListTeams returns teams owned by the current user.
func HandleListTeams(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())

		teams, err := store.ListTeamsForUser(db, user.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to list teams")
			return
		}
		if teams == nil {
			teams = []model.Team{}
		}

		writeJSON(w, http.StatusOK, teams)
	}
}

// HandleCreateTeam creates a new team. User must be a team manager.
func HandleCreateTeam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())

		if !user.IsTeamManager {
			writeError(w, http.StatusForbidden, "Team manager role required")
			return
		}

		var req teamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "Team name is required")
			return
		}

		team := model.Team{Name: req.Name, OwnerID: user.ID}
		team, err := store.CreateTeam(db, team)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to create team")
			return
		}

		writeJSON(w, http.StatusCreated, team)
	}
}

// HandleUpdateTeam renames a team. Must be the owner.
func HandleUpdateTeam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())
		teamID := r.PathValue("teamId")

		team, ok := requireTeamOwner(db, w, teamID, user.ID)
		if !ok {
			return
		}

		var req teamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "Team name is required")
			return
		}

		team.Name = req.Name
		team, err := store.UpdateTeam(db, team)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to update team")
			return
		}

		writeJSON(w, http.StatusOK, team)
	}
}

// HandleDeleteTeam deletes a team if it has no projects.
func HandleDeleteTeam(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())
		teamID := r.PathValue("teamId")

		if _, ok := requireTeamOwner(db, w, teamID, user.ID); !ok {
			return
		}

		count, err := store.CountProjectsForTeam(db, teamID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to check team projects")
			return
		}
		if count > 0 {
			writeError(w, http.StatusConflict, "Cannot delete team that owns projects. Transfer or delete the projects first.")
			return
		}

		if err := store.DeleteTeam(db, teamID); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to delete team")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// HandleListTeamMembers returns members of a team.
func HandleListTeamMembers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())
		teamID := r.PathValue("teamId")

		if _, ok := requireTeamOwner(db, w, teamID, user.ID); !ok {
			return
		}

		members, err := store.ListTeamMembers(db, teamID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to list members")
			return
		}
		if members == nil {
			members = []model.User{}
		}

		writeJSON(w, http.StatusOK, members)
	}
}

// HandleAddTeamMember adds a user to a team.
func HandleAddTeamMember(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())
		teamID := r.PathValue("teamId")

		if _, ok := requireTeamOwner(db, w, teamID, user.ID); !ok {
			return
		}

		var req addMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		if req.UserID == "" {
			writeError(w, http.StatusBadRequest, "User ID is required")
			return
		}

		if err := store.AddTeamMember(db, teamID, req.UserID); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to add member")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// HandleRemoveTeamMember removes a user from a team.
func HandleRemoveTeamMember(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := middleware.UserFromContext(r.Context())
		teamID := r.PathValue("teamId")
		memberID := r.PathValue("userId")

		if _, ok := requireTeamOwner(db, w, teamID, user.ID); !ok {
			return
		}

		if err := store.RemoveTeamMember(db, teamID, memberID); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to remove member")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
